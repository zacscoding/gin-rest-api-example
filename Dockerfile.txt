FROM golang:1.14-alpine AS build

RUN mkdir -p /go/src/github.com/zacscoding/gin-rest-api-example ~/.ssh && \
    apk add --no-cache git openssh-client make gcc libc-dev
WORKDIR /go/src/github.com/zacscoding/gin-rest-api-example
COPY . .
RUN make build

FROM alpine:3
COPY --from=build /go/src/github.com/zacscoding/gin-rest-api-example/article-server /bin/article-server
CMD /bin/article-server