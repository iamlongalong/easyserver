#!/bin/bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/easyserver-darwin-x86_64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/easyserver-darwin-arm64 .

# CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o build/easyserver-windows-386 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/easyserver-windows-amd64 .

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/easyserver-linux-x86_64 .
# CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o build/easyserver-linux-arm .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/easyserver-linux-arm64 .

# router
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o build/easyserver-linux-mips .
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -o build/easyserver-linux-mips64 .
