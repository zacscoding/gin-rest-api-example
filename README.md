# Rest API with golang and gin, gorm  
This project is an exemplary rest api server built with Go :)  

See [API Spec](./api.md) (specification modified from [RealWorld API Spec](https://github.com/gothinkster/realworld/tree/master/api) to simplify)  

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

; TBD  

---  

## TODO  

- [x] add user, article with comment api spec
- [x] add common error response
- [x] configure project layer
- [x] impl account db
- [x] impl account handler (binding, serialize, common error middleware, etc...)
- [x] impl article db
- [ ] impl article handler
- [ ] configure docker compose
- [ ] configure tests (newman or http)