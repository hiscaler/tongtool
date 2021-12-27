package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/pkg/is"
	"strings"
)

// 采购单状态
const (
	PurchaseOrderStatusDelivering        = "0" // 等待到货、未全部到货
	PurchaseOrderStatusPReceivedAndWaitM = "1" // 部分到货等待剩余
	PurchaseOrderStatusPartialReceived   = "2" // 部分到货不等待剩余
	PurchaseOrderStatusReceived          = "3" // 全部到货
	PurchaseOrderStatusCancel            = "4" // 作废
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
	ActualPayments      float64 `json:"actual_payments"`      // 实际已付款金额
	Amount              float64 `json:"amount"`               // 采购金额
	CorporationFullName string  `json:"corporation_fullname"` // 供应商名称
	CreatedDate         string  `json:"createdDate"`          // 采购单创建时间
	Currency            string  `json:"currency"`             // 币种
	GoodsIdKey          string  `json:"goodsIdKey"`           // 通途商品id key
	GoodsSKU            string  `json:"goods_sku"`            // 商品 SKU
	InQuantity          int     `json:"in_quantity"`          // 已入库数量
	PayableAmounts      float64 `json:"payableAmounts"`       // 应付金额
	PoNum               string  `json:"ponum"`                // 采购单号
	PurchaseArrivalDate string  `json:"purchaseArrivalDate"`  // 采购到货时间
	PurchaseDate        string  `json:"purchaseDate"`         // 采购日期
	PurchaseOrderId     string  `json:"purchaseOrderId"`      // 采购单id
	Quantity            int     `json:"quantity"`             // 采购数量
	ShippingCost        float64 `json:"shipping_cost"`        // 采购运费
	Status              string  `json:"status"`               // 采购单状态0-等待到货、未全部到货, 1-部分到货等待剩余, 2-部分到货不等待剩余, 3-全部到货, 4-作废
	SupplierCode        string  `json:"supplier_code"`        // 供应商代码
	TrackingNumber      string  `json:"tracking_number"`      // 跟踪号
	UnitPrice           float64 `json:"unit_price"`           // 采购单价
	WarehouseIdKey      string  `json:"warehouseIdKey"`       // 通途仓库id key
	WarehouseName       string  `json:"warehouseName"`        // 仓库名称
	WillArriveDate      string  `json:"willArriveDate"`       // 预计到达日期
}

type PurchaseOrdersQueryParams struct {
	MerchantId        string `json:"merchantId"`                   // 商户ID
	POrderStatus      string `json:"pOrderStatus,omitempty"`       // 采购单状态:delivering/等待到货 、pReceivedAndWaitM/部分到货等待剩余、partialReceived/部分到货不等待剩余、Received/全部到货、cancel/已作废、NotPaymentApply/未申请付款、paymentApply/已申请付款、paymentCancel/已取消付款、payed/已付款、partialPayed/部分付款
	PurchaseDateFrom  string `json:"purchaseDateFrom,omitempty"`   // 采购日期开始时间
	PurchaseDateTo    string `json:"purchaseDateTo,omitempty"`     // 采购日期结束时间
	PurchaseOrderCode string `json:"purchaseOrderCode,omitempty"`  // 采购单号
	SKUs              string `json:"skus,omitempty"`               // SKU数组，长度不超过10
	SupplierName      string `json:"supplierName,omitempty"`       // 供应商名称
	UpdatedDateFrom   string `json:"updatedDateFrom,omitempty"`    // 采购单更新开始时间
	UpdatedDateTo     string `json:"updatedDateTo,omitempty"`      // 采购单更新结束时间
	PageNo            int    `json:"pageNo,omitempty,omitempty"`   // 查询页数
	PageSize          int    `json:"pageSize,omitempty,omitempty"` // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
}

type purchaseOrdersResult struct {
	result
	Datas struct {
		Array    []PurchaseOrder `json:"array"`
		PageNo   int             `json:"pageNo"`
		PageSize int             `json:"pageSize"`
	} `json:"datas,omitempty"`
}

// PurchaseOrders 采购单列表
// https://open.tongtool.com/apiDoc.html#/?docId=0dd564d52ce34ad0afce1f304d6b7824
func (s service) PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	if is.Number(params.POrderStatus) {
		params.POrderStatus = PurchaseOrderStatusNtoS(params.POrderStatus)
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
			if cpr.Code == tongtool.OK {
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
