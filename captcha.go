// Copyright 2018 Kirill Zhuharev

package vk

import "context"

type CaptchaResolver interface {
	ResolveCaptcha(ctx context.Context, ssid string, imageURL string) (string, error)
}
