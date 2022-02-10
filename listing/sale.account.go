package listing

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
)

// 保存店铺信息
// https://open.tongtool.com/apiDoc.html#/?docId=41e5ab4e516944208b954c4b291c8e49

// AliExpressAccount aliexpress店铺信息
type AliExpressAccount struct {
	AliexpressAccountId string `json:"aliexpressAccountId"` // 速卖通账户顺序号
	AliexpressStatus    string `json:"aliexpressStatus"`    // 速卖通账号是否有效1-有效 0-失效 2-未知
	AppKey              string `json:"appKey"`              // APP KEY
	AppSecret           string `json:"appSecret"`           // 密钥
	CreatedBy           string `json:"createdBy"`           // 创建人
	CreatedDate         string `json:"createdDate"`         // 创建时间
	LoginId             string `json:"loginId"`             // 速卖通账号唯一标识
	MerchantId          string `json:"merchantId"`          // 商户ID
	RefreshToken        string `json:"refreshToken"`        // 授权token
	RefreshTokenTime    string `json:"refreshTokenTime"`    // 过期时间
	SaleAccountId       string `json:"saleAccountId"`       // 店铺账户编号
	SellerAliId         string `json:"sellerAliId"`         // sellerAliId(用于获取access token)
	UpdatedBy           string `json:"updatedBy"`           // 修改人
	UpdatedDate         string `json:"updatedDate"`         // 修改时间
}

// AmazonAccount Amazon店铺信息
type AmazonAccount struct {
	AccessKeyId      string `json:"accessKeyId"`      // Amazon Access Key ID
	AmazonAccountId  string `json:"amazonAccountId"`  // amazon账号主键Id
	AmazonMerchantId string `json:"amazonMerchantId"` // Amazon商户号（操作amazon账号时必填）
	CreatedBy        string `json:"createdBy"`        // 创建人
	CreatedDate      string `json:"createdDate"`      // 创建时间
	DeveloperCode    string `json:"developerCode"`    // 开发者账号
	MarketplaceIds   string `json:"marketplaceIds"`   // 站点ID集合（marketplaceId1,marketplaceId2……多个站点Id用英文逗号分隔，操作amazon账号时必填）
	MerchantId       string `json:"merchantId"`       // 商户ID
	SaleAccountId    string `json:"saleAccountId"`    // 店铺账号ID
	SecretKey        string `json:"secretKey"`        // Amazon secret key
	Status           string `json:"status"`           // 是否有效
	Token            string `json:"token"`            // 授权token
	UpdatedBy        string `json:"updatedBy"`        // 修改人
	UpdatedDate      string `json:"updatedDate"`      // 修改时间
}

// WishAccount wish店铺信息
type WishAccount struct {
	ClientId              string `json:"clientId"`              // clientId
	ClientSecret          string `json:"clientSecret"`          // clientSecret
	CreatedBy             string `json:"createdBy"`             // 创建人
	CreatedDate           string `json:"createdDate"`           // 创建时间
	LocalizedCurrencyCode string `json:"localizedCurrencyCode"` // 本地货币
	MerchantId            string `json:"merchantId"`            // 商户ID
	MerchantUserId        string `json:"merchantUserId"`        // merchantUserId
	RefreshToken          string `json:"refreshToken"`          // 授权token（操作wish账号时必填）
	SaleAccountId         string `json:"saleAccountId"`         // 账号顺序号
	UpdatedBy             string `json:"updatedBy"`             // 修改人
	UpdatedDate           string `json:"updatedDate"`           // 修改时间
	WishAccountId         string `json:"wishAccountId"`         // wish账户顺序号
	WishStatus            string `json:"wishStatus"`            // wish账号是否有效1-有效 0-失效 2-未知
}

// EbayAccount ebay店铺信息
type EbayAccount struct {
	AccessToken             string `json:"accessToken"`             // ebay Access Tokens
	AccessTokenInvalidDate  string `json:"accessTokenInvalidDate"`  // accessToken失效日期
	AppKey                  string `json:"appKey"`                  // APP KEY
	AppSecret               string `json:"appSecret"`               // 密钥
	ApplyDate               string `json:"applyDate"`               // 大中华卖家申请时间
	ApplyStatus             string `json:"applyStatus"`             // 大中华卖家申请状态; 0：未申请，1：申请成功，2：申请中，3：申请失败
	CreatedBy               string `json:"createdBy"`               // 创建人
	CreatedDate             string `json:"createdDate"`             // 创建时间
	DeveloperId             string `json:"developerId"`             // 开发者ID
	EbayAccountId           string `json:"ebayAccountId"`           // ebay店铺账号主键ID
	EIASToken               string `json:"eiasToken"`               // ebay用户唯一标识EIASToken（操作amazon账号时必填）
	IsAboutToExpire         string `json:"isAboutToExpire"`         // 密钥是否即将过期(1-即将过期) 单选：0,1
	MerchantId              string `json:"merchantId"`              // 商户ID
	NoticeState             string `json:"noticeState"`             // 通知开启状态，0：未开启，1：开启
	RefreshToken            string `json:"refreshToken"`            // ebay Access Tokens
	RefreshTokenInvalidDate string `json:"refreshTokenInvalidDate"` // refreshToken失效日期
	RegisterDate            string `json:"registerDate"`            // ebay账号注册时间
	SaleAccountId           string `json:"saleAccountId"`           // 店铺账号ID
	StoreSite               string `json:"storeSite"`               // eBay店铺站点
	Token                   string `json:"token"`                   // 安全证书(Token)
	TokenEffectiveDate      string `json:"tokenEffectiveDate"`      // 密钥生效日期
	TokenInvalidDate        string `json:"tokenInvalidDate"`        // 密钥失效日期
	UpdatedBy               string `json:"updatedBy"`               // 修改人
	UpdatedDate             string `json:"updatedDate"`             // 修改时间
}

type UpsertSaleAccountRequest struct {
	Account               string            `json:"account,omitempty"`               // 账户名称(平台账户)
	AliexpressAccountInfo AliExpressAccount `json:"aliexpressAccountInfo,omitempty"` // aliexpress店铺信息
	AmazonAccountInfo     AmazonAccount     `json:"amazonAccountInfo,omitempty"`     // Amazon店铺信息
	WishAccountInfo       WishAccount       `json:"wishAccountInfo,omitempty"`       // wish店铺信息
	EbayAccountInfo       EbayAccount       `json:"ebayAccountInfo,omitempty"`       // ebay店铺信息
	AccountCode           string            `json:"accountCode"`                     // 账户简码(自定义账户)
	CreatedBy             string            `json:"createdBy,omitempty"`             // 创建人
	CreatedDate           string            `json:"createdDate,omitempty"`           // 创建时间
	DisableTime           string            `json:"disableTime,omitempty"`           // 停用时间
	EnableTime            string            `json:"enableTime,omitempty"`            // 启用时间
	MerchantId            string            `json:"merchantId"`                      // 商户编号
	OutOfStock            string            `json:"outOfStock,omitempty"`            // ebay零库存是否不下架,Y-不下架,N-下架
	PlatformId            string            `json:"platformId"`                      // 平台编号（amazon,ebay,aliexpress,wish其中一个）
	SaleAccountId         string            `json:"saleAccountId,omitempty"`         // 店铺账户编号
	Status                string            `json:"status,omitempty"`                // 账户状态(停用0，启用1)
	StoreName             string            `json:"storeName,omitempty"`             // 店铺名称
	UpdatedBy             string            `json:"updatedBy,omitempty"`             // 修改人
	UpdatedDate           string            `json:"updatedDate,omitempty"`           // 修改时间
	UserId                string            `json:"userId,omitempty"`                // 用户Id
}

func (m UpsertSaleAccountRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.AccountCode, validation.Required.Error("账户简码不能为空")),
		validation.Field(&m.PlatformId, validation.Required.Error("平台编号不能为空"), validation.In("amazon", "ebay", "aliexpress", "wish").Error("无效的平台简码")),
	)
}

// UpsertSaleAccount 保存店铺信息
func (s service) UpsertSaleAccount(req UpsertSaleAccountRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/saleAccount/saveSaleAccount")
	if err == nil {
		if resp.IsSuccess() {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	return err
}
