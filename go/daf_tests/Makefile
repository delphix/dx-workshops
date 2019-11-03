# Makefile

.PHONY: build

local: 
	go run *.go

local-debug:
	go run *.go -vvv

build:
	echo $$GOPATH
	go get -d
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -v -ldflags="-X main.version=1.0.2" -o ./bin/linux64/daf_tests

	