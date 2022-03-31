package pui

// PackageUploadItemAmazonOnlineGoodsItemDeclare 货品关联申报信息
type PackageUploadItemAmazonOnlineGoodsItemDeclare struct {
	Currency      string  `json:"currency"`      // 申报币种
	DeclareCnName string  `json:"declareCnName"` // 申报中文名
	DeclareEnName string  `json:"declareEnName"` // 申报英文名
	HsCode        string  `json:"hsCode"`        // 海关编码
	UnitPrice     float64 `json:"unitPrice"`     // 申报价格
	Weight        float64 `json:"weight"`        // 重量
}

// PackageUploadItemAmazonOnlineGoodsItem 货品
type PackageUploadItemAmazonOnlineGoodsItem struct {
	DeclareList    []PackageUploadItemAmazonOnlineGoodsItemDeclare `json:"declareList"`    // 货品关联申报信息,必填
	Quantity       int                                             `json:"quantity"`       // 货品数量
	WebStoreItemId string                                          `json:"webstoreItemId"` // 订单明细ID
}

// PackageUploadItemAmazonOnline AmazonOnline
type PackageUploadItemAmazonOnline struct {
	AccessKeyId         string                                   `json:"accessKeyId"`         // accessKeyId
	AccessToken         string                                   `json:"accessToken"`         // 授权后的accessToken
	AmazonUndeliverable string                                   `json:"amazonUndeliverable"` // 无法投递处理方案:ABANDON-丢弃;RETURN_TO_SELLER-退回
	DeliveryExperience  string                                   `json:"deliveryExperience"`  // 交货方式(only MWS):DeliveryConfirmationWithAdultSignature;DeliveryConfirmationWithSignature;DeliveryConfirmationWithoutSignature;NoTracking
	DeliveryMode        string                                   `json:"deliveryMode"`        // 交运方式(only MWS):0-上门揽收;1-卖家自送
	GoodsItemList       []PackageUploadItemAmazonOnlineGoodsItem `json:"goodsItemList"`       // 货品信息
	IAMAccessKey        string                                   `json:"iamAccessKey"`        // AWS IAM key(签名需要必填)
	IAMSecretKey        string                                   `json:"iamSecretKey"`        // AWS IAM 秘钥(签名需要必填)
	MarketplaceId       string                                   `json:"marketplaceId"`       // marketplaceId(only MWS)
	MWSToken            string                                   `json:"mwsToken"`            // mwsToken(only MWS)
	RefreshToken        string                                   `json:"refreshToken"`        // 刷新token,传了accessToken此参数就不需要传值
	Region              string                                   `json:"region"`              // 区域，例如：us-east-1
	SecretAccessKey     string                                   `json:"secretAccessKey"`     // secretAccessKey
	SellerId            string                                   `json:"sellerId"`            // sellerId(only MWS)
	ServiceURL          string                                   `json:"serviceURL"`          // serviceURL
	Version             string                                   `json:"version"`             // 亚马逊接口版本：MWS,SPAPI 默认：MWS
}
