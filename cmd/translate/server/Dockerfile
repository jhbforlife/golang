FROM golang:latest AS builder

MAINTAINER Jacob Brown @jhbforlife

RUN apt update && apt upgrade -y
RUN git clone https://github.com/jhbforlife/golang && cd golang/cmd/translate/server && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM alpine:latest

EXPOSE 8080

WORKDIR /server

COPY --from=builder /go/golang/cmd/translate/server/server .

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ["./server"]
