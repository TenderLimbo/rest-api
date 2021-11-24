# rest-api
It's a simple REST API application written in Golang using GORM and PostgreSQL. To run you need [docker](https://www.docker.com/) and [golang-migrate](https://github.com/golang-migrate/migrate) to be installed  
## How to run
create .env file from example
```
cp .env.example .env
```
build and run
```
make build && make run
```
if you run for the first time, make migration-up. Before it make sure your PostgreSQL server is stopped.
```
make migration-up
```
you can also migration-down
```
make migration-down
```
stop all containers
```
make stop
```
## In addition
run tests
```
make test
```
linter
```
make lint
```
