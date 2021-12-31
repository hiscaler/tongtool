package config

type Config struct {
	Debug       bool   // 是否为调试模式
	AppKey      string // 通途 APP Key
	AppSecret   string // 通途 APP Secret
	EnableCache bool   // 是否激活缓存
}
