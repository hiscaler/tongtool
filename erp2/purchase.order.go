package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/pkg/cache"
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

// PurchaseOrders 采购单列表
// https://open.tongtool.com/apiDoc.html#/?docId=0dd564d52ce34ad0afce1f304d6b7824
func (s service) PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	if is.Number(params.POrderStatus) {
		params.POrderStatus = PurchaseOrderStatusNtoS(params.POrderStatus)
	}
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = cache.GenerateKey(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = json.Unmarshal(b, &items); e == nil {
				return
			} else {
				s.tongTool.Logger.Printf(`cache data unmarshal error
 DATA: %s
ERROR: %s
`, string(b), e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("get cache %s error: %s", cacheKey, e.Error())
		}
	}
	items = make([]PurchaseOrder, 0)
	res := struct {
		result
		Datas struct {
			Array    []PurchaseOrder `json:"array"`
			PageNo   int             `json:"pageNo"`
			PageSize int             `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
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

	if err == nil && s.tongTool.EnableCache {
		if b, e := json.Marshal(&items); e == nil {
			e = s.tongTool.Cache.Set(cacheKey, b)
			if e != nil {
				s.tongTool.Logger.Printf("set cache %s error: %s", cacheKey, e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("items marshal error: %s", err.Error())
		}
	}
	return
}

// 采购单创建
// https://open.tongtool.com/apiDoc.html#/?docId=bcfd5d50a664486298b7fb0c1d08f714

type PurchaseOrderGoodDetail struct {
	GoodsDetailId string  `json:"goodsDetailId"` // 通途货品 ID
	Quantity      int     `json:"quantity"`      // 采购数量
	UnitPrice     float64 `json:"unitPrice"`     // 采购单价
}

type CreatePurchaseOrderRequest struct {
	Currency       string                    `json:"currency"`       // 币种
	GoodsDetail    []PurchaseOrderGoodDetail `json:"goodsDetail"`    // 采购货品信息
	ExternalNumber string                    `json:"externalNumber"` // 外部流水号
	MerchantId     string                    `json:"merchantId"`     // 商户 ID
	PurchaseUserId string                    `json:"purchaseUserId"` // 采购员 ID
	Remark         string                    `json:"remark"`         // 采购备注
	ShippingFee    float64                   `json:"shippingFee"`    // 运费
	SupplierId     string                    `json:"supplierId"`     // 通途供应商 ID
	TrackingNumber string                    `json:"trackingNumber"` // 跟踪号
	WarehouseIdKey string                    `json:"warehouseIdKey"` // 通途仓库 ID
}

func (m CreatePurchaseOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Currency, validation.Required.Error("币种不能为空")),
		validation.Field(&m.GoodsDetail, validation.Required.Error("采购货品信息不能为空"), validation.By(func(value interface{}) error {
			items, ok := value.([]PurchaseOrderGoodDetail)
			if !ok {
				return errors.New("无效的采购货品信息")
			}
			for _, item := range items {
				if item.GoodsDetailId == "" {
					return errors.New("采购货品中通途货品 ID 不能为空")
				}
				if item.Quantity <= 0 {
					return errors.New("采购货品中采购数量不能小于 1")
				}
			}
			return nil
		})),
		validation.Field(&m.PurchaseUserId, validation.Required.Error("采购员 ID 不能为空")),
		validation.Field(&m.SupplierId, validation.Required.Error("通途供应商 ID 不能为空")),
		validation.Field(&m.WarehouseIdKey, validation.Required.Error("通途仓库 ID 不能为空")),
	)
}

// CreatePurchaseOrder 创建采购单
func (s service) CreatePurchaseOrder(req CreatePurchaseOrderRequest) (number string, err error) {
	if err = req.Validate(); err != nil {
		return
	}
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

// 采购单入库处理
// https://open.tongtool.com/apiDoc.html#/?docId=21d1c988af2d4dc5940d1faf105d5a46

// PurchaseOrderStockInItem 采购单入库项
type PurchaseOrderStockInItem struct {
	GoodsDetailId string `json:"goodsDetailId"` // 通途货品 ID
	Quantity      int    `json:"quantity"`      // 采购数量
}

// PurchaseOrderStockInRequest 采购单入库
type PurchaseOrderStockInRequest struct {
	ArrivalInfoList []PurchaseOrderStockInItem `json:"arrivalInfoList"` // 到货货品信息
	Freight         float64                    `json:"freight"`         // 运费
	MerchantId      string                     `json:"merchantId"`      // 商家 ID
	PurchaseOrderId string                     `json:"purchaseOrderId"` // 采购单 ID
}

// PurchaseOrderStockIn 采购单入库
func (s service) PurchaseOrderStockIn(req PurchaseOrderStockInRequest) error {
	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		result
		Datas interface{} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/purchaseOrderStockIn")
	if err == nil {
		if resp.IsSuccess() {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}
	return err
}

// 采购单入库记录查询
// https://open.tongtool.com/apiDoc.html#/?docId=eeaf137bae9049b5b087dc0de1ded27a

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

// PurchaseOrderStockInLogs 采购单入库记录查询
func (s service) PurchaseOrderStockInLogs(params PurchaseOrderLogQueryParams) (items []PurchaseOrderLog, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = cache.GenerateKey(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = json.Unmarshal(b, &items); e == nil {
				return
			} else {
				s.tongTool.Logger.Printf(`cache data unmarshal error
 DATA: %s
ERROR: %s
`, string(b), e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("get cache %s error: %s", cacheKey, e.Error())
		}
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

	if err == nil && s.tongTool.EnableCache {
		if b, e := json.Marshal(&items); e == nil {
			e = s.tongTool.Cache.Set(cacheKey, b)
			if e != nil {
				s.tongTool.Logger.Printf("set cache %s error: %s", cacheKey, e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("items marshal error: %s", err.Error())
		}
	}
	return
}

// 采购单到货处理
// https://open.tongtool.com/apiDoc.html#/?docId=ee942453af114a7686d0c8d5187988f2

// PurchaseOrderArrivalItem 采购到货项
type PurchaseOrderArrivalItem struct {
	ArrivalGoodsList  []PurchaseOrderArrivalGoodsItem `json:"arrivalGoodsList"`  // 采购到货明细
	Freight           float64                         `json:"freight"`           // 运费
	PurchaseOrderCode string                          `json:"purchaseOrderCode"` // 采购单号
	Remark            string                          `json:"remark"`            // 到货备注
}

// PurchaseOrderArrivalGoodsItem 采购到货项
type PurchaseOrderArrivalGoodsItem struct {
	GoodsDetailId        string `json:"goodsDetailId"`        // 通途货品ID
	InQuantity           int    `json:"inQuantity"`           // 到货数量
	IsReplace            string `json:"isReplace"`            // 是否是变参替换到货：[Y：是]
	ReplaceGoodsDetailId string `json:"replaceGoodsDetailId"` // 变参替换的通途货品ID
	ReplaceQuantity      int    `json:"replaceQuantity"`      // 变参替换的到货数量
}

type PurchaseOrderArrivalRequest struct {
	MerchantId          string                     `json:"merchantId"`          // 商户ID
	PurchaseArrivalList []PurchaseOrderArrivalItem `json:"purchaseArrivalList"` // 采购到货列表
}

// PurchaseOrderArrival 采购单到货
func (s service) PurchaseOrderArrival(req PurchaseOrderArrivalRequest) error {
	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		result
		Datas interface{} `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/purchaseArrival")
	if err == nil {
		if resp.IsSuccess() {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	return err
}
