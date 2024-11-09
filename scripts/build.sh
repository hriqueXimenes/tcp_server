#!/bin/bash
BUILD_DIR="../build"

if [ ! -d "$BUILD_DIR" ]; then
  mkdir "$BUILD_DIR"
fi

rm -f "$BUILD_DIR/sumologic_server"

echo "Compiling for Linux..."
GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/sumologic_server" ../main.go

echo "Compiling for Windows..."
GOOS=windows GOARCH=amd64 go build -o "$BUILD_DIR/sumologic_server.exe" ../main.go

echo "Build completed. Executables are located in the $BUILD_DIR directory."