VERSION := $(shell git describe --tags --always --dirty="dev")

clean:
	rm -rf dist

test:
	go test -v ./...

dist: test
	mkdir dist
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/secretly-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/secretly-$(VERSION)-linux-amd64

