#!/bin/sh
GOPATH=$(go env GOPATH)
GOBIN=$(echo "${GOPATH%%:*}/bin" | sed s#//*#/#g)
LINT_DIRS='
  api/...
  lib/...
'
echo "GOBIN = $GOBIN"
echo "PWD = $(echo $(pwd))"
echo "LINT_DIRS = $LINT_DIRS"
$GOBIN/golangci-lint --deadline=5m run $LINT_DIRS
