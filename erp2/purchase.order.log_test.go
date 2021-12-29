package erp2

import "testing"

func TestService_PurchaseOrderLogs(t *testing.T) {
	_, ttService := newTestTongTool()
	params := PurchaseOrderLogQueryParams{
		PurchaseOrderCode:   "PO002057",
		WarehousingDateFrom: "2021-11-01 00:00:00",
		WarehousingDateTo:   "2021-12-31 23:59:59",
	}
	orders := make([]PurchaseOrderLog, 0)
	for {
		pageOrders, isLastPage, err := ttService.PurchaseOrderLogs(params)
		if err != nil {
			t.Errorf("ttService.PurchaseOrderLogs error: %s", err.Error())
		} else {
			orders = append(orders, pageOrders...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
}
