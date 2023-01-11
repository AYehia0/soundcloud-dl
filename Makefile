TRACK_URL := https://soundcloud.com/sobhi-mohamed5/99-118-mp4
DL_PATH := ./download
LANG=en_US.UTF-8
SHELL=/bin/bash

run:
	go run main.go ${TRACK_URL} ${DL_PATH}
build:
	go build -o bin/sc-dl main.go
test:
	go test -v # No tests for now
compile:
	echo "Compiling for multiple Platforms"
	GOOS=linux GOARCH=386 go build -o bin/sc-dl-linux main.go
	GOOS=windows GOARCH=386 go build -o bin/sc-dl-windoos main.go
