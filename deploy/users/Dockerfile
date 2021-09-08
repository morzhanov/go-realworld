# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ../api ./
COPY ../configs/env.users ./
COPY ../configs/users.yml ./

COPY ../cmd/users/*.go ./
COPY ../internal/users ./
COPY ../internal/common ./

RUN go build -o /users

EXPOSE 3004
EXPOSE 5004

CMD [ "/users" ]
