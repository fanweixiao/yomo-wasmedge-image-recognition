package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/yomorun/yomo/pkg/client"

	"github.com/yomorun/y3-codec-golang"
)

var (
	codec   = y3.NewCodec(0x10)
	coder64 = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
)

const (
	packetSize = 1024
)

type ImageData struct {
	ImageID       string `y3:"0x12"`
	ImageHash     string `y3:"0x13"`
	PacketCount   int64  `y3:"0x14"`
	PacketId      int64  `y3:"0x15"`
	PacketContent string `y3:"0x17"`
}

func main() {
	fmt.Println("Go: Args:", os.Args)
	filePath := os.Args[1]

	// connect to yomo-zipper.
	cli, err := client.NewSource("image-recognition-source").Connect("localhost", 9000)
	if err != nil {
		log.Printf("❌ Emit the data to yomo-zipper failure with err: %v", err)
		return
	}
	defer cli.Close()

	loadImageAndSendData(cli, filePath)
}

func loadImageAndSendData(stream io.Writer, filePath string) {
	for {
		// load image data
		img, _ := ioutil.ReadFile(filePath)
		img64 := coder64.EncodeToString(img)
		groups := split([]rune(img64), int64(len(img64)/packetSize))

		count := len(groups)
		data := ImageData{
			ImageID:     genUUID(),
			ImageHash:   genSha1(img),
			PacketCount: int64(count),
		}

		// send data via QUIC stream.
		for i := 0; i < count; i++ {
			r := groups[i]
			data.PacketId = int64(i)
			data.PacketContent = string(r)

			sendingBuf, _ := codec.Marshal(data)
			_, err := stream.Write(sendingBuf)
			if err != nil {
				log.Printf("❌ Send %v [%v] to yomo-zipper failure with err: %v", data.ImageID, i, err)
			} else {
				log.Printf("✅ Send %v [%v] to yomo-zipper, ContentSize=%v, ImageHash=%v",
					data.ImageID, i, len(data.PacketContent), data.ImageHash)
			}
			time.Sleep(1 * time.Millisecond)
		}

		time.Sleep(2 * time.Second)
	}
}

func genSha1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func genUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func split(runes []rune, num int64) [][]rune {
	max := int64(len(runes))
	if max < num {
		return nil
	}
	var segments = make([][]rune, 0)
	quantity := max / num
	end := int64(0)
	for i := int64(1); i <= num; i++ {
		qu := i * quantity
		if i != num {
			segments = append(segments, runes[i-1+end:qu])
		} else {
			segments = append(segments, runes[i-1+end:])
		}
		end = qu - i
	}
	return segments
}
