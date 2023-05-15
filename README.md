# go-rest-boilerplate

# Running backend

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
