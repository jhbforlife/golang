FROM alpine:latest AS builder

MAINTAINER Jacob Brown

RUN apk update && apk upgrade
RUN apk add go && apk add git
RUN git clone https://github.com/jhbforlife/golang && cd golang/cmd/translate/cli && go build -o translate

FROM alpine:latest

COPY --from=builder /golang/cmd/translate/cli/translate /bin

RUN apk update && apk upgrade

CMD ["sh"]