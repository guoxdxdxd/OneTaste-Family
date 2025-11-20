package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/services"
	"onetaste-family/backend/internal/utils"
)

// MediaHandler 媒体上传处理器
type MediaHandler struct {
	mediaService *services.MediaService
}

// NewMediaHandler 构造函数
func NewMediaHandler() *MediaHandler {
	return &MediaHandler{
		mediaService: services.NewMediaService(),
	}
}

// UploadMedia 上传文件
// @Summary 上传媒体文件
// @Description 上传图片等媒体文件，并返回可访问地址。需要 Bearer Token。
// @Tags 媒体
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "待上传文件"
// @Param path formData string true "文件目录前缀，例如 /user/head/{userId}、/family/{familyId}/dishes/{dishId}"
// @Success 200 {object} utils.Response{data=models.MediaUploadResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /media/upload [post]
func (h *MediaHandler) UploadMedia(c *gin.Context) {
	if _, ok := getUserIDFromContext(c); !ok {
		return
	}

	dir := c.PostForm("path")
	if dir == "" {
		c.JSON(http.StatusBadRequest, utils.BadRequest("请提供文件存储路径"))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequest("请选择要上传的文件"))
		return
	}

	if fileHeader.Size <= 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequest("文件内容为空"))
		return
	}

	if fileHeader.Size > services.MaxMediaFileSize {
		maxMB := services.MaxMediaFileSize / (1024 * 1024)
		c.JSON(http.StatusBadRequest, utils.BadRequest(fmt.Sprintf("文件大小不能超过 %dMB", maxMB)))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("读取文件失败"))
		return
	}
	defer file.Close()

	contentType := detectContentType(file, fileHeader.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := &models.MediaUploadInput{
		Directory:    dir,
		OriginalName: fileHeader.Filename,
		ContentType:  contentType,
		Size:         fileHeader.Size,
	}

	resp, err := h.mediaService.UploadMedia(c.Request.Context(), file, input)
	if err != nil {
		switch err {
		case models.ErrInvalidDirectory:
			c.JSON(http.StatusBadRequest, utils.BadRequest("文件存储路径不合法"))
		case services.ErrUnsupportedMediaType:
			c.JSON(http.StatusBadRequest, utils.BadRequest("仅支持上传 JPG/PNG/WebP/GIF/SVG 等常见图片格式"))
		case services.ErrInvalidMediaSize:
			c.JSON(http.StatusBadRequest, utils.BadRequest("文件大小不合法"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("上传文件失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.Success(resp))
}

func detectContentType(file io.ReadSeeker, headerType string) string {
	if headerType != "" && headerType != "application/octet-stream" {
		return headerType
	}

	var buf [512]byte
	n, _ := file.Read(buf[:])
	if n > 0 {
		file.Seek(0, io.SeekStart)
		return http.DetectContentType(buf[:n])
	}
	file.Seek(0, io.SeekStart)
	return headerType
}
