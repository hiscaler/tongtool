package pui

// PackageUploadItemAmloutStore 艾姆勒俄罗斯仓
type PackageUploadItemAmloutStore struct {
	AppKey         string `json:"appKey"`         // APP KEY
	AppToken       string `json:"appToken"`       // 令牌
	Doorplate      string `json:"doorplate"`      // 门牌号
	ForceVerify    int    `json:"forceVerify"`    // 是否强制校验
	IsSignature    int    `json:"isSignature"`    // 是否需要签名
	Platform       string `json:"platform"`       // 平台类型
	Service        string `json:"service"`        // 接口名称
	TransactionId  string `json:"transactionId"`  // 交易号
	WebStoreItemId string `json:"webstoreItemId"` // 货品编号
}
