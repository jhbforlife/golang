FROM golang:alpine AS builder

MAINTAINER Jacob Brown @jhbforlife

RUN apk update && apk upgrade && apk add git
RUN git clone https://github.com/jhbforlife/golang && cd golang/cmd/translate/server && env GOOS=linux GOARCH=amd64 go build

FROM alpine:latest

EXPOSE 8080

WORKDIR /server

COPY --from=builder /go/golang/cmd/translate/server/server .

CMD ["./server"]
