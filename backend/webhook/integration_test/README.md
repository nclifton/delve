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
        WEBHOOK_RABBIT_EXCHANGE=webhook
        WEBHOOK_RABBIT_EXCHANGE_TYPE=direct
        WEBHOOK_MIGRATION_ROOT="file://../migration/sql"
        INTEGRATION_TEST_FIXTURE_POSTGRES_USER=gnomock
        INTEGRATION_TEST_FIXTURE_POSTGRES_USER_PASSWORD=gnomick
        INTEGRATION_TEST_FIXTURE_RABBITMQ_USER=gnomock
        INTEGRATION_TEST_FIXTURE_RABBITMQ_USER_PASSWORD=gnomick

```

The publish integration test should assert that the http requests that are emitted conform to the expected patterns defined in the [webhook specs](../spec/webhooks.md)