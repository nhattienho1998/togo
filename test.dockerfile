FROM golang:1.14-alpine

# Install git
RUN set -ex; \
    apk update; \
    apk add --no-cache git

# Set working directory
WORKDIR /go/src/togo

# Run tests
CMD CGO_ENABLED=0 go test ./...
