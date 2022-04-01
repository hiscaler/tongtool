package pui

// PackageUploadItemChuKou1  出口易
type PackageUploadItemChuKou1 struct {
	CustomUserKey string `json:"customUserKey"` // 物流商提供账号信息
	Token         string `json:"token"`         // 物流商提供账号信息
	UserKey       string `json:"userKey"`       // 物流商提供账号信息
}

// PackageUploadItemChuKou1OutStore 出口易海外仓
type PackageUploadItemChuKou1OutStore struct {
	BuyerAccountId string `json:"buyerAccountId"` // 买家帐号ID
	CustomUserKey  string `json:"customUserKey"`  // 用户名
	Key            string `json:"key"`            // key
	Token          string `json:"token"`          // 令牌
}
