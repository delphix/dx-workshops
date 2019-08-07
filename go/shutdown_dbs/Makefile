# Makefile

local: 
	go run *.go -c conf.txt

local-debug:
	go run *.go -vvv -c config.txt

build:
	echo $$GOPATH
	go get -d
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -v -o ./bin/linux64/shutdown_dbs

	