package config

type Config struct {
	Debug            bool   // 是否为调试模式
	Timeout          int    // 请求超市时间（秒）
	RetryCount       int    // 重试次数
	RetryWaitTime    int    // 重试等待时间
	RetryMaxWaitTime int    // 重试最大等待时间
	ForceWaiting     bool   // 是否强制等待
	AppKey           string // 通途 APP Key
	AppSecret        string // 通途 APP Secret
	EnableCache      bool   // 是否激活缓存
}
