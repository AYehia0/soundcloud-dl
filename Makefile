TRACK_URL :=  https://soundcloud.com/alheweny-official/alminshawy_al-aaraf103_128
PL_URL := https://soundcloud.com/adam00alakad/sets/k1rc41mibizn 
DL_PATH := ./download
LANG=en_US.UTF-8
SHELL=/bin/bash

run:
	rm -f download/*
	go run main.go ${TRACK_URL} --download-path ${DL_PATH} --quality mp3
build:
	go build -o bin/sc-dl main.go
test:
	go test ./pkg/* -v # No tests for now
install:
	go install
compile:
	echo "Compiling for multiple Platforms"
	GOOS=linux GOARCH=386 go build -o bin/sc-dl-linux main.go
	GOOS=windows GOARCH=386 go build -o bin/sc-dl-windoos main.go
