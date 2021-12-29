package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
	"strings"
	"testing"
)

func TestService_Orders(t *testing.T) {
	_, ttService := newTestTongTool()
	params := OrderQueryParams{
		SaleDateFrom: "2021-12-01 00:00:00",
		SaleDateTo:   "2021-12-31 23:59:59",
	}
	orders := make([]Order, 0)
	for {
		pageOrders, isLastPage, err := ttService.Orders(params)
		if err != nil {
			t.Errorf("ttService.Orders error: %s", err.Error())
		} else {
			orders = append(orders, pageOrders...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
}

func TestService_Order(t *testing.T) {
	orderNumber := "L-M20211208145011174"
	_, ttService := newTestTongTool()
	order, err := ttService.Order(orderNumber)
	if err != nil {
		if errors.Is(err, tongtool.ErrNotFound) {
			t.Errorf("%s not exists in tongtool", orderNumber)
		} else {
			t.Errorf("ttService.Order error: %s", err.Error())
		}
	} else if !strings.EqualFold(order.OrderIdCode, orderNumber) {
		t.Errorf("order.OrderIdKey %s not match %s", order.OrderIdCode, orderNumber)
	}
}

func TestService_OrderNotFound(t *testing.T) {
	orderNumber := "L-M20211208145011174-bad-number"
	_, ttService := newTestTongTool()
	_, err := ttService.Order(orderNumber)
	if err == nil {
		t.Errorf("ttService.Order except error")
	} else {
		if !errors.Is(err, tongtool.ErrNotFound) {
			t.Error("ttService.Order except not found error")
		}
	}
}

func TestService_CreateOrder(t *testing.T) {
	_, ttService := newTestTongTool()
	req := CreateOrderRequest{
		BuyerInfo: OrderBuyer{
			BuyerAccount:     "test",
			BuyerAddress1:    "test address1",
			BuyerAddress2:    "test address2",
			BuyerAddress3:    "test address3",
			BuyerCity:        "深圳",
			BuyerCountryCode: "CN",
			BuyerEmail:       "happy__snow@126.com",
			BuyerMobilePhone: "15211111111",
			BuyerName:        "张三",
			BuyerPhone:       "15211111113",
			BuyerPostalCode:  "510000",
			BuyerState:       "test",
		},
		Currency:                "CNY",
		InsuranceIncome:         0,
		InsuranceIncomeCurrency: "CNY",
		NeedReturnOrderId:       "1",
		Notes:                   "test notes",
		OrderCurrency:           "CNY",
		PaymentInfos: []OrderPayment{
			{
				OrderAmount:           0,
				OrderAmountCurrency:   "CNY",
				PaymentAccount:        "abc",
				PaymentDate:           "2021-12-29 16:53:08",
				PaymentMethod:         "",
				PaymentNotes:          "text",
				PaymentTransactionNum: "132456798",
				RecipientAccount:      "cba",
				URL:                   "https://www.example.com",
			},
		},
		PlatformCode:       "ebay_api",
		Remarks:            "",
		SaleRecordNum:      "O123456",
		SellerAccountCode:  "test",
		ShippingMethodId:   "",
		TaxIncome:          0,
		TaxIncomeCurrency:  "CNY",
		TotalPrice:         0,
		TotalPriceCurrency: "CNY",
		Transactions: []OrderTransaction{
			{
				GoodsDetailId:              "",
				GoodsDetailRemark:          "",
				ProductsTotalPrice:         0,
				ProductsTotalPriceCurrency: "CNY",
				Quantity:                   1,
				ShipType:                   "速卖通线上发货",
				ShippingFeeIncome:          0,
				ShippingFeeIncomeCurrency:  "CNY",
				SKU:                        "sku-abc",
			},
		},
		WarehouseId: "test",
	}
	orderId, orderNumber, err := ttService.CreateOrder(req)
	if err != nil {
		t.Errorf("ttService.CreateOrder error: %s", err.Error())
	} else {
		t.Logf("orderId: %s, orderNumber: %s", orderId, orderNumber)
	}
}
