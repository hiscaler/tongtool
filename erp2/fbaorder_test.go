package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/cast"
	"testing"
)

func TestService_FBAOrders(t *testing.T) {
	_, ttService := newTestTongTool()
	params := FBAOrderQueryParams{
		PurchaseDateFrom: "2021-12-01 00:00:00",
		PurchaseDateTo:   "2021-12-10 23:59:59",
	}
	for {
		orders, isLastPage, err := ttService.FBAOrders(params)
		if err != nil {
			t.Errorf("ttService.FBAOrders error: %s", err.Error())
		} else {
			fmt.Println(cast.ToJson(orders))
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
}
