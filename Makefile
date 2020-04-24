.PHONY: dep run lint build

run:
	go run main.go

dep:
	go mod download
	go mod verify

build:
	go build .

lint:
	go fmt ./...
