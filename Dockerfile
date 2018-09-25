FROM golang:1.11-alpine3.7 AS build

ENV CGO_ENABLED=1
ENV GOOS=linux
# ENV GO111MODULE=on

WORKDIR /go/src/github.com/mvl-at/rest
RUN apk add --no-cache \
    git \
    musl-dev \
    build-base
COPY . /go/src/github.com/mvl-at/rest
RUN go get ./...
RUN go install -ldflags '-s -w' ./cmd/serve

# ---

FROM alpine
COPY --from=build /go/bin/serve /rest
WORKDIR /rest-data
VOLUME  /rest-data
EXPOSE  7301
ENTRYPOINT [ "/rest" ]
