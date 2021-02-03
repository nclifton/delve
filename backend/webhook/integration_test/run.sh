#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd $DIR &> /dev/null

export $(grep -v '^#' ../rpc/infra/values/dev-values.env | xargs)

export MIGRATION_ROOT="file://../migration/sql"
export TEST_FIXTURE_POSTGRES_USER=gnomock
export TEST_FIXTURE_POSTGRES_USER_PASSWORD=gnomick
export TEST_FIXTURE_RABBITMQ_USER=gnomock
export TEST_FIXTURE_RABBITMQ_USER_PASSWORD=gnomick

go test -timeout 300s -tags integration -v -run ^Test_.*$

popd &> /dev/null