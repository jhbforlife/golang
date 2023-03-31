FROM alpine:latest AS builder

MAINTAINER Jacob Brown @jhbforlife

RUN apk update && apk upgrade
RUN apk add go && apk add git
RUN git clone https://github.com/jhbforlife/golang && cd golang/cmd/translate/server && go build

FROM alpine:latest

EXPOSE 8080

WORKDIR /server

COPY --from=builder /golang/cmd/translate/server/server .

CMD ["./server"]