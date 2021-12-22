package erp2

import (
	"fmt"
	"testing"
)

func TestService_ShopifyOrders(t *testing.T) {
	_, ttService := newTestTongTool()
	params := ShopifyOrderQueryParams{
		PayDateFrom: "2021-12-01 00:00:00",
		PayDateTo:   "2021-12-11 23:59:59",
	}
	orders := make([]ShopifyOrder, 0)
	for {
		pageOrders, isLastPage, err := ttService.ShopifyOrders(params)
		if err != nil {
			t.Errorf("ttService.ShopifyOrders error: %s", err.Error())
		} else {
			fmt.Println(fmt.Sprintf("%#v", pageOrders))
			orders = append(orders, pageOrders...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	fmt.Println(fmt.Sprintf("Total found %d orders", len(orders)))
}