package s3

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

var (
	BucketNameEnv = os.Getenv("MINIO_BUCKET_NAME")
	RegionEnv     = os.Getenv("MINIO_REGION")
)

type S3Client struct {
	Client *minio.Client
}

func New() *S3Client {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	secureStr := os.Getenv("MINIO_SECURE")
	secure := true
	if secureStr != "" {
		if s, err := strconv.ParseBool(secureStr); err == nil {
			secure = s
		}
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create S3/MinIO client")
		return nil
	}

	log.Debug().Msg("Connected to S3/MinIO successfully")

	return &S3Client{
		Client: client,
	}
}

func (s *S3Client) GetPresignedURL(key string) (string, error) {
	url, err := s.Client.PresignedGetObject(context.Background(), BucketNameEnv, key, time.Hour*24, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *S3Client) GetPresignedUploadURL(key string, contentType string) (string, error) {
	url, err := s.Client.PresignedPutObject(context.Background(), BucketNameEnv, key, time.Hour*1)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *S3Client) DeleteObject(key string) error {
	return s.Client.RemoveObject(context.Background(), BucketNameEnv, key, minio.RemoveObjectOptions{})
}
