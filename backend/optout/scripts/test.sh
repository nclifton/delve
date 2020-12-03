#!/bin/bash

export OPTOUT_TEST_DATABASE_URL="postgres://postgres:example@localhost:54321/optout"

echo "Unit test running..."

go test ./...

echo "Integration test running..."

go test -tags=integration ./...
