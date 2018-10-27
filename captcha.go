// Copyright 2018 Kirill Zhuharev

package vk

type CaptchaResolver interface {
	ResolveCaptcha(ssid string, imageURL string) (string, error)
}
