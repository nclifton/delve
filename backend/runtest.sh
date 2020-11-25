#!/bin/bash

set -e

packages='
    api
    tecloo
    tualet
    sms
    account
    webhook
    tecloo_receiver
'

echo "" >coverage_report.txt

for p in $packages; do
  go test -mod=vendor -race -coverprofile=profile.out -covermode=atomic github.com/burstsms/mtmo-tp/backend/$p/...
  if [ -f profile.out ]; then
    cat profile.out >>coverage_report.txt
    rm profile.out
  fi
done
