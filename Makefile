BINARY_NAME=ooh-bot
BUILD_DIR=output

.PHONY: test clean

build:
	GOARCH=amd64 GOOS=darwin go build -o "${BUILD_DIR}/${BINARY_NAME}-darwin" .
	GOARCH=amd64 GOOS=linux go build -o "${BUILD_DIR}/${BINARY_NAME}-linux" .
	GOARCH=amd64 GOOS=windows go build -o "${BUILD_DIR}/${BINARY_NAME}-windows" .

clean:
	go clean
	rm -f ${BUILD_DIR}/**/*
	rm -f ${BUILD_DIR}/*

test:
	go test -v internal/bot
	go test -v internal/client
	go test -v internal/parser
	go test -v internal/scraper
	go test -v internal/storage

install:
	chmod +x githooks/**
	cp githooks/* .git/hooks/
