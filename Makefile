.PHONY: clean install test dist

clean:
	rm -rf dist

install:
	go get -u github.com/kardianos/govendor
	govendor sync

test:
	go test -v ./...

dist:
	mkdir dist
	GOOS=darwin GOARCH=amd64 go build -o dist/secretly-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -o dist/secretly-linux-amd64
