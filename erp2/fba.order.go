package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
)

// FBAOrder 通途 FBA 订单
type FBAOrder struct {
	BuyerEmail         string  `json:"buyerEmail"`         // 买家邮箱
	BuyerName          string  `json:"buyerName"`          // 买家姓名
	BuyerPhoneNumber   string  `json:"buyerPhoneNumber"`   // 买家电话
	Currency           string  `json:"currency"`           // 币种
	OrderId            string  `json:"orderId"`            // 订单号
	PageNo             int     `json:"pageNo"`             // 查询页数
	PageSize           int     `json:"pageSize"`           // 查询数量
	PaymentsDate       string  `json:"paymentsDate"`       // 付款时间
	PurchaseDate       int     `json:"purchaseDate"`       // 购买时间
	RecipientName      string  `json:"recipientName"`      // 收件人姓名
	SalesChannel       string  `json:"salesChannel"`       // 销售站点
	ShipAddress1       string  `json:"shipAddress1"`       // 地址1
	ShipAddress2       string  `json:"shipAddress2"`       // 地址2
	ShipAddress3       string  `json:"shipAddress3"`       // 地址3
	ShipCity           string  `json:"shipCity"`           // 城市
	ShipCountry        string  `json:"shipCountry"`        // 国家
	ShipPhoneNumber    string  `json:"shipPhoneNumber"`    // 收件人电话
	ShipPostalCode     string  `json:"shipPostalCode"`     // 邮编
	ShipServiceLevel   string  `json:"shipServiceLevel"`   // 物流服务等级
	ShipState          string  `json:"shipState"`          // 州/省
	TotalItemPrice     float64 `json:"totalItemPrice"`     // 货品总计
	TotalItemTax       float64 `json:"totalItemTax"`       // 商品税费总计
	TotalShippingPrice float64 `json:"totalShippingPrice"` // 物流费用总计
	TotalShippingTax   float64 `json:"totalShippingTax"`   // 物流税费总计
}

type FBAOrdersQueryParams struct {
	Paging
	Account          string `json:"account,omitempty"`          // 速卖通登录账号
	MerchantId       string `json:"merchantId"`                 // 商户ID
	PurchaseDateFrom string `json:"purchaseDateFrom,omitempty"` // 订单购买时间开始时间
	PurchaseDateTo   string `json:"purchaseDateTo,omitempty"`   // 订单购买时间结束时间
}

func (m FBAOrdersQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PurchaseDateFrom, validation.Required.Error("订单购买开始时间不能为空"), validation.Date(constant.DatetimeFormat).Error("无效的订单购买开始时间格式")),
		validation.Field(&m.PurchaseDateTo, validation.Required.Error("订单购买结束时间不能为空"), validation.Date(constant.DatetimeFormat).Error("无效的订单购买结束时间格式")),
	)
}

// FBAOrders FBA 订单列表
// https://open.tongtool.com/apiDoc.html#/?docId=c33e7bd4e73d4d2d9a27de56f794cc82
func (s service) FBAOrders(params FBAOrdersQueryParams) (items []FBAOrder, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

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
	items = make([]FBAOrder, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array []FBAOrder `json:"array"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/fbaOrderQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
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
