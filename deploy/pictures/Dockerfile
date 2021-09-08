# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ../api ./
COPY ../configs/env.pictures ./
COPY ../configs/pictures.yml ./

COPY ../cmd/pictures/*.go ./
COPY ../internal/pictures ./
COPY ../internal/common ./

RUN go build -o /auth

EXPOSE 3004
EXPOSE 5004

CMD [ "/pictures" ]
