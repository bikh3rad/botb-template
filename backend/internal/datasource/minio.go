package datasource

import (
	"application/app"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// minioConfig is the koanf sub-tree `datasource.minio`. Env overrides follow the
// APP_ convention, e.g. APP_DATASOURCE_MINIO_ENDPOINT → datasource.minio.endpoint.
type minioConfig struct {
	Endpoint      string        `koanf:"endpoint"`
	AccessKey     string        `koanf:"access_key"`
	SecretKey     string        `koanf:"secret_key"`
	Bucket        string        `koanf:"bucket"`
	Region        string        `koanf:"region"`
	UseSSL        bool          `koanf:"use_ssl"`
	PresignExpiry time.Duration `koanf:"presign_expiry"`
}

const defaultPresignExpiry = 15 * time.Minute

// NewMinioConfig loads the `datasource.minio` sub-tree.
func NewMinioConfig(_ context.Context, c *app.KConfig) (*minioConfig, error) {
	cfg := &minioConfig{
		Region:        "us-east-1",
		Bucket:        "botb-media",
		PresignExpiry: defaultPresignExpiry,
	}
	if err := c.Unmarshal("datasource.minio", cfg); err != nil {
		return nil, err
	}

	if cfg.PresignExpiry <= 0 {
		cfg.PresignExpiry = defaultPresignExpiry
	}

	return cfg, nil
}

// ObjectInfo is the storage-agnostic metadata returned by Stat.
type ObjectInfo struct {
	Key         string
	Size        int64
	ContentType string
}

// MinioStorage is an S3-compatible object store backed by MinIO. It is used
// through the (structural) storage interface owned by the media biz layer, so
// it can be swapped for AWS S3 later without touching callers.
type MinioStorage struct {
	client        *minio.Client
	bucket        string
	presignExpiry time.Duration
	logger        *slog.Logger
}

var ErrObjectNotFound = errors.New("object not found")

// NewMinioStorage constructs the client, ensures the bucket exists, and
// registers liveness + shutdown hooks on the shared controller.
func NewMinioStorage(
	ctx context.Context,
	logger *slog.Logger,
	controller app.Controller,
	cfg *minioConfig,
) (*MinioStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region}); err != nil {
			return nil, err
		}
	}

	s := &MinioStorage{
		client:        client,
		bucket:        cfg.Bucket,
		presignExpiry: cfg.PresignExpiry,
		logger:        logger.With("layer", "MinioStorage"),
	}

	controller.RegisterHealthz("minio", s.healthz)
	controller.RegisterShutdown("minio", s.shutdown)

	return s, nil
}

// Put uploads an object of the given size and content type.
func (s *MinioStorage) Put(
	ctx context.Context,
	key string,
	r io.Reader,
	size int64,
	contentType string,
) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

// PresignGet returns a time-limited URL granting read access to the object.
func (s *MinioStorage) PresignGet(ctx context.Context, key string) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, s.bucket, key, s.presignExpiry, url.Values{})
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// Stat returns object metadata, or ErrObjectNotFound if it is missing.
func (s *MinioStorage) Stat(ctx context.Context, key string) (ObjectInfo, error) {
	info, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return ObjectInfo{}, ErrObjectNotFound
		}

		return ObjectInfo{}, err
	}

	return ObjectInfo{Key: info.Key, Size: info.Size, ContentType: info.ContentType}, nil
}

// Remove deletes an object.
func (s *MinioStorage) Remove(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

// Bucket returns the configured bucket name (used when persisting metadata).
func (s *MinioStorage) Bucket() string {
	return s.bucket
}

func (s *MinioStorage) healthz(ctx context.Context) error {
	_, err := s.client.BucketExists(ctx, s.bucket)

	return err
}

func (s *MinioStorage) shutdown(_ context.Context) error {
	s.logger.Info("shutting down MinioStorage")

	return nil
}
