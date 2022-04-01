package pui

// PackageUploadItemAPAC 亚太物流获取token所需字段，调用亚太物流必填
type PackageUploadItemAPAC struct {
	Account        string `json:"account"`        // 账号
	APIAppId       string `json:"apiAppId"`       // 绑定通途的 appId
	APICertId      string `json:"apiCertId"`      // 绑定通途的 certId
	APIDevId       string `json:"apiDevId"`       // 绑定通途的 devId
	BuyerAccountId string `json:"buyerAccountId"` // 买家账号 ID
	Carrier        string `json:"carrier"`        // 承运商：CNPOST、UBI、DHL eCommerce
}
