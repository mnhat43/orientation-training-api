ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-alpine as builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

ENTRYPOINT ["go", "run", "/app/cmd/frontapi/main.go"]