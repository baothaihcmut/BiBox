package storage

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
	"github.com/samber/lo"
)

type PresignUrlMethod string

const (
	PresignUrlGetMethod PresignUrlMethod = "get_method"
	PresignUrlPutMethod PresignUrlMethod = "put_method"
)

type GetPresignUrlArg struct {
	Method      PresignUrlMethod
	Key         string
	ContentType enums.MimeType
	Expiry      time.Duration
	Preview     bool
	fileName    string
}

type StorageService interface {
	GetPresignUrl(context.Context, GetPresignUrlArg) (string, error)
	GetStorageProviderName() string
	GetStorageBucket() string
	GetFile(context.Context, string) (io.ReadCloser, error)
	BulkDelete(context.Context, []string) error
}

type S3StorageService struct {
	client *s3.Client
	logger logger.Logger
	cfg    *config.S3Config
}

// BulkDelete implements StorageService.
func (s *S3StorageService) BulkDelete(ctx context.Context, keys []string) error {
	objectIds := lo.Map(keys, func(item string, _ int) types.ObjectIdentifier {
		return types.ObjectIdentifier{Key: aws.String(item)}
	})
	deleteInput := &s3.DeleteObjectsInput{
		Bucket: aws.String(s.cfg.Bucket),
		Delete: &types.Delete{
			Objects: objectIds,
			Quiet:   aws.Bool(true),
		},
	}
	_, err := s.client.DeleteObjects(context.TODO(), deleteInput)
	if err != nil {
		return err
	}
	return nil

}

func (s *S3StorageService) GetPresignUrl(ctx context.Context, args GetPresignUrlArg) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	if args.Method == PresignUrlPutMethod {
		url, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
			Key:         aws.String(args.Key),
			Bucket:      aws.String(s.cfg.Bucket),
			ContentType: aws.String(string(args.ContentType)),
		}, s3.WithPresignExpires(args.Expiry))
		if err != nil {
			s.logger.Errorf(ctx, map[string]any{
				"key":    args.Key,
				"bucket": s.cfg.Bucket,
			}, "Error get presign url for put object:", err)
		}
		return url.URL, nil
	} else {
		arg := &s3.GetObjectInput{
			Key:                 aws.String(args.Key),
			Bucket:              aws.String(s.cfg.Bucket),
			ResponseContentType: aws.String(string(args.ContentType)),
		}
		if args.Preview {
			arg.ResponseContentDisposition = aws.String("inline")
		} else {
			arg.ResponseContentDisposition = aws.String("attachment")

		}
		url, err := presigner.PresignGetObject(ctx, arg, s3.WithPresignExpires(args.Expiry))
		if err != nil {
			s.logger.Errorf(ctx, map[string]any{
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

func (s *S3StorageService) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(s.cfg.Bucket),
	})
	if err != nil {
		s.logger.Errorf(ctx, map[string]any{
			"key":    key,
			"bucket": s.cfg.Bucket,
		}, "Error get object from storage: ", err)
		return nil, err
	}
	return resp.Body, err

}
func NewS3StorageService(client *s3.Client, logger logger.Logger, cfg *config.S3Config) StorageService {
	return &S3StorageService{
		client: client,
		logger: logger,
		cfg:    cfg,
	}
}
