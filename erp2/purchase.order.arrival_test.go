package erp2

import (
	"testing"
)

func TestService_PurchaseOrderArrival(t *testing.T) {
	_, ttService := newTestTongTool()
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
