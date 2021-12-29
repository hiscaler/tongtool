package erp2

import (
	"testing"
)

func TestService_PurchaseOrderArrival(t *testing.T) {
	_, ttService := newTestTongTool()
	req := PurchaseOrderArrivalRequest{
		PurchaseArrivalList: []PurchaseOrderArrivalItem{},
	}
	ttService.PurchaseOrders()

}
