package models

import (
	"errors"
	"path"
	"path/filepath"
	"strings"
)

var mediaMimeExtensions = map[string]string{
	"image/jpeg":    "jpg",
	"image/png":     "png",
	"image/webp":    "webp",
	"image/gif":     "gif",
	"image/svg+xml": "svg",
}

var allowedRootDirs = map[string]struct{}{
	"user":    {},
	"family":  {},
	"dishes":  {},
	"dish":    {},
	"general": {},
}

// MediaUploadInput 上传文件时的业务输入
type MediaUploadInput struct {
	Directory    string
	OriginalName string
	ContentType  string
	Size         int64
}

// MediaUploadResponse 上传成功返回的数据
type MediaUploadResponse struct {
	URL         string `json:"url"`
	Bucket      string `json:"bucket"`
	ObjectKey   string `json:"object_key"`
	Directory   string `json:"directory"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// IsAllowedMediaContentType 判断 MIME 是否允许
func IsAllowedMediaContentType(contentType string) bool {
	ct := normalizeContentType(contentType)
	_, ok := mediaMimeExtensions[ct]
	return ok
}

// ResolveFileExtension 依据文件名或 MIME 得到扩展名，默认返回 ".bin"
func ResolveFileExtension(originalName, contentType string) string {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(originalName)))
	if ext != "" {
		return ext
	}

	ct := normalizeContentType(contentType)
	if suffix, ok := mediaMimeExtensions[ct]; ok {
		return "." + suffix
	}

	return ".bin"
}

// CleanStorageDirectory 规范化用户传入的目录
func CleanStorageDirectory(dir string) (string, error) {
	if dir == "" {
		return "", ErrInvalidDirectory
	}

	cleaned := path.Clean("/" + dir)
	cleaned = strings.TrimPrefix(cleaned, "/")
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" || cleaned == "." {
		return "", ErrInvalidDirectory
	}
	if strings.HasPrefix(cleaned, "..") {
		return "", ErrInvalidDirectory
	}

	root := strings.Split(cleaned, "/")[0]
	if _, ok := allowedRootDirs[root]; !ok {
		return "", ErrInvalidDirectory
	}

	return cleaned, nil
}

var ErrInvalidDirectory = errors.New("invalid storage directory")

func normalizeContentType(contentType string) string {
	if contentType == "" {
		return ""
	}
	parts := strings.Split(contentType, ";")
	return strings.ToLower(strings.TrimSpace(parts[0]))
}
