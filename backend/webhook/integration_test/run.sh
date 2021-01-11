#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd $DIR &> /dev/null

export RABBIT_EXCHANGE=webhook
export RABBIT_EXCHANGE_TYPE=direct
export MIGRATION_ROOT="file://../migration/sql"
export TEST_FIXTURE_POSTGRES_USER=gnomock
export TEST_FIXTURE_POSTGRES_USER_PASSWORD=gnomick
export TEST_FIXTURE_RABBITMQ_USER=gnomock
export TEST_FIXTURE_RABBITMQ_USER_PASSWORD=gnomick

go test -timeout 30s -tags integration -run ^Test_.*$

popd &> /dev/null