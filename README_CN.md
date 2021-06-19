# Streaming Image Recognition by WebAssembly

该项目是一个Show Case，展示如何借助WebAssembly技术，实时解析视频流，并将每一帧的图片调用深度学习模型，判断该帧中是否存在食物。

项目使用的相关技术：

- 流式计算框架是使用[YoMo Streaming Serverless Framework](https://github.com/yomorun/yomo)构建
- Serverless通过[WasmEdge](https://github.com/WasmEdge/WasmEdge)引入WebAssembly，执行深度学习模型
- **TODO** 深度学习模型来自于 

## 如何运行

### 1. Clone Repository

```bash
git clone git@github.com:yomorun/yomo-wasmedge-image-recognition.git
```

### 2. 安装YoMo CLI

```bash
$ go install github.com/yomorun/cli/yomo@latest
```

执行下面的命令，确保yomo已经在环境变量中，如果有任何问题可以参考[YoMo的详细文档](https://github.com/yomorun/yomo)

```bash
$ yomo -v

YoMo CLI version: v0.0.1

```

当然也可以直接下载可执行文件: [Linux](https://github.com/yomorun/yomo-app-image-recognition-example/releases/download/v0.1.0/yomo) [MacOS](https://github.com/yomorun/yomo-app-image-recognition-example/releases/download/v0.1.0/yomo)

### 3. 安装WasmEdge

**TODO** 将具体[安装过程官方文档](https://github.com/second-state/WasmEdge-go#option-2-install-the-pre-released-library) 的必要文档贴在这里

### 4. 编写 Streaming Serverless

**TODO** @chenjunbiao 这里补充一下过程吧（我没找到`.pb`文件怎么来的，`.wasm`文件貌似应该是生成的，可以放在.gitignore中)

```bash
$ cd flow
$ go get -u github.com/second-state/WasmEdge-go/wasmedge
```

### 5. 运行YoMo Streaming Orchestrator

```bash
$ yomo serve -c ./zipper/workflow.yaml
```

### 6. 通过 WasmEdge 运行 Streaming Serverless

```bash
$ cd flow
$ go run --tags tensorflow app.go
```

### 7. 模拟视频流并查看运行结果

```bash
$ go run ./source/main.go ./source/hot-dog.mp4
```

