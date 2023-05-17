# go-rest-boilerplate

## Required packages
- mockgen: `go install github.com/golang/mock/mockgen@v1.6.0`
- golangci-lint: https://golangci-lint.run/usage/install/#local-installation

# Running backend

Run docker for backend with the command: `docker compose up`
To run separate instances of API, DB, use: `docker compose up {instance}`
- Instance being `api` or `postgres`

To run cobra commands `go run main.go {command}`
Available commands:
- `migrate {migrations-dir}`
- `new-migration {name}`
- `print-config`
- `rollback {migrations-dir}`
- `serve`
- `test-worker`
- `runner-schedule`
- `worker`

### Start server
```golang
go run . server
```

### Schema Migration
```golang
go run . migrate
```
