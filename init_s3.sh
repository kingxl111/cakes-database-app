#!/bin/bash

echo "Waiting for MinIO to start..."
until curl -s http://minio:9000; do
  sleep 2
done

BUCKET_NAME=${S3_BUCKET}

mc alias set myminio http://minio:9000 "${S3_ACCESS_KEY}" "${S3_SECRET_KEY}"

if ! mc ls myminio/"${BUCKET_NAME}" > /dev/null 2>&1; then
  echo "Creating bucket ${BUCKET_NAME}"
  mc mb myminio/"${BUCKET_NAME}"
else
  echo "Bucket ${BUCKET_NAME} already exists."
fi

echo "Setting public policy on bucket ${BUCKET_NAME}"
mc anonymous set public myminio/"${BUCKET_NAME}"

echo "Checking anonymous policy for ${BUCKET_NAME}"
mc anonymous info myminio/"${BUCKET_NAME}"

if [ -z "$(ls /images/*.jpg 2>/dev/null)" ]; then
  echo "No .jpg files found in /images"
else
  for file in /images/*.jpg; do
    echo "Uploading $file to ${BUCKET_NAME}"
    mc cp "$file" myminio/"${BUCKET_NAME}"/"$(basename "$file")"
    echo "Public URL: http://localhost:9000/${BUCKET_NAME}/$(basename "$file")"
  done
fi

echo "S3 bucket ${BUCKET_NAME} is set up and files are uploaded."
