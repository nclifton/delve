# This Dockerfile does the following:
# - Adds certs to enable HTTPS
# - Forces build to use vendor folder
# - Compiles into the smallest possible binary
# - Uses scratch for the smallest runnable image
# - Runs the binary as an unpriviledged user

# First stage: Build the executable
FROM golang:alpine AS builder

# Create unpriviledged user
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

RUN apk add --no-cache ca-certificates

# Set the working directory outside $GOPATH to enable the support for modules
WORKDIR /src
COPY ./backend /src/

RUN CGO_ENABLED=0 go build -mod=vendor -ldflags '-s -w' -o /app /src/sms/migration

FROM scratch

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
COPY --from=builder /src/sms/migration/sql /sql

USER nobody:nobody
ENTRYPOINT ["/app"]
