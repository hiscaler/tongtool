package erp2

import (
	"fmt"
	"testing"
)

func TestService_Orders(t *testing.T) {
	_, ttService := newTestTongTool()
	params := OrderQueryParams{
		SaleDateFrom: "2021-12-01 00:00:00",
		SaleDateTo:   "2021-12-11 23:59:59",
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
	fmt.Println(fmt.Sprintf("Total found %d orders", len(orders)))
}
