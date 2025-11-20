package models

// CaptchaRequest 图形验证码请求参数
type CaptchaRequest struct {
	Width  int `form:"width" example:"220"` // 图片宽度，单位像素
	Height int `form:"height" example:"70"` // 图片高度，单位像素
}

// CaptchaResponse 图形验证码响应
// @Description 图形验证码返回数据
type CaptchaResponse struct {
	CaptchaKey  string `json:"captcha_key" example:"01J0XYZABCD1234EFG567HIJK"`                        // 验证码编码
	ImageBase64 string `json:"image_base64" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUg..."` // 图形验证码Base64
	ExpireIn    int64  `json:"expire_in" example:"30"`                                                 // 过期时间（秒）
}
