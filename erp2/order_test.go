package erp2

import (
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
		t.Errorf(err.Error())
	} else if !strings.EqualFold(order.OrderIdCode, orderNumber) {
		t.Errorf("order.OrderIdKey %s not match %s", order.OrderIdCode, orderNumber)
	}
}
