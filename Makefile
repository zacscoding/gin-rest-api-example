MODULE = $(shell go list -m)

.PHONY: generate build test lint build-docker  compose compose-down
generate:
	go generate ./...

build: # build a server
	go build -a -o article-server $(MODULE)/cmd/server

test:
	go clean -testcache
	go test ./... -v

lint:
	gofmt -l .

build-docker: # build docker image
	docker build -f cmd/server/Dockerfile -t gin-example/article-server .

compose: # run with docker-compose
	docker-compose up --force-recreate

compose-down: # down docker-compose
	docker-compose down -v

