#!/bin/bash
echo 'build for darwin'
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/easyserver-darwin-x86_64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/easyserver-darwin-arm64 .

echo 'build for windows'
# CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o build/easyserver-windows-386 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/easyserver-windows-amd64 .

echo 'build for linux'
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/easyserver-linux-x86_64 .
# CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o build/easyserver-linux-arm .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o build/easyserver-linux-arm64 .

echo 'build for linux mips'
# router
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags="-s -w" -o build/easyserver-linux-mips .
# CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags="-s -w" -o build/easyserver-linux-mips64 .

echo '🎉🎉🎉🎉  finish build'

ls -alh build/*

upx -9 build/*
