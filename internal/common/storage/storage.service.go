package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/baothaihcmut/Storage-app/internal/config"
)

type PresignUrlMethod string

const (
	PresignUrlGetMethod PresignUrlMethod = "get_method"
	PresignUrlPutMethod PresignUrlMethod = "put_method"
)

type GetPresignUrlArg struct {
	Method PresignUrlMethod
	Key    string
}

type StorageService interface {
	GetPresignUrl(context.Context, GetPresignUrlArg) (string, error)
	GetStorageProviderName() string
	GetStorageBucket() string
}

type S3StorageService struct {
	client *s3.Client
	logger logger.Logger
	cfg    *config.S3Config
}

func (s *S3StorageService) GetPresignUrl(ctx context.Context, args GetPresignUrlArg) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	if args.Method == PresignUrlPutMethod {
		url, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
			Key:    aws.String(args.Key),
			Bucket: aws.String(s.cfg.Bucket),
		}, s3.WithPresignExpires(time.Hour*3))
		if err != nil {
			s.logger.Errorf(ctx, map[string]interface{}{
				"key":    args.Key,
				"bucket": s.cfg.Bucket,
			}, "Error get presign url for put object:", err)
		}
		return url.URL, nil
	} else {
		url, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
			Key:    aws.String(args.Key),
			Bucket: aws.String(s.cfg.Bucket),
		}, s3.WithPresignExpires(time.Hour*3))
		if err != nil {
			s.logger.Errorf(ctx, map[string]interface{}{
				"key":    args.Key,
				"bucket": s.cfg.Bucket,
			}, "Error get presign url for get object:", err)
		}
		return url.URL, nil
	}
}
func (s *S3StorageService) GetStorageProviderName() string {
	return "S3"
}

func (s *S3StorageService) GetStorageBucket() string {
	return s.cfg.Bucket
}
func NewS3StorageService(client *s3.Client, logger logger.Logger, cfg *config.S3Config) StorageService {
	return &S3StorageService{
		client: client,
		logger: logger,
		cfg:    cfg,
	}
}
