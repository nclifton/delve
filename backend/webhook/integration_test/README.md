# Integration Test For GRPC Webhook Service

## References

Uses gnomock containerised dependency services for RabbitMQ and Postgres

 - https://github.com/orlangure/gnomock

## Running test: 
```
cd webhook/integration_test
go test -timeout 30s -tags integration -run ^Test_.*$  github.com/burstsms/mtmo-tp/backend/webhook/integration_test
```

## Notes:
requires environment variables set for the tests
```
        RABBIT_EXCHANGE=webhook
        RABBIT_EXCHANGE_TYPE=direct
        MIGRATION_ROOT="file://../migration/sql"
        TEST_FIXTURE_POSTGRES_USER=gnomock
        TEST_FIXTURE_POSTGRES_USER_PASSWORD=gnomick
        TEST_FIXTURE_RABBITMQ_USER=gnomock
        TEST_FIXTURE_RABBITMQ_USER_PASSWORD=gnomick

```

The publish integration test should assert that the http requests that are emitted conform to the expected patterns defined in the [webhook specs](../spec/webhooks.md)