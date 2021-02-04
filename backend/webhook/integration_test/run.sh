#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd $DIR &> /dev/null

go test -timeout 300s -tags integration -v -run ^Test_.*$
status=$?

popd &> /dev/null

exit $status