# syntax=docker/dockerfile:1

FROM debian:latest

WORKDIR /user_hub

COPY ./ ./

RUN apt-get update && apt-get install -y ca-certificates

ADD server.crt /container/cert/path

RUN update-ca-certificates

EXPOSE 5623

EXPOSE 465

WORKDIR /user_hub/src

CMD [ "./executable" ]