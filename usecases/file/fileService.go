package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileService struct {
	s3Client *s3.Client
	ctx context.Context
}

func NewFileService(s3Client *s3.Client, ctx context.Context) FileService {
	return FileService{s3Client, ctx}
}

func (s *FileService) UploadToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	key := fmt.Sprintf("%d-%s", time.Now().Unix(), fileHeader.Filename)
	_, err := s.s3Client.PutObject(s.ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(fmt.Sprintf("%d-%s", time.Now().Unix(), fileHeader.Filename)),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    "public-read",
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}
	
	return key, nil
}