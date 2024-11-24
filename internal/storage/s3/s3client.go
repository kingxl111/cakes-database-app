package s3

import (
	"context"
	"mime/multipart"
)

type ClientS3 interface {
	UploadFile(ctx context.Context, file multipart.File, fileName string) (string, error)
	GetFileURL(fileName string) string
}
