package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
)

// PurchaseOrderLog 采购单入库日志
type PurchaseOrderLog struct {
	ActualPayments    float64 `json:"actualPayments"`    // 实际已付款金额
	Amount            float64 `json:"amount"`            // 采购金额
	Currency          string  `json:"currency"`          // 币种
	GoodsDetailId     string  `json:"goodsDetailId"`     // 通途商品 ID Key
	PurchaseDate      string  `json:"purchaseDate"`      // 采购单创建时间
	PurchaseOrderCode string  `json:"purchaseOrderCode"` // 采购单号
	PurchaseOrderId   string  `json:"purchaseOrderId"`   // 采购单ID
	Quantity          int     `json:"quantity"`          // 采购数量
	ShippingCost      float64 `json:"shippingCost"`      // 采购运费
	SKU               string  `json:"sku"`               // SKU
	SupplierName      string  `json:"supplierName"`      // 供应商
	TrackingNum       string  `json:"trackingNum"`       // 跟踪号
	UnitPrice         float64 `json:"unitPrice"`         // 采购单价
	WarehouseId       string  `json:"warehouseId"`       // 通途仓库 ID key
	WarehouseName     string  `json:"warehouseName"`     // 仓库名称
	WarehousingDate   string  `json:"warehousingDate"`   // 当前入库时间
	WarehousingNum    int     `json:"warehousingNum"`    // 当前入库数量
}

type PurchaseOrderLogQueryParams struct {
	MerchantId          string `json:"merchantId"`                  // 商户ID
	PageNo              int    `json:"pageNo,omitempty"`            // 当前页
	PageSize            int    `json:"pageSize,omitempty"`          // 每页数量
	PurchaseOrderCode   string `json:"purchaseOrderCode,omitempty"` // 采购单号
	WarehousingDateFrom string `json:"warehousingDateFrom"`         // 起始入库时间
	WarehousingDateTo   string `json:"warehousingDateTo"`           // 截止入库时间
}

// PurchaseOrderLogs 采购单入库查询
// https://open.tongtool.com/apiDoc.html#/?docId=eeaf137bae9049b5b087dc0de1ded27a
func (s service) PurchaseOrderLogs(params PurchaseOrderLogQueryParams) (items []PurchaseOrderLog, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	items = make([]PurchaseOrderLog, 0)
	res := struct {
		result
		Datas struct {
			Array    []PurchaseOrderLog `json:"array"`
			PageNo   int                `json:"pageNo"`
			PageSize int                `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/purchaseStockQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				isLastPage = len(items) < params.PageSize
			}
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}
	return
}
