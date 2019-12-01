# Makefile

local: 
	go run *.go -c conf.txt

local-debug:
	go run *.go -vvv -c config.txt

build:
	echo $$GOPATH
	go get -d
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -v -ldflags="-X main.version=1.0.5" -o ./bin/linux64/tw_prep

	