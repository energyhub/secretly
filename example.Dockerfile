FROM alpine:3.7
LABEL maintainer="dev@energyhub.net"

# necessary for SSL w/ AWS
RUN apk add --no-cache ca-certificates

ADD https://github.com/energyhub/secretly/releases/download/0.0.6/secretly-linux-amd64 /usr/local/bin/secretly

RUN chmod +x /usr/local/bin/secretly

ENTRYPOINT ["secretly"]
