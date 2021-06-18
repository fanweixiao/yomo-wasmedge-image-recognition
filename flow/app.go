package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/second-state/WasmEdge-go/wasmedge"
	"github.com/yomorun/yomo/pkg/client"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/rx"
)

var (
	coder64 = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	vm      *wasmedge.VM
	vmConf  *wasmedge.Configure
	counter uint64
)

const ImageDataKey = 0x10

func main() {
	initVM()
	defer vm.Delete()
	defer vmConf.Delete()

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

	res, err := vm.ExecuteBindgen("infer", wasmedge.Bindgen_return_array, img)
	if err == nil {
		fmt.Println("GO: Run bindgen -- infer:", string(res.([]byte)))
	} else {
		fmt.Println("GO: Run bindgen -- infer FAILED")
	}

	hash := genSha1(img)
	log.Printf("✅ received image-%d hash %v, img64_size=%d \n", atomic.AddUint64(&counter, 1), hash, len(img64))

	return hash, nil
}

func genSha1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func initVM() {
	wasmedge.SetLogErrorLevel()
	/// Set Tensorflow not to print debug info
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "3")
	os.Setenv("TF_CPP_MIN_VLOG_LEVEL", "3")

	/// Create configure
	vmConf = wasmedge.NewConfigure(wasmedge.WASI)

	/// Create VM with configure
	vm = wasmedge.NewVMWithConfig(vmConf)

	/// Init WASI
	var wasi = vm.GetImportObject(wasmedge.WASI)
	wasi.InitWasi(
		os.Args[1:],     /// The args
		os.Environ(),    /// The envs
		[]string{".:."}, /// The mapping directories
		[]string{},      /// The preopens will be empty
	)

	/// Register WasmEdge-tensorflow
	var tfobj = wasmedge.NewTensorflowImportObject()
	var tfliteobj = wasmedge.NewTensorflowLiteImportObject()
	vm.RegisterImport(tfobj)
	vm.RegisterImport(tfliteobj)

	/// Instantiate wasm
	vm.LoadWasmFile("rust_mobilenet_food_lib_bg.wasm")
	vm.Validate()
	vm.Instantiate()
}
