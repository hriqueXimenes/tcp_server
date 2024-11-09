#!/bin/bash
BUILD_DIR="../build"

if [ ! -d "$BUILD_DIR" ]; then
  mkdir "$BUILD_DIR"
fi

rm -f "$BUILD_DIR/sumologic_server" "$BUILD_DIR/sumologic_server.exe" "$BUILD_DIR/sumologic_server_mac"

echo "Compiling for Linux..."
GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/sumologic_server" ../main.go

echo "Compiling for Windows..."
GOOS=windows GOARCH=amd64 go build -o "$BUILD_DIR/sumologic_server.exe" ../main.go

echo "Compiling for macOS..."
GOOS=darwin GOARCH=amd64 go build -o "$BUILD_DIR/sumologic_server_mac" ../main.go

echo "Build completed. Executables are located in the $BUILD_DIR directory."