# Integration Test For GRPC Account Service

## References

Uses gnomock containerised dependency services for RabbitMQ and Postgres

 - https://github.com/orlangure/gnomock

## Running test: 
```
cd account/integration_test
go test -timeout 30s .
```