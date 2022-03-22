package config

import (
	"github.com/faryoo/cloudrun-wechat/cache"
)

// Config for 企业微信
type Config struct {
	AppID string `json:"app_id"`

	Cache cache.Cache
}
