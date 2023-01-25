FROM alpine:latest AS builder

MAINTAINER Jacob Brown

RUN apk update && apk upgrade
RUN apk add go && apk add git
RUN git clone https://github.com/jhbatshipt/golang && cd golang/cmd/capstone/server && go build

FROM alpine:latest

EXPOSE 8080

WORKDIR /server

COPY --from=builder /golang/cmd/capstone/server/server .

CMD ["./server"]