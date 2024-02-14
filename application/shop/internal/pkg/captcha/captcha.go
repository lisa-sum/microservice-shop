package captcha

import (
	"context"
	"fmt"
	"github.com/mojocn/base64Captcha"
)

// 验证码文件

var Store = base64Captcha.DefaultMemStore

type CapInfo struct {
	CaptchaId string
	PicPath   string
}

// GetCaptcha 生成验证码
func GetCaptcha(ctx context.Context) (*CapInfo, error) {
	driver := base64Captcha.NewDriverDigit(80, 250, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, Store)
	id, b64s, answer, err := cp.Generate()
	if err != nil {
		return nil, err
	}
	fmt.Println(id, b64s, answer)

	return &CapInfo{
		CaptchaId: id,
		PicPath:   b64s,
	}, nil
}
