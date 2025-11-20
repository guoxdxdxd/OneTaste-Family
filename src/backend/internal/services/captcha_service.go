package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/utils"
	"onetaste-family/backend/pkg/cache"
)

const (
	captchaKeyPrefix = "captcha:"
	captchaTTL       = 30 * time.Second
	captchaCharset   = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz"

	defaultCaptchaWidth  = 220
	defaultCaptchaHeight = 70
	minCaptchaWidth      = 120
	maxCaptchaWidth      = 360
	minCaptchaHeight     = 40
	maxCaptchaHeight     = 160
	minNoiseCount        = 25
	maxNoiseCount        = 220
)

const (
	defaultNoiseCount = 140
	defaultLineOption = base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSineLine
)

// CaptchaService 图形验证码服务
type CaptchaService struct {
	redis *redis.Client
	fonts []string
}

// NewCaptchaService 创建验证码服务
func NewCaptchaService() *CaptchaService {
	return &CaptchaService{
		redis: cache.GetRedis(),
		fonts: []string{"wqy-microhei.ttc", "Flim-Flam.ttf", "RitaSmith.ttf"},
	}
}

// GenerateCaptcha 生成图形验证码
func (s *CaptchaService) GenerateCaptcha(ctx context.Context, req *models.CaptchaRequest) (*models.CaptchaResponse, error) {
	if s.redis == nil {
		return nil, errors.New("redis client is not initialized")
	}

	cfg := s.buildDriverConfig(req)
	driver := base64Captcha.NewDriverString(
		cfg.height,
		cfg.width,
		cfg.noiseCount,
		cfg.lineOption,
		4,
		captchaCharset,
		nil,
		nil,
		s.fonts,
	)

	_, content, answer := driver.GenerateIdQuestionAnswer()
	item, err := driver.DrawCaptcha(content)
	if err != nil {
		return nil, fmt.Errorf("failed to draw captcha: %w", err)
	}

	key := utils.GenerateULID()
	if err := s.redis.Set(ctx, s.redisKey(key), strings.ToUpper(answer), captchaTTL).Err(); err != nil {
		return nil, fmt.Errorf("failed to store captcha: %w", err)
	}

	return &models.CaptchaResponse{
		CaptchaKey:  key,
		ImageBase64: item.EncodeB64string(),
		ExpireIn:    int64(captchaTTL.Seconds()),
	}, nil
}

// ValidateCaptcha 验证图形验证码
func (s *CaptchaService) ValidateCaptcha(ctx context.Context, key, value string) error {
	if s.redis == nil {
		return errors.New("redis client is not initialized")
	}

	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if key == "" || value == "" {
		return ErrInvalidVerifyCode
	}

	redisKey := s.redisKey(key)
	stored, err := s.redis.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		return ErrCaptchaExpired
	}
	if err != nil {
		return fmt.Errorf("failed to get captcha: %w", err)
	}

	if !strings.EqualFold(stored, value) {
		return ErrInvalidVerifyCode
	}

	s.redis.Del(ctx, redisKey)
	return nil
}

func (s *CaptchaService) redisKey(key string) string {
	return captchaKeyPrefix + key
}

type driverConfig struct {
	width      int
	height     int
	noiseCount int
	lineOption int
}

func (s *CaptchaService) buildDriverConfig(req *models.CaptchaRequest) driverConfig {
	width := defaultCaptchaWidth
	height := defaultCaptchaHeight

	if req != nil {
		if req.Width != 0 {
			width = req.Width
		}
		if req.Height != 0 {
			height = req.Height
		}
	}

	width = clampDimension(width, minCaptchaWidth, maxCaptchaWidth, defaultCaptchaWidth)
	height = clampDimension(height, minCaptchaHeight, maxCaptchaHeight, defaultCaptchaHeight)

	baseArea := defaultCaptchaWidth * defaultCaptchaHeight
	area := width * height
	noise := defaultNoiseCount * area / baseArea
	if noise < minNoiseCount {
		noise = minNoiseCount
	}
	if noise > maxNoiseCount {
		noise = maxNoiseCount
	}

	return driverConfig{
		width:      width,
		height:     height,
		noiseCount: noise,
		lineOption: defaultLineOption,
	}
}

func clampDimension(value, min, max, def int) int {
	if value == 0 {
		return def
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
