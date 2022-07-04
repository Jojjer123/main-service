# main-service

## Development

### Generation of protobuf
#### Only structures
```
protoc --go_out=. --go_opt=paths=source_relative filename.proto
```
#### With gRPC server/client
```
protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative filename.proto
```