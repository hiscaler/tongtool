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
			fmt.Println(cast.ToJson(orders))
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

// 创建采购单
func TestService_CreatePurchaseOrder(t *testing.T) {
	_, ttService := newTestTongTool()
	details := []PurchaseOrderGoodDetail{
		{GoodsDetailId: "8309050530202104270001946885", Quantity: 1, UnitPrice: 1.1},
		{GoodsDetailId: "8309050530202106100002312298", Quantity: 2, UnitPrice: 2.2},
	}
	req := CreatePurchaseOrderRequest{
		Currency:       "CNY",
		GoodsDetail:    details,
		ExternalNumber: "",
		PurchaseUserId: "202012180006653303",
		Remark:         "test for purchase order create",
		ShippingFee:    6.6,
		SupplierId:     "8309050530202107230004245350",
		TrackingNumber: "",
		WarehouseIdKey: "8151050530202008250000047045",
	}

	number, err := ttService.CreatePurchaseOrder(req)
	if err != nil {
		t.Errorf("create purchase order error: %s", err.Error())
	} else {
		fmt.Println(fmt.Sprintf("Purchase number: %s", number))
	}
}
