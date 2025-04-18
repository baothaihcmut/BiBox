package services

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/baothaihcmut/BiBox/file-process/internal/configs"
)

type StorageService interface {
	GetFile(context.Context, string) (bytes.Buffer, error)
}

type S3Service struct {
	s3    *s3.Client
	s3Cfg *configs.S3Config
}

func NewS3Service(s3 *s3.Client) *S3Service {
	return &S3Service{
		s3: s3,
	}
}

func (s *S3Service) GetFile(ctx context.Context, key string) (bytes.Buffer, error) {
	output, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.s3Cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer output.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, output.Body)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}
