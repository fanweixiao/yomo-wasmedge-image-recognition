package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/yomorun/y3-codec-golang"

	"github.com/yomorun/yomo/pkg/client"
	"github.com/yomorun/yomo/pkg/rx"
)

var (
	zipperHost = getEnvString("YOMO_ZIPPER_HOST", "localhost")
	zipperPort = getEnvInt("YOMO_ZIPPER_Port", 9000)
)

const ImageDataKey = 0x10

var decode = func(v []byte) (interface{}, error) {
	img64, err := y3.ToUTF8String(v)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("img64=%v\n", img64)
	img, err := base64.StdEncoding.DecodeString(img64)
	if err != nil {
		return nil, err
	}

	//create file
	err = ioutil.WriteFile("./temp.jpg", img, 0644)
	if err != nil {
		return nil, err
	}

	hash := genSha1(img)
	log.Printf("✅ received image hash %v, img64_size=%d \n", hash, len(img64))

	return hash, nil
}

func Handler(rxstream rx.RxStream) rx.RxStream {
	stream := rxstream.
		Subscribe(ImageDataKey).
		OnObserve(decode).
		Encode(0x11)

	return stream
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) != 0 {
		result, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}

		return result
	}
	return defaultValue
}

func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) != 0 {
		return value
	}
	return defaultValue
}

func genSha1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
	cli, err := client.NewServerless("image-recognition").Connect(zipperHost, zipperPort)
	if err != nil {
		log.Print("❌ Connect to zipper failure: ", err)
		return
	}

	defer cli.Close()
	cli.Pipe(Handler)
}
