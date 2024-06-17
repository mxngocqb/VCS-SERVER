# VCS-SERVER

## Run service

### Migrate Postgres Databse
``` 
go run ./cmd/migration
```

### API Service
```
go run ./cmd/api
```

### Report Servicee
```
go run ./cmd/report
```

## Genarate code

### Create Proto

```
protoc -Ipkg/service/report/proto --go_out=. --go_opt=module=github.com/mxngocqb/VCS-SERVER/back-end  --go-grpc_out=. --go-grpc_opt=module=github.com/mxngocqb/VCS-SERVER/back-end pkg/service/report/proto/report.proto
```

### Initialize Swagger

```
swag init -g ./cmd/api/main.go --parseDependency true
```