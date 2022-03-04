package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

type Supplier struct {
	AccountName         string  `json:"accountName"`         // 开户名
	Bank                string  `json:"bank"`                // 开户行
	BillingCycle        float64 `json:"billingCycle"`        // 结算周期
	BillingCycleUnit    string  `json:"billingCycleUnit"`    // 结算周期单位
	CityCnName          string  `json:"cityCnName"`          // 市中文名称
	ClearingForm        string  `json:"clearingForm"`        // 结算方式（货到付款、款到发货、快递代收、定期结算）
	ClearingRemark      string  `json:"clearingRemark"`      // 结算方式备注
	CorporationFullName string  `json:"corporationFullname"` // 企业全称
	CountryCnName       string  `json:"countryCnName"`       // 国家中文名称
	Description         string  `json:"description"`         // 经营范围介绍
	DetailAddress       string  `json:"detailAddress"`       // 详细地址
	Email               string  `json:"email"`               // Email
	FaxNumber           string  `json:"faxNumber"`           // 传真号
	FullAddress         string  `json:"fullAddress"`         // 完整地址
	IsDefault           string  `json:"isDefult"`            // 是否是默认供应商
	IsDefaultBoolean    bool    `json:"isDefaultBoolean"`    // 是否是默认供应商布尔值
	Linkman             string  `json:"linkman"`             // 联系人
	PayeeAccount        string  `json:"payeeAccount"`        // 收款账号
	PaymentMode         string  `json:"paymentMode"`         // 支付方式
	PostalCode          string  `json:"postalCode"`          // 邮编
	QQNumber            string  `json:"qqNumber"`            // QQ号
	StateCnName         string  `json:"stateCnName"`         // 省/州中文名称
	SupplierCode        string  `json:"supplierCode"`        // 供应商代码
	SupplierGrade       string  `json:"supplierGrade"`       // 供应商等级
	SupplierId          string  `json:"supplierId"`          // 通途供应商id
	Telephone           string  `json:"telephone"`           // 联系电话
	WangWangNumber      string  `json:"wwNumber"`            // 旺旺号
	ZipCode             string  `json:"zipCode"`             // 电话区号
}

type SuppliersQueryParams struct {
	Paging
	MerchantId string `json:"merchantId"`
}

// Suppliers 供应商列表
// https://open.tongtool.com/apiDoc.html#/?docId=1456c221fcbf4632b06d4810e8e0d4e4
func (s service) Suppliers(params SuppliersQueryParams) (items []Supplier, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
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
	items = make([]Supplier, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []Supplier `json:"array"`
			PageNo   int        `json:"pageNo"`
			PageSize int        `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/supplierQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
			for i, item := range items {
				items[i].IsDefaultBoolean = item.IsDefault == "1"
			}
			isLastPage = len(items) < params.PageSize
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
