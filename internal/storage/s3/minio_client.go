package s3

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client     *minio.Client
	bucketName string
	baseURL    string
	publicURL  string
}

func NewS3Client(endpoint, accessKey, secretKey, bucketName, region, publicUrl string) (ClientS3, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	return &MinioClient{
		client:     client,
		bucketName: bucketName,
		baseURL:    fmt.Sprintf("http://%s/%s", endpoint, bucketName),
		publicURL:  publicUrl,
	}, nil
}

func (s *MinioClient) UploadFile(ctx context.Context, file multipart.File, fileName string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, fileName, file, -1, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return s.GetFileURL(fileName), nil
}

func (s *MinioClient) GetFileURL(fileName string) string {
	return fmt.Sprintf("%s/%s/%s.jpg", s.publicURL, s.bucketName, fileName)
}
