package logistics

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
	"strconv"
	"strings"
)

// 上传包裹信息
// https://open.tongtool.com/apiDoc.html#/?docId=f4efa880c11242c1a69bf8db4d763fec

// PackageUploadItem 批量包裹信息
type PackageUploadItem struct {
	ACS PackageUploadItemACS `json:"acs"` // ACS物流
	AJ  PackageUploadItemAJ  `json:"aj"`  // AnJun
}

// PackageUplodItemACSGoods 包裹物品明細
type PackageUplodItemACSGoods struct {
	BarCode  string  `json:"barCode"`  // 商品條形碼
	Brand    string  `json:"brand"`    // 品牌
	Code     string  `json:"code"`     // 商品貨號
	Country  string  `json:"country"`  // 原產國
	Currency string  `json:"currency"` // 币种，比如：RMB,USD,AUD,EUR,KRW,JPY,THB,GBP
	HsCode   string  `json:"hsCode"`   // HS編碼
	Name     string  `json:"name"`     // 品名
	Num      int     `json:"num"`      // 商品數量
	Spec     string  `json:"spec"`     // 規格型號，如：120粒/瓶
	TaxNo    string  `json:"taxNo"`    // 行郵稅號
	Unit     string  `json:"unit"`     // 計量單位：如瓶，个，件等
	Value    float64 `json:"value"`    // 單價，單位：元
}

type PackageUploadItemACS struct {
	CodCurrency string                     `json:"codCurrency"` // 代收款项 币种
	CodMoney    float64                    `json:"codMoney"`    // 代收款项 金额
	Cuscode     string                     `json:"cuscode"`     // 客户编码
	GoodsList   []PackageUplodItemACSGoods `json:"goodsList"`   // 包裹物品明細
	Key         string                     `json:"key"`         // 令牌,非通途用户必填
	SiteCode    string                     `json:"siteCode"`    // TMS的倉庫代碼
	Transport   string                     `json:"transport"`   // 運輸方式，值為：空運，陸運，海運，海快
}

// PackageUploadItemAJ AnJun
type PackageUploadItemAJ struct {
	ApiUser       string `json:"apiUser"`       // 物流商提供的客户名
	Battery       int    `json:"battery"`       // 是否带电（ 0：不带电、1：带电）
	DeclareRemark string `json:"declareRemark"` // 配货信息，安骏(Title2)限制不超过126字符
	Token         string `json:"token"`         // 物流商提供的 Token
	WishAccountId string `json:"wishAccountId"` // wish 订单所在对应的 wish account ID
}

// SenderAddress 寄件人信息
type SenderAddress struct {
	SenderCity          string `json:"senderCity"`          // 寄件人城市
	SenderCompany       string `json:"senderCompany"`       // 寄件人公司名称,部分物流商必填
	SenderCountry       string `json:"senderCountry"`       // 寄件人国家
	SenderCountryEnName string `json:"senderCountryEnName"` // 寄件人国家名称(英文名),部分物流商必填
	SenderDistrict      string `json:"senderDistrict"`      // 寄件人区域,部分物流商必填
	SenderEmail         string `json:"senderEmail"`         // 寄件人邮箱,部分物流商必填
	SenderName          string `json:"senderName"`          // 寄件人姓名
	SenderPhone         string `json:"senderPhone"`         // 寄件人电话
	SenderPostcode      string `json:"senderPostcode"`      // 寄件人邮编
	SenderState         string `json:"senderState"`         // 寄件人省州
	SenderStreet        string `json:"senderStreet"`        // 寄件人地址
}

// ReturnAddress 退件人信息
type ReturnAddress struct {
	ReturnCity     string `json:"returnCity"`     // 退件人城市
	ReturnCompany  string `json:"returnCompany"`  // 揽收人公司,部分物流商必填
	ReturnCountry  string `json:"returnCountry"`  // 退件人国家
	ReturnDistrict string `json:"returnDistrict"` // 揽收人区域,部分物流商必填
	ReturnName     string `json:"returnName"`     // 退件人姓名
	ReturnPhone    string `json:"returnPhone"`    // 退件人电话
	ReturnPostcode string `json:"returnPostcode"` // 退件人邮编
	ReturnState    string `json:"returnState"`    // 退件人省州
	ReturnStreet   string `json:"returnStreet"`   // 退件人地址
}

// PickupAddress 揽收人信息
type PickupAddress struct {
	PickupCity          string `json:"pickupCity"`          // 揽收人城市
	PickupCompany       string `json:"pickupCompany"`       // 揽收人公司名称,部分物流商必填
	PickupCountry       string `json:"pickupCountry"`       // 揽收人国家
	PickupCountryEnName string `json:"pickupCountryEnName"` // 揽件人国名称(英文名),部分物流商必填
	PickupDistrict      string `json:"pickupDistrict"`      // 揽收人区域,部分物流商必填
	PickupEmail         string `json:"pickupEmail"`         // 揽收人邮箱,部分物流商必填
	PickupName          string `json:"pickupName"`          // 揽收人姓名
	PickupPhone         string `json:"pickupPhone"`         // 揽收人电话
	PickupPostcode      string `json:"pickupPostcode"`      // 揽收人邮编
	PickupState         string `json:"pickupState"`         // 揽收人省州
	PickupStreet        string `json:"pickupStreet"`        // 揽收人地址
}

type PackageUploadRequest struct {
	CarrierCode   string              `json:"carrierCode"`             // 物流商代码
	PackageItems  []PackageUploadItem `json:"packageItems"`            // 上传包裹信息
	SenderAddress SenderAddress       `json:"senderAddress,omitempty"` //	寄件人信息
	ReturnAddress ReturnAddress       `json:"returnAddress,omitempty"` // 退件人信息
	PickupAddress PickupAddress       `json:"pickupAddress,omitempty"` // 揽收人信息
	MerchantId    string              `json:"merchantId"`              // 商户号
	Source        string              `json:"source"`                  // 订单来源
	WarehouseCode string              `json:"warehouseCode"`           // 仓库代码，海外仓独有属性,不要随意传值!!
}

func (m PackageUploadRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CarrierCode, validation.Required.Error("物流商代码不能为空")),
		validation.Field(&m.PackageItems,
			validation.Required.Error("包裹信息不能为空"),
		),
		validation.Field(&m.Source, validation.Required.Error("订单来源不能为空")),
	)
}

// PackageUpload
// todo
func (s service) PackageUpload(req PackageUploadRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	res := struct {
		tongtool.Response
		Datas struct {
			ACK          string `json:"ack"`          // 响应结果（Success：成功、Failure：失败）
			ErrorCode    string `json:"code"`         // 错误代码
			ErrorMessage string `json:"errorMessage"` // 错误信息
		} `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/logi/packageUpload")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			if strings.EqualFold(res.Datas.ACK, "Failure") {
				errorCode, _ := strconv.Atoi(res.Datas.ErrorCode)
				err = tongtool.ErrorWrap(errorCode, res.Datas.ErrorMessage)
			}
		}
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}
