# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ../api ./
COPY ../configs/env.apigw ./
COPY ../configs/api.yml ./

COPY ../cmd/apigw/*.go ./
COPY ../internal/apigw ./
COPY ../internal/common ./

RUN go build -o /apigw

EXPOSE 3005

CMD [ "/apigw" ]
