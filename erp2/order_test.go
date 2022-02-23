package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
	"strings"
	"sync"
	"testing"
)

func TestService_Orders(t *testing.T) {
	params := OrdersQueryParams{
		AccountCode:  "LDXAUS",
		SaleDateFrom: "2021-12-01 00:00:00",
		SaleDateTo:   "2021-12-31 23:59:59",
	}
	params.PageNo = 1
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

func TestService_GoroutineOrders(t *testing.T) {
	params := OrdersQueryParams{
		SaleDateFrom: "2021-12-01 00:00:00",
		SaleDateTo:   "2021-12-01 23:59:59",
	}
	var wg sync.WaitGroup
	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
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
			t.Logf("%d: found %d orders", i, len(orders))
		}()
	}
	wg.Wait()
}

func TestService_Order(t *testing.T) {
	orderNumber := "abc"
	order, _, err := ttService.Order(orderNumber)
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
	_, exists, err := ttService.Order(orderNumber)
	if err == nil {
		t.Errorf("ttService.Order except error")
	} else if exists {
		t.Errorf("ttService.Order: this is a invalid order number, but return it.")
	} else if !errors.Is(err, tongtool.ErrNotFound) {
		t.Error("ttService.Order except tongtool.ErrNotFound error")
	}
}

func TestService_CreateOrder(t *testing.T) {
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
		Remarks:            []string{},
		SaleRecordNum:      "O1234567",
		SellerAccountCode:  "test",
		ShippingMethodId:   "",
		TaxIncome:          0,
		TaxIncomeCurrency:  "CNY",
		TotalPrice:         0,
		TotalPriceCurrency: "CNY",
		Transactions: []OrderTransaction{
			{
				GoodsDetailId:              "",
				GoodsDetailRemark:          "货品备注",
				ProductsTotalPrice:         2,
				ProductsTotalPriceCurrency: "CNY",
				Quantity:                   2,
				ShipType:                   "速卖通线上发货",
				ShippingFeeIncome:          2,
				ShippingFeeIncomeCurrency:  "CNY",
				SKU:                        "goods_sku",
			},
		},
		WarehouseId: "0001000007201303040000013106",
	}
	orderId, orderNumber, err := ttService.CreateOrder(req)
	if err != nil {
		t.Errorf("ttService.CreateOrder error: %s", err.Error())
	} else {
		t.Logf("orderId: %s, orderNumber: %s", orderId, orderNumber)
	}
}

func TestService_UpdateOrder(t *testing.T) {
	req := UpdateOrderRequest{
		OrderId: "abc",
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
		Remarks:          []string{},
		ShippingMethodId: "",
		Transactions: []UpdateOrderTransaction{
			{
				GoodsDetailId: "",
				Quantity:      2,
			},
		},
		WarehouseId: "0001000007201303040000013106",
	}
	err := ttService.UpdateOrder(req)
	if err != nil {
		t.Errorf("ttService.UpdateOrder error: %s", err.Error())
	}
}
