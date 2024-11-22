FROM alpine:3.18

RUN apk add --no-cache \
    postgresql15-client \
    bash \
    dcron \
    tini

COPY backup.sh /usr/local/bin/backup.sh
RUN chmod +x /usr/local/bin/backup.sh

COPY backup-cron /etc/cron.d/backup-cron

RUN chmod 0644 /etc/cron.d/backup-cron && crontab /etc/cron.d/backup-cron

RUN mkdir -p /backups
VOLUME /backups

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["crond", "-f"]
