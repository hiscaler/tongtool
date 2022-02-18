package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	"strconv"
)

type ShopifyOrderItem struct {
	ItemId        string  `json:"itemId"`        // 订单交易号
	Price         string  `json:"price"`         // 价格
	PriceValue    float64 `json:"priceValue"`    // 价格
	ProductId     string  `json:"productId"`     // 产品ID
	Quantity      string  `json:"quantity"`      // 产品数量
	QuantityValue int     `json:"quantityValue"` // 产品数量
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

type ShopifyOrdersQueryParams struct {
	Paging
	BuyerEmail          string `json:"buyerEmail,omitempty"`          // 买家邮箱
	MerchantId          string `json:"merchantId"`                    // 商户ID
	PayDateFrom         string `json:"payDateFrom,omitempty"`         // 付款起始时间
	PayDateTo           string `json:"payDateTo,omitempty"`           // 付款结束时间
	PaypalTransactionId string `json:"paypalTransactionId,omitempty"` // Paypal 交易号/Shopify 订单号/付款时间范围 必传其一
	ShopifyOrderId      string `json:"shopifyOrderId,omitempty"`      // Shopify 订单号
}

func (m ShopifyOrdersQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PayDateFrom, validation.When(m.PayDateFrom != "", validation.Date(constant.DatetimeFormat).Error("付款起始时间格式错误"))),
		validation.Field(&m.PayDateTo, validation.When(m.PayDateTo != "", validation.Date(constant.DatetimeFormat).Error("付款结束时间格式错误"))),
		validation.Field(&m.PaypalTransactionId, validation.When(m.PayDateFrom == "" && m.PayDateTo == "" && m.ShopifyOrderId == "", validation.Required.Error("Paypal 交易号/Shopify 订单号/付款时间范围"))),
		validation.Field(&m.ShopifyOrderId, validation.When(m.PayDateFrom == "" && m.PayDateTo == "" && m.PaypalTransactionId == "", validation.Required.Error("Paypal 交易号/Shopify 订单号/付款时间范围"))),
	)
}

// ShopifyOrders Shopify 订单列表
// https://open.tongtool.com/apiDoc.html#/?docId=e949a88561e7471785cccef86feb3e6d
func (s service) ShopifyOrders(params ShopifyOrdersQueryParams) (items []ShopifyOrder, isLastPage bool, err error) {
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
	items = make([]ShopifyOrder, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []ShopifyOrder `json:"array"`
			PageNo   int            `json:"pageNo"`
			PageSize int            `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/shopifyOrderQuery")
	if err != nil {
		return
	}

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
