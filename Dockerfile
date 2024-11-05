FROM golang:alpine

WORKDIR /user-hub

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /user-hub/user-hub .

EXPOSE 5623

ENTRYPOINT [ "./user-hub" ]