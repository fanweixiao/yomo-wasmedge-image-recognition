package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/yomorun/y3-codec-golang"

	"github.com/yomorun/yomo/pkg/client"
)

var (
	zipperHost = getEnvString("YOMO_ZIPPER_HOST", "localhost")
	zipperPort = getEnvInt("YOMO_ZIPPER_Port", 9000)
	codec      = y3.NewCodec(0x10)
)

const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var coder = base64.NewEncoding(base64Table)

//func Base64Encode(encode_byte []byte) []byte {
//	return []byte(coder.EncodeToString(encode_byte))
//}

func main() {
	fmt.Println("Go: Args:", os.Args)
	filePath := os.Args[1]

	// connect to yomo-zipper.
	cli, err := client.NewSource("image-recognition-source").Connect(zipperHost, zipperPort)
	if err != nil {
		log.Printf("❌ Emit the data to yomo-zipper failure with err: %v", err)
		return
	}
	defer cli.Close()

	loadImageAndSendData(cli, filePath)
}

func loadImageAndSendData(stream io.Writer, filePath string) {
	// load image.
	img, _ := ioutil.ReadFile(filePath)

	for {
		// encode image.
		img64 := coder.EncodeToString(img)

		fmt.Printf("img64=%v\n", img64)
		sendingBuf, _ := codec.Marshal(img64)

		// end data via QUIC stream.
		_, err := stream.Write(sendingBuf)
		if err != nil {
			log.Printf("❌ Send %v to yomo-zipper failure with err: %v", filePath, err)
		} else {
			log.Printf("✅ Send %v to yomo-zipper, hash=%s, img64_size=%v", filePath, genSha1(img), len(img64))
		}

		time.Sleep(2 * time.Second)
	}
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
