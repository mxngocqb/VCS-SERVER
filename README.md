# VCS-SERVER

```
protoc -Ipkg/service/report/proto --go_out=. --go_opt=module=github.com/mxngocqb/VCS-SERVER/back-end  --go-grpc_out=. --go-grpc_opt=module=github.com/mxngocqb/VCS-SERVER/back-end pkg/service/report/proto/report.proto
```