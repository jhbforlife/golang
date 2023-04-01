FROM golang:latest AS builder

MAINTAINER Jacob Brown @jhbforlife

RUN apt update && apt upgrade -y
RUN git clone https://github.com/jhbforlife/golang && cd golang/cmd/translate/server && env GOOS=linux GOARCH=amd64 go build

FROM ubuntu:latest

EXPOSE 8080

WORKDIR /server

COPY --from=builder /go/golang/cmd/translate/server/server .

RUN apt update && apt upgrade && apt install ca-certificates -y

CMD ["./server"]
