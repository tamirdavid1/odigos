1. generate pb.go files from the .proto file

```shell
cd k8sutils/pkg
protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative instrumented_process/proto/instrumented_process.proto
```