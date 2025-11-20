package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOConfig 定义 MinIO 初始化参数
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	UseSSL    bool
	BaseURL   string
}

var (
	minioClient *minio.Client
	minioCfg    MinIOConfig
)

// InitMinIO 初始化 MinIO 客户端并确保桶存在
func InitMinIO(cfg MinIOConfig) error {
	if cfg.Endpoint == "" {
		return errors.New("minio endpoint is required")
	}
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return errors.New("minio credentials are required")
	}
	if cfg.Bucket == "" {
		return errors.New("minio bucket is required")
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return fmt.Errorf("failed to check minio bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{
			Region: cfg.Region,
		}); err != nil {
			return fmt.Errorf("failed to create minio bucket: %w", err)
		}
	}

	minioClient = client
	minioCfg = cfg
	return nil
}

// Upload 上传文件到 MinIO，返回可访问 URL
func Upload(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	if minioClient == nil {
		return "", errors.New("minio client is not initialized")
	}

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	if _, err := minioClient.PutObject(ctx, minioCfg.Bucket, objectName, reader, objectSize, opts); err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	return buildFileURL(objectName), nil
}

// BucketName 返回默认桶名称
func BucketName() string {
	return minioCfg.Bucket
}

// ObjectURL 根据 objectName 返回完整访问地址
func ObjectURL(objectName string) string {
	return buildFileURL(objectName)
}

func buildFileURL(objectName string) string {
	baseURL := strings.TrimRight(minioCfg.BaseURL, "/")
	if baseURL == "" {
		scheme := "http"
		if minioCfg.UseSSL {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, minioCfg.Endpoint)
	}

	return fmt.Sprintf("%s/%s/%s", baseURL, minioCfg.Bucket, objectName)
}
