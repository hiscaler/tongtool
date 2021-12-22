package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
	"strconv"
)

type ShopifyOrderItem struct {
	ItemId        string  `json:"itemId"`        // 订单交易号
	Price         string  `json:"price"`         // 价格
	PriceValue    float64 `json:"priceValue"`    // 价格
	ProductId     string  `json:"productId"`     // 产品ID
	Quantity      string  `json:"quantity"`      // 产品数量
	QuantityValue int     `json:"QuantityValue"` // 产品数量
	SKU           string  `json:"sku"`           // SKU
	Title         string  `json:"title"`         // 产品名称
	Weight        string  `json:"weight"`        // 重量
	WeightValue   float64 `json:"weightValue"`   // 重量
}

type ShopifyOrder struct {
	Cod                 string             `json:"cod"`                 // 是否是货到付款订单
	CodBoolean          bool               `json:"codBoolean"`          // 是否是货到付款订单布尔值
	FinancialStatus     string             `json:"financialStatus"`     // 支付状态：pending-未付款,paid-已付款,partially_paid-部分付款
	Gateway             string             `json:"gateway"`             // 网关
	Items               []ShopifyOrderItem `json:"items"`               // 订单明细
	OrderName           string             `json:"order_name"`          // 订单名称
	OrderStatus         string             `json:"orderStatus"`         // 订单状态
	PaymentTime         string             `json:"paymentTime"`         // 付款时间
	PaypalAccount       string             `json:"paypalAccount"`       // Paypal账号
	PaypalTransactionId string             `json:"paypalTransactionId"` // Paypal交易号
	SalesOrderNumber    string             `json:"salesOrderNumber"`    // 销售单号
	ShopifyOrderId      string             `json:"shopifOrderId"`       // Shopify订单号
}

type ShopifyOrderQueryParams struct {
	BuyerEmail          string `json:"buyerEmail,omitempty"`
	MerchantId          string `json:"merchantId"`
	PageNo              int    `json:"pageNo,omitempty"`
	PageSize            int    `json:"pageSize,omitempty"`
	PayDateFrom         string `json:"payDateFrom,omitempty"`
	PayDateTo           string `json:"payDateTo,omitempty"`
	PaypalTransactionId string `json:"paypalTransactionId,omitempty"`
	ShopifyOrderId      string `json:"shopifyOrderId,omitempty"`
}

type shopifyOrderResult struct {
	result
	Datas struct {
		Array    []ShopifyOrder `json:"array"`
		PageNo   int            `json:"pageNo"`
		PageSize int            `json:"pageSize"`
	}
}

// ShopifyOrders Shopify 订单列表
// https://open.tongtool.com/apiDoc.html#/?docId=e949a88561e7471785cccef86feb3e6d
func (s service) ShopifyOrders(params ShopifyOrderQueryParams) (items []ShopifyOrder, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]ShopifyOrder, 0)
	res := shopifyOrderResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/shopifyOrderQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				for i, item := range items {
					items[i].CodBoolean = item.Cod == "true"
					for j, orderItem := range item.Items {
						orderItem.PriceValue, _ = strconv.ParseFloat(orderItem.Price, 64)
						orderItem.QuantityValue, _ = strconv.Atoi(orderItem.Quantity)
						orderItem.WeightValue, _ = strconv.ParseFloat(orderItem.Weight, 64)
						items[i].Items[j] = orderItem
					}
				}
				isLastPage = len(items) < params.PageSize
			}
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}