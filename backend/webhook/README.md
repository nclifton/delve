# Webhook

Webhook service

## gRPC - Proto3

Do you know what are gRPC and Proto3?

ref: https://grpc.io/about/
ref: https://grpc.io/docs/languages/go/

## Prerequisite:

ref https://grpc.io/docs/protoc-installation/

> Warning
>
> Check the version of protoc (as indicated below) after installation to ensure that it is sufficiently recent. The versions of protoc installed by some package managers can be quite dated.
> 
> Installing from pre-compiled binaries, as indicated below, is the best way to ensure that you’re using the latest release of protoc.

## Install pre-compiled binaries (any OS)

To install the latest release of the protocol compiler from pre-compiled binaries, follow these instructions:

1. Manually download from github.com/google/protobuf/releases the zip file corresponding to your operating system and computer architecture (protoc-<version>-<os><arch>.zip), or fetch the file using commands such as the following: (assuming version is 3.14.0, change as required)
```
PB_REL="https://github.com/protocolbuffers/protobuf/releases"
curl -LO $PB_REL/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip
```
2. Unzip the file under $HOME/.local or a directory of your choice. For example: 
```
unzip protoc-3.14.0-linux-x86_64.zip -d $HOME/.local
```
3. Update your environment’s path variable to include the path to the protoc executable. For example:
```
export PATH="$PATH:$HOME/.local/bin"
```

## Go plugins for the protocol compiler:

Install the protocol compiler plugins for Go using the following commands:

```
export GO111MODULE=on  # Enable module mode
go get google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

Update your PATH so that the protoc compiler can find the plugins:

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Script to Generate Go Code from the Proto File:
```
cd backend/webhook
.generate.sh
```

## Webhook Integration Test
Tests the operations of the webhook service

see [integration test readme](./integration_test/README.md)