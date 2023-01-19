.PHONY: all test testc format

all: format test

test:
	go test -v ./...

testc:
	go test -v -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

format:
	go fmt ./...
	goarrange run -r .
