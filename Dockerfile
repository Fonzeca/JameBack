FROM golang:1.23.2-alpine3.20 AS build

WORKDIR /user-hub

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /user-hub/user-hub .

FROM alpine:3.20 AS runtime

WORKDIR /user-hub

COPY --from=build /user-hub/user-hub .

ENTRYPOINT [ "./user-hub" ]