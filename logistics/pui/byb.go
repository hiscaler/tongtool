package pui

// PackageUploadItemByb 万色物流获取token所需字段,调用万色物流必填
type PackageUploadItemByb struct {
	AppKey string `json:"appKey"` // API密钥
	Otype  string `json:"otype"`  // 是否挂号件(0=否，1=是)
	TypeNo string `json:"typeNo"` // 内件类型代码(1=礼品,2=文件,3=商品货样,4=其他)
	// 如传分仓代码必须在贝邮宝或WISH邮网站配置揽收地址信息并设置默认）贝邮宝：(传空值为北京仓)
	// 0=北京仓
	// 1=上海仓
	// 2=广州仓
	// 3=深圳仓
	// 4=义乌仓
	// 5=南京仓
	// 6=福州仓;
	// WISH邮：(传空值为上海仓)1=上海仓2=广州仓3=深圳仓4=义乌仓5=北京仓6=南京仓7=福州仓
	WareCode int `json:"wareCode"` // 分仓代码
}
