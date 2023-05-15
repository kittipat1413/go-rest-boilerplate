# go-rest-boilerplate

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
