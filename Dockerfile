# syntax=docker/dockerfile:1

FROM golang:alpine

WORKDIR /user_hub

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /user_hub/user_hub .

EXPOSE 5623

ENTRYPOINT [ "./user_hub", "server" ]