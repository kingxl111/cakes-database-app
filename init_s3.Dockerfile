FROM alpine:latest

RUN apk update && apk add --no-cache \
    bash \
    curl \
    jq \
    && curl -sSL https://dl.min.io/client/mc/release/linux-amd64/mc > /usr/bin/mc \
    && chmod +x /usr/bin/mc \
    && rm -rf /var/cache/apk/*

COPY init_s3.sh /init_s3.sh
RUN chmod +x /init_s3.sh

ENTRYPOINT ["/bin/sh", "/init_s3.sh"]
