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
else
  echo "Backup failed!" >&2
  exit 1
fi
