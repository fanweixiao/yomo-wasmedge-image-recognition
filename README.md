# yomo-app-image-recognition-example



## Getting Started

### 1. Zipper

#### Run

```bash
yomo serve -c ./zipper/workflow.yaml
```



### 2. Flow

#### Run

```bash
cd flow
go get -u github.com/second-state/WasmEdge-go/wasmedge
go run --tags tensorflow app.go
```



### 3.Source

#### Run

```bash
go run ./source/main.go ./source/hot-dog.mp4
```



