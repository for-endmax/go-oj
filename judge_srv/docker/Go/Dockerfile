FROM golang:1.20-alpine

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"

RUN mkdir -p /judge
WORKDIR /judge


RUN go install std

ENTRYPOINT ["sleep", "infinity"]
