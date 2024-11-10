FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    apk add curl &&  \
    apk add tar &&    \
    rm -rf /var/cache/apk/*


RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz

RUN chmod +x migrate && \
    mv migrate /bin/migrate

WORKDIR /app

ADD /internal/storage/migrations/*.sql migrations/
ADD migration_local.sh .
ADD local.env .

RUN chmod +x migration_local.sh

ENTRYPOINT ["bash", "migration_local.sh"]