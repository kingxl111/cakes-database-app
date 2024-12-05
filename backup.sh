#!/bin/bash

set -e

PGUSER="${DB_USER}"
PGPASSWORD="${DB_PASSWORD}"
PGHOST="${DB_HOST}"
PGDATABASE="${DB_NAME}"

export PGPASSWORD

BACKUP_FILE="/backups/db_backup_$(date +'%Y%m%d%H%M%S').sql"

pg_dump -h "${PGHOST}" -U "${PGUSER}" -d "${PGDATABASE}" -F c > "${BACKUP_FILE}"

if [ $? -eq 0 ]; then
  echo "Backup successful: ${BACKUP_FILE}"
  find "/backups" -name "db_backup_*.sql" -type f -mmin +10 -exec echo rm {} \;
  find "/backups" -name "db_backup_*.sql" -type f -mmin +10 -exec rm {} \;
else
  echo "Backup failed!" >&2
  exit 1
fi
