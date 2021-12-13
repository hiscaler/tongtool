package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
)

// FBAOrder 通途 FBA 订单
type FBAOrder struct {
	BuyerEmail         string  `json:"buyerEmail"`
	BuyerName          string  `json:"buyerName"`
	BuyerPhoneNumber   string  `json:"buyerPhoneNumber"`
	Currency           string  `json:"currency"`
	OrderId            string  `json:"orderId"`
	PageNo             int     `json:"pageNo"`
	PageSize           int     `json:"pageSize"`
	PaymentsDate       string  `json:"paymentsDate"`
	PurchaseDate       int     `json:"purchaseDate"`
	RecipientName      string  `json:"recipientName"`
	SalesChannel       string  `json:"salesChannel"`
	ShipAddress2       string  `json:"shipAddress2"`
	ShipAddress3       string  `json:"shipAddress3"`
	ShipCity           string  `json:"shipCity"`
	ShipCountry        string  `json:"shipCountry"`
	ShipPhoneNumber    string  `json:"shipPhoneNumber"`
	ShipPostalCode     string  `json:"shipPostalCode"`
	ShipServiceLevel   string  `json:"shipServiceLevel"`
	ShipState          string  `json:"shipState"`
	TotalItemPrice     float64 `json:"totalItemPrice"`
	TotalItemTax       string  `json:"totalItemTax"`
	TotalShippingPrice string  `json:"totalShippingPrice"`
	TotalShippingTax   string  `json:"totalShippingTax"`
}

type FBAOrderQueryParams struct {
	Account          string `json:"account"`
	MerchantId       string `json:"merchantId"`
	PageNo           int    `json:"pageNo,omitempty"`
	PageSize         int    `json:"pageSize,omitempty"`
	PurchaseDateFrom string `json:"purchaseDateFrom,omitempty"`
	PurchaseDateTo   string `json:"purchaseDateTo,omitempty"`
}

type fbaOrderResult struct {
	result
	Datas struct {
		Array []FBAOrder `json:"array"`
	}
}

func (s service) FBAOrders(params FBAOrderQueryParams) (items []FBAOrder, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = s.tongTool.QueryDefaultValues.PageNo
	}
	if params.PageSize <= 0 {
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
			if err = tongtool.HasError(res.Code); err == nil {
				items = res.Datas.Array
				isLastPage = len(items) < params.PageSize
			}
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}
