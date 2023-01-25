FROM alpine:latest AS builder

MAINTAINER Jacob Brown

RUN apk update && apk upgrade
RUN apk add go && apk add git
RUN git clone https://github.com/jhbatshipt/golang && cd golang/cmd/capstone/client && go build -o translate

FROM alpine:latest

COPY --from=builder /golang/cmd/capstone/client/translate /bin

RUN apk update && apk upgrade

CMD ["sh"]