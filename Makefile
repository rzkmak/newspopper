# VARS
REPO=rzkmak
TAG=0.0.3
IMAGENAME=newspopper
IMAGEFULLNAME=${REPO}/${IMAGENAME}:${TAG}


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

docker-build:
	docker build -f docker/Dockerfile -t ${IMAGEFULLNAME} .

docker-compose-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-compose-down:
	docker-compose -f docker/docker-compose.yml down

docker-logs-follow:
	docker logs -f newspopper-bot
