#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd $DIR &> /dev/null

go test -timeout 30s -tags integration -run ^Test_.*$ -v
status=$?

popd &> /dev/null
exit $status
