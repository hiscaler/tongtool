package pui

// PackageUploadItemCaiNiao 调用菜鸟物流
type PackageUploadItemCaiNiao struct {
	TemplateURL string `json:"templateUrl"` // 模版地址，调用获取面单模版接口获取/openapi/tongtool/logi/getLabelTemplate
	Token       string `json:"token"`       // 令牌
	UserId      string `json:"userId"`      // 用户ID
}
