package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/yomorun/yomo/pkg/client"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/rx"
)

var (
	coder64 = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
)

const ImageDataKey = 0x10

func main() {
	cli, err := client.NewServerless("image-recognition").Connect("localhost", 9000)
	if err != nil {
		log.Print("❌ Connect to zipper failure: ", err)
		return
	}

	defer cli.Close()
	cli.Pipe(Handler)
}

func Handler(rxStream rx.RxStream) rx.RxStream {
	stream := rxStream.
		Subscribe(ImageDataKey).
		OnObserve(decode).
		Encode(0x11)

	return stream
}

var decode = func(v []byte) (interface{}, error) {
	img64, err := y3.ToUTF8String(v)
	if err != nil {
		return nil, err
	}

	img, err := coder64.DecodeString(img64)
	if err != nil {
		return nil, err
	}

	//create file
	err = ioutil.WriteFile(fmt.Sprintf("./.out/%s.jpg", genSha1(img)), img, 0644)
	if err != nil {
		return nil, err
	}

	hash := genSha1(img)
	log.Printf("✅ received image hash %v, img64_size=%d \n", hash, len(img64))

	return hash, nil
}

func genSha1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}
