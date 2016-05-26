#!/bin/bash
#GOOS=darwin GOARCH=amd64 go build -o binaries/distributor-osx-arm64
#GOOS=linux GOARCH=arm64 go build -o binaries/distributor-linux-arm64
GOOS=linux GOARCH=amd64 go build -o binaries/distributor-linux-amd64
#GOOS=linux GOARCH=386 go build -o binaries/distributor-linux-386
#GOOS=windows GOARCH=386 go build -o binaries/distributor-windows-386.exe
#GOOS=windows GOARCH=amd64 go build -o binaries/distributor-windows-amd64.exe
