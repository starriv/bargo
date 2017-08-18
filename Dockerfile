FROM golang:1.8

LABEL maintainer="https://github.com/dawniii/bargo"

RUN mkdir -p /usr/local/opt/bargo
COPY ./ /usr/local/opt/bargo
WORKDIR /usr/local/opt/bargo/

ENTRYPOINT ["bin/bargo-linux-amd64"]
