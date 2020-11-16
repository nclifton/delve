FROM golang:1.13.0-stretch

RUN apt-get install -y ca-certificates git
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4 && \
	echo "deb http://repo.mongodb.org/apt/debian stretch/mongodb-org/4.0 main" | tee /etc/apt/sources.list.d/mongodb-org-4.0.list && \
	apt-get update && \
	apt-get install -y mongodb-org 

RUN env GO111MODULE=on go get -u github.com/cortesi/modd/cmd/modd
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.19.1

WORKDIR /backend
ENV DISABLE_NEW_RELIC true
CMD /go/bin/modd 
