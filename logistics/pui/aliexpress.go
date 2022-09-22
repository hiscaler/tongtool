package pui

// PackageUploadItemAliExpress 速卖通线上物流获取token所需字段,调用速卖通线上物流必填
type PackageUploadItemAliExpress struct {
	AppKey                string `json:"appKey"`                // 速卖通开发者 APPKEY（如果不传入默认使用tongtool的开发者账号）
	AppSecret             string `json:"appSecret"`             // 速卖通开发者 APPSECRET（如果不传入默认使用tongtool的开发者账号）
	DomesticTrackingNo    string `json:"domesticTrackingNo"`    // 国内快递运单号（长度 1 ~ 32）
	ExtendData            string `json:"extendData"`            // 拣单信息[{"imageUrl":"http://xxxxxx","productDescription":"ALIBAB ALIBABA ALIBABA"}]
	IsAneroidMarkup       string `json:"isAneroidMarkup"`       // 是否含非液体化妆品(必填，填0代表不含非液体化妆品;填1代表含非液体化妆品；默认为0)
	IsContainsBattery     string `json:"isContainsBattery"`     // 是否包含锂电池(必填0/1)
	LoginId               string `json:"loginId"`               // 授权速卖通的账号信息
	PickupAddressId       int    `json:"pickupAddressId"`       // 卖家后台揽收人地址id
	RefundAddressId       int    `json:"refundAddressId"`       // 卖家后台退件人地址id
	SenderAddressId       int    `json:"senderAddressId"`       // 卖家后台寄件人地址id
	UndeliverableDecision string `json:"undeliverableDecision"` // 不可达处理(退回:0/销毁:1)
}
