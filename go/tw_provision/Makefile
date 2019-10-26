# Makefile

local: 
	go run *.go -c conf.txt

local-debug:
	go run *.go -vvv -c config.txt

build:
	echo $$GOPATH
	go get -d
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=1.0.4" -o ./bin/linux64/tw_provision
	
buildmac:
	echo $$GOPATH
	go get -d -v
	env GOOS=darwin GOARCH=amd64 go build -installsuffix 'static' -ldflags="-X main.version=1.0.4" -v -o ./bin/macos/tw_provision
	