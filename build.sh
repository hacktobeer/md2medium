#!/bin/bash
mkdir builds
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o builds/md2medium.windows main.go 
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/md2medium.linux main.go 
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o builds/md2medium.darwin main.go 