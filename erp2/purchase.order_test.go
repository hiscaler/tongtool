package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"strings"
	"testing"
)

// 采购单查询
func TestService_PurchaseOrders(t *testing.T) {
	number := "PO000007"
	params := PurchaseOrdersQueryParams{
		PurchaseOrderCode: number,
	}
	orders, _, err := ttService.PurchaseOrders(params)
	if err == nil {
		exists := false
		for _, order := range orders {
			if strings.EqualFold(number, order.PONum) {
				exists = true
				break
			}
		}
		if exists {
			fmt.Println(jsonx.ToJson(orders, "[]"))
		} else {
			t.Errorf("not found %s", number)
		}
	} else {
		t.Error(err)
	}
}

func TestService_PurchaseOrdersByStatus(t *testing.T) {
	status := "delivering"
	params := PurchaseOrdersQueryParams{
		POrderStatus: status,
	}
	orders, _, err := ttService.PurchaseOrders(params)
	if err == nil {
		fmt.Println(jsonx.ToJson(orders, "[]"))
	} else {
		t.Error(err)
	}
}

// 创建采购单
func TestService_CreatePurchaseOrder(t *testing.T) {
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

func TestService_PurchaseOrderStockInLogs(t *testing.T) {
	params := PurchaseOrderLogsQueryParams{
		PurchaseOrderCode:   "PO002057",
		WarehousingDateFrom: "2021-11-01 00:00:00",
		WarehousingDateTo:   "2021-12-31 23:59:59",
	}
	orders := make([]PurchaseOrderLog, 0)
	for {
		pageOrders, isLastPage, err := ttService.PurchaseOrderStockInLogs(params)
		if err != nil {
			t.Errorf("ttService.PurchaseOrderStockInLogs error: %s", err.Error())
		} else {
			orders = append(orders, pageOrders...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
}

func TestService_PurchaseOrderArrival(t *testing.T) {
	req := PurchaseOrderArrivalRequest{
		PurchaseArrivalList: []PurchaseOrderArrivalItem{
			{
				ArrivalGoodsList: []PurchaseOrderArrivalGoodsItem{
					{
						GoodsDetailId:        "123",
						InQuantity:           1,
						IsReplace:            "N",
						ReplaceGoodsDetailId: "",
						ReplaceQuantity:      0,
					},
				},
				Freight:           1,
				PurchaseOrderCode: "PO123",
				Remark:            "备注",
			},
		},
	}
	err := ttService.PurchaseOrderArrival(req)
	if err != nil {
		t.Errorf("ttService.PurchaseOrderArrival error: %s", err.Error())
	}
}
