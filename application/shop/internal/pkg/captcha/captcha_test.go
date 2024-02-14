package captcha

import (
	"context"
	"testing"
)

func TestGetCaptcha(t *testing.T) {
	captcha, err := GetCaptcha(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(captcha.CaptchaId)
	t.Log(captcha.PicPath)
}
