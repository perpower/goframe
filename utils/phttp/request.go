package phttp

import (
	"github.com/go-resty/resty/v2"
)

// NewClient 构建一个client
func NewClient() *resty.Client {
	return resty.New()
}

// NewRequest 构建一个请求实例
func NewRequest() *resty.Request {
	return resty.New().R()
}
