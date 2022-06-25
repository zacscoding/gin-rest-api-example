![workflow](https://github.com/zacscoding/gin-rest-api-example/actions/workflows/check.yaml/badge.svg)

# Rest API with golang and gin, gorm  
This project is an exemplary rest api server built with Go :)  

See [API Spec](./api.md) (modified from [RealWorld API Spec](https://github.com/gothinkster/realworld/tree/master/api) to simplify)  

API Server technology stack is  

- Server code: `golang`
- REST Server: `gin`
- Database: `MySQL` with `golang-migrate` to migrate  
- ORM: `gorm v2`  
- Dependency Injection: `fx`  
- Unit Testing: `go test` and `testify`
- Configuration management: `cobra`

---  

> ## Getting started  

### 1. Start with docker compose  

> ####  build a docker image  

```bash
// docker build -f cmd/server/Dockerfile -t gin-example/article-server .
$ make build-docker
``` 

> #### run api server with mysql (see docker-compose.yaml)  

```bash
// docker-compose up --force-recreate
$ make compose.up

$  docker ps -a
CONTAINER ID        IMAGE                        COMMAND                  CREATED             STATUS              PORTS                               NAMES
e01564708984        gin-example/article-server   "article-server --co…"   40 seconds ago      Up 39 seconds       0.0.0.0:3000->3000/tcp              article-server
105cb25b6d3a        mysql:8.0.17                 "docker-entrypoint.s…"   40 seconds ago      Up 39 seconds       0.0.0.0:3306->3306/tcp, 33060/tcp   my-mysql

$ make compose.down
```  

> #### Check apis  

Run intellij's .http files in `tools/http/sample directory`(./tools/http/sample)  

---  

## TODO  

- [x] add user, article with comment api spec
- [x] add common error response
- [x] configure project layer
- [x] impl account db
- [x] impl account handler (binding, serialize, common error middleware, etc...)
- [x] impl article db
- [x] impl article handler
- [ ] refactor binding and validation of request
- [x] configure docker compose
- [ ] add metrics
- [ ] configure tests (newman or http)
- [ ] another project layer with different branch