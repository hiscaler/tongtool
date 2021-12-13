package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
	"strings"
)

// Order 通途订单
type Order struct {
	ActualTotalPrice        float64       `json:"actualTotalPrice"`
	AssignStockCompleteTime string        `json:"assignstockCompleteTime"`
	BuyerAccountId          string        `json:"buyerAccountId"`
	BuyerCity               string        `json:"buyerCity"`
	BuyerCountry            string        `json:"buyerCountry"`
	BuyerEmail              string        `json:"buyerEmail"`
	BuyerMobile             string        `json:"buyerMobile"`
	BuyerName               string        `json:"buyerName"`
	BuyerPassportCode       string        `json:"buyerPassportCode"`
	BuyerPhone              string        `json:"buyerPhone"`
	BuyerState              string        `json:"buyerState"`
	ReceiveAddress          string        `json:"receiveAddress"`
	PostalCode              string        `json:"postalCode"`
	Carrier                 string        `json:"carrier"`
	CarrierType             string        `json:"carrierType"`
	CarrierURL              string        `json:"carrierUrl"`
	DespatchCompleteTime    string        `json:"despatchCompleteTime"`
	DispathTypeName         string        `json:"dispathTypeName"`
	EbayNotes               string        `json:"ebayNotes"`
	EbaySiteEnName          string        `json:"ebaySiteEnName"`
	FirstTariff             float64       `json:"firstTariff"`
	InsuranceIncome         float64       `json:"insuranceIncome"`
	InsuranceIncomeCurrency float64       `json:"insuranceIncomeCurrency"`
	ProductsTotalPrice      float64       `json:"productsTotalPrice"`
	OrderAmount             float64       `json:"orderAmount"`
	OrderAmountCurrency     string        `json:"orderAmountCurrency"`
	OrderDetails            []OrderDetail `json:"orderDetails"`
	GoodsInfo               GoodsInfo     `json:"goodsInfo"`
	OrderIdCode             string        `json:"orderIdCode"`
	OrderIdKey              string        `json:"orderIdKey"`
	OrderStatus             string        `json:"orderStatus"`
	IsInvalid               string        `json:"isInvalid"`
	SaleTime                string        `json:"saleTime"`
	PaidTime                string        `json:"paidTime"`
	RefundedTime            string        `json:"refundedTime"` // 退款日期
	ShippingFeeIncome       float64       `json:"shippingFeeIncome"`
}

// OrderDetail 通途订单详情
type OrderDetail struct {
	GoodsMatchedQuantity int     `json:"goodsMatchedQuantity"`
	GoodsMatchedSku      string  `json:"goodsMatchedSku"`
	OrderDetailsId       string  `json:"orderDetailsId"`
	Quantity             int     `json:"quantity"`
	TransactionPrice     float64 `json:"transaction_price"`
	WebstoreCustomLabel  string  `json:"webstore_custom_label"`
	WebstoreItemId       string  `json:"webstore_item_id"`
	WebstoreSKU          string  `json:"webstore_sku"`
}

type TongToolGoodsInfoList struct {
	GoodsSku           string  `json:"goodsSku"`
	ProductCurrentCost float64 `json:"productCurrentCost"`
	GoodsCurrentCost   float64 `json:"goodsCurrentCost"`
	GoodsAverageCost   float64 `json:"goodsAverageCost"`
	GoodsImageGroupId  string  `json:"goodsImageGroupId"`
	ProductName        string  `json:"productName"`
}
type GoodsInfo struct {
	TongToolGoodsInfoList []TongToolGoodsInfoList `json:"tongToolGoodsInfoList"`
}

type OrderQueryParams struct {
	AccountCode      string `json:"accountCode"`
	BuyerEmail       string `json:"buyerEmail,omitempty"`
	MerchantId       string `json:"merchantId"`
	OrderId          string `json:"orderId,omitempty"`
	OrderStatus      string `json:"orderStatus,omitempty"`
	PageNo           int    `json:"pageNo,omitempty"`
	PageSize         int    `json:"pageSize,omitempty"`
	PayDateFrom      string `json:"payDateFrom,omitempty"`
	PayDateTo        string `json:"payDateTo,omitempty"`
	PlatformCode     string `json:"platformCode,omitempty"`
	RefundedDateFrom string `json:"refundedDateFrom,omitempty"`
	RefundedDateTo   string `json:"refundedDateTo,omitempty"`
	SaleDateFrom     string `json:"saleDateFrom,omitempty"`
	SaleDateTo       string `json:"saleDateTo,omitempty"`
	StoreFlag        int    `json:"storeFlag"`
	UpdatedDateFrom  string `json:"updatedDateFrom,omitempty"`
	UpdatedDateTo    string `json:"updatedDateTo,omitempty"`
}

type orderResult struct {
	result
	Datas struct {
		Array    []Order `json:"array"`
		PageNo   int     `json:"pageNo"`
		PageSize int     `json:"pageSize"`
	}
}

func (s service) Orders(params OrderQueryParams) (items []Order, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = s.tongTool.QueryDefaultValues.PageNo
	}
	if params.PageSize <= 0 {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]Order, 0)
	res := orderResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/ordersQuery")
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

func (s service) Order(id string) (item Order, err error) {
	params := OrderQueryParams{
		OrderId:  id,
		PageNo:   1,
		PageSize: s.tongTool.QueryDefaultValues.PageSize,
	}
	for {
		items := make([]Order, 0)
		isLastPage := false
		items, isLastPage, err = s.Orders(params)
		if err == nil {
			if len(items) == 0 {
				err = errors.New("not found")
			} else {
				for _, order := range items {
					if strings.EqualFold(order.OrderIdKey, id) {
						item = order
						return
					}
				}
			}
		}
		if err != nil || isLastPage {
			break
		}
		params.PageNo++
	}
	return
}
