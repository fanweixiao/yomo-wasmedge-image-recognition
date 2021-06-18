#!/usr/bin/env bash

set -e

sudo apt-get update
sudo apt-get install -y libjpeg-dev libpng-dev build-essential

wget https://github.com/second-state/wasmedge-image/releases/download/0.8.0/WasmEdge-image-deps-0.8.0-manylinux1_x86_64.tar.gz
sudo tar -C /usr/local/lib -zxvf WasmEdge-image-deps-0.8.0-manylinux1_x86_64.tar.gz
sudo ln -sf libjpeg.so.8.3.0 /usr/local/lib/libjpeg.so
sudo ln -sf libjpeg.so.8.3.0 /usr/local/lib/libjpeg.so.8
sudo ln -sf libpng16.so.16.37.0 /usr/local/lib/libpng.so
sudo ln -sf libpng16.so.16.37.0 /usr/local/lib/libpng16.so
sudo ln -sf libpng16.so.16.37.0 /usr/local/lib/libpng16.so.16
sudo ldconfig

wget https://github.com/second-state/WasmEdge-tensorflow-deps/releases/download/0.8.0/WasmEdge-tensorflow-deps-TF-0.8.0-manylinux2014_x86_64.tar.gz
wget https://github.com/second-state/WasmEdge-tensorflow-deps/releases/download/0.8.0/WasmEdge-tensorflow-deps-TFLite-0.8.0-manylinux2014_x86_64.tar.gz
sudo tar -C /usr/local/lib -xzf WasmEdge-tensorflow-deps-TF-0.8.0-manylinux2014_x86_64.tar.gz
sudo tar -C /usr/local/lib -xzf WasmEdge-tensorflow-deps-TFLite-0.8.0-manylinux2014_x86_64.tar.gz
sudo ln -sf libtensorflow.so.2.4.0 /usr/local/lib/libtensorflow.so.2
sudo ln -sf libtensorflow.so.2 /usr/local/lib/libtensorflow.so
sudo ln -sf libtensorflow_framework.so.2.4.0 /usr/local/lib/libtensorflow_framework.so.2
sudo ln -sf libtensorflow_framework.so.2 /usr/local/lib/libtensorflow_framework.so
sudo ldconfig

wget https://github.com/second-state/WasmEdge-tensorflow/releases/download/0.8.0/WasmEdge-tensorflow-0.8.0-manylinux2014_x86_64.tar.gz
wget https://github.com/second-state/WasmEdge-tensorflow/releases/download/0.8.0/WasmEdge-tensorflowlite-0.8.0-manylinux2014_x86_64.tar.gz
sudo tar -C /usr/local/ -xzf WasmEdge-tensorflow-0.8.0-manylinux2014_x86_64.tar.gz
sudo tar -C /usr/local/ -xzf WasmEdge-tensorflowlite-0.8.0-manylinux2014_x86_64.tar.gz
sudo ldconfig

wget https://github.com/WasmEdge/WasmEdge/releases/download/0.8.0/WasmEdge-0.8.0-manylinux2014_x86_64.tar.gz
tar -xzf WasmEdge-0.8.0-manylinux2014_x86_64.tar.gz
sudo cp WasmEdge-0.8.0-Linux/include/wasmedge.h /usr/local/include
sudo cp WasmEdge-0.8.0-Linux/lib64/libwasmedge_c.so /usr/local/lib
sudo ldconfig