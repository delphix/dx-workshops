# Makefile

local: 
	go run *.go -c conf.txt

local-debug:
	go run *.go -vvv -c config.txt

build:
	echo $$GOPATH
	go get -d
	env GOOS=linux GOARCH=amd64 go build -i -v -ldflags="-X main.version=1.0.3" -o ./bin/linux64/tw_prep

	