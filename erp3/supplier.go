package erp3

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

// 供应商列表
// https://open.tongtool.com/apiDoc.html#/?docId=992cbecb0c3b4dcd9618c24ccbdf20ab

// SupplierContacts 供应商联系人
type SupplierContacts struct {
	ContactMan      string `json:"contactMan"`  // 联系人
	Email           string `json:"email"`       // 邮箱
	FaxNumber       string `json:"faxNumber"`   // 传真号
	IsDefault       int    `json:"isDefault"`   // 是否默认（0：否、1：是）
	MobilePhone     string `json:"mobilePhone"` // 手机
	QQNumber        string `json:"qqNumber"`    // QQ 号码
	Telephone       string `json:"telephone"`   // 联系电话
	WangWangAccount string `json:"wwAccount"`   // 旺旺帐号
}

// SupplierPayment 供应商付款信息
type SupplierPayment struct {
	BankName           string `json:"bankName"`           // 支付银行/平台名称
	IsDefault          int    `json:"isDefault"`          // 是否默认（0：否、1：是）
	PaymentAccount     string `json:"paymentAccount"`     // 支付账号
	PaymentAccountName string `json:"paymentAccountName"` // 支付帐户名
	PaymentBank        string `json:"paymentBank"`        // 支付银行/平台
	PaymentType        string `json:"paymentType"`        // 支付方式（01：现金、02：银行转账、03：Paypal、04：支付宝）
	Subbranch          string `json:"subbranch"`          // 支行
}

type Supplier struct {
	Address       string            `json:"address"`       // 供应商地址
	CityId        string            `json:"cityId"`        // 城市 ID
	Contacts      SupplierContacts  `json:"contacts"`      // 联系人
	CountryId     string            `json:"countryId"`     // 国家ID
	CreatedBy     string            `json:"createdBy"`     // 创建人
	CreatedTime   string            `json:"createdTime"`   // 创建时间
	DeveloperId   string            `json:"developerId"`   // 开发人ID
	MerchantId    string            `json:"merchantId"`    // 商户编号
	Payments      []SupplierPayment `json:"payments"`      // 付款信息
	ProvinceId    string            `json:"provinceId"`    // 省份ID
	PurchaserId   string            `json:"purchaserId"`   // 采购人ID
	StoreUrl      string            `json:"storeUrl"`      // 店铺网址
	SupplierCode  string            `json:"supplierCode"`  // 供应商代码
	SupplierId    string            `json:"supplierId"`    // 供应商ID
	SupplierLevel string            `json:"supplierLevel"` // 供应商等级（01：一级、02：二级、03：三级）
	SupplierName  string            `json:"supplierName"`  // 供应商名称
	SupplierType  string            `json:"supplierType"`  // 供应商类型（01：工厂采购、02：市场采购、03：网络采购）
	UpdatedBy     string            `json:"updatedBy"`     // 更新人
	UpdatedTime   string            `json:"updatedTime"`   // 更新时间
}

type SuppliersQueryParams struct {
	MerchantId string `json:"merchantId"` // 商户编号
	Paging
}

func (m SuppliersQueryParams) Validate() error {
	return nil
}

// Suppliers 供应商列表
func (s service) Suppliers(params SuppliersQueryParams) (items []Supplier, nextToken string, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.NextToken, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = keyx.Generate(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = json.Unmarshal(b, &items); e == nil {
				return
			} else {
				s.tongTool.Logger.Printf(`cache data unmarshal error
 DATA: %s
ERROR: %s
`, string(b), e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("get cache %s error: %s", cacheKey, e.Error())
		}
	}
	res := struct {
		tongtool.Response
		Datas struct {
			NextToken string     `json:"nextToken"`
			Suppliers []Supplier `json:"suppliers"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/supplier/query")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Suppliers
			nextToken = res.Datas.NextToken
			isLastPage = nextToken == ""
		}
	} else {
		if e := json.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	if err != nil {
		return
	}

	if s.tongTool.EnableCache && len(items) > 0 {
		if b, e := json.Marshal(&items); e == nil {
			e = s.tongTool.Cache.Set(cacheKey, b)
			if e != nil {
				s.tongTool.Logger.Printf("set cache %s error: %s", cacheKey, e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("items marshal error: %s", err.Error())
		}
	}
	return
}
