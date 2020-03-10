
.PHONY: test install

lint:
	github.com/mdempsky/maligned

install:
	go get ./...

test: install
	go test -race -cover -v ./...
