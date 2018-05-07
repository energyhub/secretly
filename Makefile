.PHONY: clean install test dist
VERSION:=$(shell git describe --tags --always)

clean:
	rm -rf dist

install:
	go get -u github.com/kardianos/govendor
	govendor sync

test:
	go test -v ./...

dist:
	mkdir dist
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dist/secretly-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/secretly-$(VERSION)-linux-amd64

