FROM golang:1.13.0-stretch

RUN apt-get install -y ca-certificates git

RUN env GO111MODULE=on go get -u github.com/cortesi/modd/cmd/modd
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.19.1

WORKDIR /backend
ENV DISABLE_NEW_RELIC true
CMD /go/bin/modd 
