package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
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
	TotalItemTax       string  `json:"totalItemTax"`       // 商品税费总计
	TotalShippingPrice string  `json:"totalShippingPrice"` // 物流费用总计
	TotalShippingTax   string  `json:"totalShippingTax"`   // 物流税费总计
}

type FBAOrderQueryParams struct {
	Account          string `json:"account,omitempty"`          // 速卖通登录账号
	MerchantId       string `json:"merchantId"`                 // 商户ID
	PageNo           int    `json:"pageNo,omitempty"`           // 查询页数
	PageSize         int    `json:"pageSize,omitempty"`         // 每页数量
	PurchaseDateFrom string `json:"purchaseDateFrom,omitempty"` // 订单购买时间开始时间
	PurchaseDateTo   string `json:"purchaseDateTo,omitempty"`   // 订单购买时间结束时间
}

type fbaOrderResult struct {
	result
	Datas struct {
		Array []FBAOrder `json:"array"`
	}
}

// FBAOrders FBA 订单列表
// https://open.tongtool.com/apiDoc.html#/?docId=c33e7bd4e73d4d2d9a27de56f794cc82
func (s service) FBAOrders(params FBAOrderQueryParams) (items []FBAOrder, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]FBAOrder, 0)
	res := fbaOrderResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/fbaOrderQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				isLastPage = len(items) < params.PageSize
			}
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}
