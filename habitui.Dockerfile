FROM golang:1.22.4-alpine

COPY . /habitui-src

WORKDIR /habitui-src

RUN go build -o /habitui cmd/habitui/main.go

ENV TERM=xterm-256color

ENTRYPOINT ["/habitui"]
