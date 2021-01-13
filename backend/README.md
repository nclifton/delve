# Backend

## Services

### gRPC - Proto3

Do you know what gRPC and Proto3 are?

ref: https://grpc.io/about/
ref: https://grpc.io/docs/languages/go/

## Code Generation from proto file

When a service's container is running, `modd` will be able to use the service's proto file (`<service>/rpc/<service>pb/<service>.proto`) to re-generate go code (`<service>/rpc/<service>pb/<service>_grpc.pb.go` and `<service>/rpc/<service>pb/<service>.pb.go`).


Well behaved developers will have all the project's containers running using the project's docker-compose, so for them they should not need to be concerned about such details.


PS: But if a developer wants to live dangerously and without friends, and wants to use local protobuf code generation they must ensure the `protoc` version is consistent with the version being used in the docker service container. Failure to do so will rain death and damnation upon themselves and their family.
