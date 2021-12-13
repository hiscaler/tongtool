package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
	"strings"
)

type PurchaseOrderGoodDetail struct {
	GoodsDetailId string  `json:"goodsDetailId"`
	Quantity      int     `json:"quantity"`
	UnitPrice     float64 `json:"unitPrice"`
}

type CreatePurchaseOrderRequest struct {
	Currency       string                    `json:"currency"`
	GoodsDetail    []PurchaseOrderGoodDetail `json:"goodsDetail"`
	ExternalNumber string                    `json:"externalNumber"`
	MerchantId     string                    `json:"merchantId"`
	PurchaseUserId string                    `json:"purchaseUserId"`
	Remark         string                    `json:"remark"`
	ShippingFee    float64                   `json:"shippingFee"`
	SupplierId     string                    `json:"supplierId"`
	TrackingNumber string                    `json:"trackingNumber"`
	WarehouseIdKey string                    `json:"warehouseIdKey"`
}

type PurchaseOrder struct {
	ActualPayments      float64 `json:"actual_payments"`
	Amount              float64 `json:"amount"`
	CorporationFullName string  `json:"corporation_fullname"`
	CreatedDate         string  `json:"createdDate"`
	Currency            string  `json:"currency"`
	GoodsIdKey          string  `json:"goodsIdKey"`
	GoodsSKU            string  `json:"goods_sku"`
	InQuantity          int     `json:"in_quantity"`
	PayableAmounts      float64 `json:"payableAmounts"`
	PoNum               string  `json:"ponum"`
	PurchaseArrivalDate string  `json:"purchaseArrivalDate"`
	PurchaseDate        string  `json:"purchaseDate"`
	PurchaseOrderId     string  `json:"purchaseOrderId"`
	Quantity            int     `json:"quantity"`
	ShippingCost        float64 `json:"shipping_cost"`
	Status              string  `json:"status"`
	SupplierCode        string  `json:"supplier_code"`
	TrackingNumber      string  `json:"tracking_number"`
	UnitPrice           float64 `json:"unit_price"`
	WarehouseIdKey      string  `json:"warehouseIdKey"`
	WarehouseName       string  `json:"warehouseName"`
	WillArriveDate      string  `json:"willArriveDate"`
}

type PurchaseOrdersQueryParams struct {
	MerchantId        string `json:"merchantId"`
	POrderStatus      string `json:"pOrderStatus,omitempty"`
	PurchaseDateFrom  string `json:"purchaseDateFrom,omitempty"`
	PurchaseDateTo    string `json:"purchaseDateTo,omitempty"`
	PurchaseOrderCode string `json:"purchaseOrderCode,omitempty"`
	SKUs              string `json:"skus,omitempty"`
	SupplierName      string `json:"supplierName,omitempty"`
	UpdatedDateFrom   string `json:"updatedDateFrom,omitempty"`
	UpdatedDateTo     string `json:"updatedDateTo,omitempty"`
	PageNo            int    `json:"pageNo,omitempty,omitempty"`
	PageSize          int    `json:"pageSize,omitempty,omitempty"`
}

type purchaseOrdersResult struct {
	result
	Datas struct {
		Array    []PurchaseOrder `json:"array"`
		PageNo   int             `json:"pageNo"`
		PageSize int             `json:"pageSize"`
	}
}

func (s service) PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = s.tongTool.QueryDefaultValues.PageNo
	}
	if params.PageSize <= 0 {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]PurchaseOrder, 0)
	res := purchaseOrdersResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/purchaseOrderQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.HasError(res.Code); err == nil {
				items = res.Datas.Array
				isLastPage = len(items) < params.PageSize
			}
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}

func (s service) CreatePurchaseOrder(req CreatePurchaseOrderRequest) (number string, err error) {
	type createPurchaseOrderResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Datas   string
	}
	cpr := createPurchaseOrderResponse{}
	req.MerchantId = s.tongTool.MerchantId
	r, err := s.tongTool.Client.R().SetResult(&cpr).SetBody(req).Post("/openapi/tongtool/purchaseOrderCreate")
	if err == nil {
		if r.IsSuccess() {
			if cpr.Code == 200 {
				number = strings.TrimSpace(cpr.Datas)
				if number == "" {
					err = errors.New("not found number in http response")
				}
			} else {
				err = errors.New(cpr.Message)
			}
		} else {
			err = errors.New(r.Status())
		}
	}

	return
}
