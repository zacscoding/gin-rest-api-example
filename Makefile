MODULE = $(shell go list -m)

.PHONY: generate build test lint build-docker compose compose-down migrate
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

compose.%:
	$(eval CMD = ${subst compose.,,$(@)})
	tools/script/compose.sh $(CMD)

migrate:
	docker run --rm -v migrations:/migrations --network host migrate/migrate -path=/migrations/ \
	-database mysql://root:password@localhost:3306/local_db?charset=utf8&parseTime=True&multiStatements=true up 2

