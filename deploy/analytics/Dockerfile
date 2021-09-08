# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ../api ./
COPY ../configs/env.analytics ./
COPY ../configs/api.yml ./

COPY ../cmd/analytics/*.go ./
COPY ../internal/analytics ./
COPY ../internal/common ./

RUN go build -o /analytics

EXPOSE 3001
EXPOSE 5001

CMD [ "/analytics" ]
