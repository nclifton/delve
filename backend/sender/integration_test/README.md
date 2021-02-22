# Integration Test For GRPC Sender Service

## References

Uses gnomock containerised dependency services for RabbitMQ and Postgres

 - https://github.com/orlangure/gnomock

## Running test: 
```
cd sender/integration_test
go test -timeout 30s -tags integration -run ^Test_.*$  github.com/burstsms/mtmo-tp/backend/sender/integration_test
```