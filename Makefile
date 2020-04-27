.PHONY: dep run lint build simulate

run:
	go run main.go

dep:
	go mod download
	go mod verify

build:
	go build .

lint:
	go fmt ./...

simulate:
	go run main.go simulate
