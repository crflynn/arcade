VERSION=0.1.0

.PHONY: setup
setup:
	brew install asdf || True
	asdf install

.PHONY: install
install:
	go mod download
	go mod tidy

.PHONY: run
run:
	go fmt ./...
	go run main.go

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	docker-compose build

.PHONY: up
serve:
	docker-compose up

.PHONY: docker
docker:
	docker build -t crflynn/arcade .
	docker tag crflynn/arcade crflynn/arcade:$(VERSION)
	docker tag crflynn/arcade crflynn/arcade:latest

.PHONY: push
push:
	docker push crflynn/arcade:$(VERSION)
	docker push crflynn/arcade:latest