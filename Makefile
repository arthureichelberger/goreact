.DEFAULT_GOAL := help
GRC=$(shell which grc)
VERSION = ${VERSION}

help: ## Makefile help
help:
	@grep -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

build: ## Build the server and the client
build: build-client build-server

build-client: ## Build the react client
build-client: 
	rm -rf server/cmd/build
	npm -C client run build
	mv client/build server/cmd/build

build-server: ## Build the go binary
build-server:
	go build -o deployment/server server/cmd/server.go

build-server-linux: ## Build the go binary for linux
	@CGO_ENABLED=0 GOOS=linux go build -o deployment/server server/cmd/server.go

build-linux: ## Build for linux
build-linux: build-client build-server-linux

docker: ## Build the docker image
docker: build-linux
	docker build -t goreact:$(VERSION) ./deployment