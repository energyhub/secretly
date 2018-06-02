FROM alpine:3.7
LABEL maintainer="dev@energyhub.net"

# necessary for SSL w/ AWS
RUN apk add --no-cache ca-certificates

COPY dist/secretly-linux-amd64 /usr/local/bin/secretly

ENTRYPOINT ["secretly"]

CMD ["env"]

