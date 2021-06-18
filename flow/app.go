package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/yomorun/yomo/pkg/client"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/rx"
)

var (
	m   = map[string][]ImageData{}
	mux sync.Mutex

	coder64 = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
)

const (
	ImageDataKey = 0x10
)

type ImageData struct {
	ImageID       string `y3:"0x12"`
	ImageHash     string `y3:"0x13"`
	PacketCount   int64  `y3:"0x14"`
	PacketId      int64  `y3:"0x15"`
	PacketContent string `y3:"0x17"`
}

func main() {
	cli, err := client.NewServerless("image-recognition").Connect("localhost", 9000)
	if err != nil {
		log.Print("❌ Connect to zipper failure: ", err)
		return
	}

	defer cli.Close()
	cli.Pipe(Handler)
}

func Handler(rxstream rx.RxStream) rx.RxStream {
	stream := rxstream.
		Subscribe(ImageDataKey).
		OnObserve(decode).
		Encode(0x11)

	return stream
}

var decode = func(v []byte) (interface{}, error) {
	mux.Lock()
	defer mux.Unlock()

	// parse ImageData
	var mold ImageData
	err := y3.ToObject(v, &mold)
	if err != nil {
		return nil, err
	}

	// gather packages
	dataArray, ok := m[mold.ImageID]
	if ok == false {
		dataArray = make([]ImageData, 0)
		dataArray = append(dataArray, mold)
		m[mold.ImageID] = dataArray
	} else {
		dataArray = append(dataArray, mold)
		m[mold.ImageID] = dataArray
	}

	packetCount := dataArray[0].PacketCount

	if int64(len(dataArray)) == packetCount {
		groupId := dataArray[0].ImageID
		imageHash := dataArray[0].ImageHash

		// combine packages
		img64 := ""
		for _, item := range dataArray {
			img64 += item.PacketContent
		}
		delete(m, groupId)

		// restore image
		img, err := coder64.DecodeString(img64)

		//create file
		err = ioutil.WriteFile("./temp.jpg", img, 0644)
		if err != nil {
			return nil, err
		}

		log.Printf("✅ received image %s, cal_hash=%s, imageHash=%s\n", groupId, genSha1(img), imageHash)

		return true, nil
	}

	return false, nil
}

func genSha1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}
