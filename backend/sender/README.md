# Sender

Sender service uses gRPC protobuf to define the RPC messages and service end points.

see [sender.proto](./rpc/senderpb/sender.proto)

## Sender Integration Test
Tests the operations of the sender service

see [integration test readme](./integration_test/README.md)

## Mock RPC Service Client

Generated mock for the ServiceClient for use in API unit tests and possibly in integration tests of domain services where the sender service is a dependency.

mock service client (`senderpb.MockServiceClient`) is generated using [ vektra/mockery ](https://github.com/vektra/mockery/blob/master/pkg/generator.go).

command line:

```
cd backend/sender/rpc/senderpb
mockery --name ServiceClient --inpackage
```