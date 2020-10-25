FROM golang:1.14

ENV GO111MODULE=on
ENV CONFIG=/etc/advertisement_crud/config.json
WORKDIR /go/src/advertisement_crud

# Prevent overwrite vendor foreach build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
COPY etc/config/config.json /etc/advertisement_crud/config.json

RUN go get github.com/swaggo/swag/cmd/swag
RUN swag init

RUN go get github.com/vektra/mockery/...
RUN mockery -dir controller -name IApplication
RUN mockery -dir application -name IModel