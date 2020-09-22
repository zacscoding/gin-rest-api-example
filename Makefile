MODULE = $(shell go list -m)

.PHONY: generate build build-docker compose compose-down test
generate:
	go generate ./...

build: # build a server
	go build -a -o article-server $(MODULE)/cmd/server

build-docker: # build docker image
	docker build -f cmd/server/Dockerfile -t gin-example/article-server .

compose: # run with docker-compose
	docker-compose up --force-recreate

compose-down: # down docker-compose
	docker-compose down -v

test:
	go test ./... -v