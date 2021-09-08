# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ../api ./
COPY ../configs/env.auth ./
COPY ../configs/auth.yml ./

COPY ../cmd/auth/*.go ./
COPY ../internal/auth ./
COPY ../internal/common ./

RUN go build -o /auth

EXPOSE 3002
EXPOSE 5002

CMD [ "/auth" ]
