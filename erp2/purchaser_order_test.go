package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"strings"
	"testing"
)

// 采购单查询
func TestService_PurchaseOrders(t *testing.T) {
	_, ttService := newTestTongTool()
	number := "PO000007"
	params := PurchaseOrdersQueryParams{
		PurchaseOrderCode: number,
	}
	orders, _, err := ttService.PurchaseOrders(params)
	if err == nil {
		exists := false
		for _, order := range orders {
			if strings.EqualFold(number, order.PoNum) {
				exists = true
				break
			}
		}
		if exists {
			fmt.Println(fmt.Sprintf("Orders: %#v", orders))
		} else {
			t.Errorf("not found %s", number)
		}
	} else {
		t.Error(err)
	}
}

func TestService_PurchaseOrdersByStatus(t *testing.T) {
	_, ttService := newTestTongTool()
	status := "delivering"
	params := PurchaseOrdersQueryParams{
		POrderStatus: status,
	}
	orders, _, err := ttService.PurchaseOrders(params)
	if err == nil {
		fmt.Println(cast.ToJson(orders))
	} else {
		t.Error(err)
	}
}
