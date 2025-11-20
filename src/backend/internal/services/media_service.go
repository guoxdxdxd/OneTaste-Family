package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/utils"
	"onetaste-family/backend/pkg/storage"
)

var (
	// ErrUnsupportedMediaType 不支持的文件类型
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	// ErrInvalidMediaSize 文件大小不合法
	ErrInvalidMediaSize = errors.New("invalid media size")
)

// MaxMediaFileSize 单个文件最大 10MB
const MaxMediaFileSize int64 = 10 * 1024 * 1024

// MediaService 媒体业务逻辑
type MediaService struct{}

// NewMediaService 创建 MediaService
func NewMediaService() *MediaService {
	return &MediaService{}
}

// UploadMedia 处理文件上传逻辑
func (s *MediaService) UploadMedia(ctx context.Context, reader io.Reader, input *models.MediaUploadInput) (*models.MediaUploadResponse, error) {
	if reader == nil {
		return nil, fmt.Errorf("invalid reader")
	}

	if input.Size <= 0 || input.Size > MaxMediaFileSize {
		return nil, ErrInvalidMediaSize
	}

	contentType := input.ContentType
	if !models.IsAllowedMediaContentType(contentType) {
		return nil, ErrUnsupportedMediaType
	}

	dir, err := models.CleanStorageDirectory(input.Directory)
	if err != nil {
		return nil, err
	}

	objectKey := buildObjectKey(dir, input.OriginalName, contentType)

	url, err := storage.Upload(ctx, objectKey, reader, input.Size, contentType)
	if err != nil {
		return nil, fmt.Errorf("upload media to storage failed: %w", err)
	}

	return &models.MediaUploadResponse{
		URL:         url,
		Bucket:      storage.BucketName(),
		ObjectKey:   objectKey,
		Directory:   dir,
		Filename:    path.Base(objectKey),
		ContentType: contentType,
		Size:        input.Size,
	}, nil
}

func buildObjectKey(directory, originalName, contentType string) string {
	filename := utils.GenerateULID() + models.ResolveFileExtension(originalName, contentType)
	directory = strings.TrimSuffix(directory, "/")
	return path.Join(directory, filename)
}
