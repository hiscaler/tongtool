package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/pkg/in"
	"strings"
)

const (
	OrderStoreFlagActive   = "0" // 活跃表
	OrderStoreFlagOneYear  = "1" // 一年表
	OrderStoreFlagArchived = "2" // 归档表
)

// OrderDetail 通途订单详情
type OrderDetail struct {
	GoodsMatchedQuantity int     `json:"goodsMatchedQuantity"`
	GoodsMatchedSku      string  `json:"goodsMatchedSku"`
	OrderDetailsId       string  `json:"orderDetailsId"`
	Quantity             int     `json:"quantity"`
	TransactionPrice     float64 `json:"transaction_price"`
	WebStoreCustomLabel  string  `json:"webstore_custom_label"`
	WebStoreItemId       string  `json:"webstore_item_id"`
	WebStoreSKU          string  `json:"webstore_sku"`
}

// OrderPackage 订单包裹信息
type OrderPackage struct {
	PackageId            string `json:"packageId"`            // 包裹号
	TrackingNumber       string `json:"trackingNumber"`       // 物流跟踪号
	TrackingNumberStatus string `json:"trackingNumberStatus"` // 物流跟踪号获取状态(00:未就绪 01:就绪 02:处理中 03:处理成功 04:处理失败)
	TrackingNumberTime   int    `json:"trackingNumberTime"`   // 物流跟踪号获取时间
}

type TongToolGoodsInfoList struct {
	GoodsSku           string  `json:"goodsSku"`
	ProductCurrentCost float64 `json:"productCurrentCost"`
	GoodsCurrentCost   float64 `json:"goodsCurrentCost"`
	GoodsAverageCost   float64 `json:"goodsAverageCost"`
	GoodsImageGroupId  string  `json:"goodsImageGroupId"`
	ProductName        string  `json:"productName"`
}

type GoodsInfo struct {
	TongToolGoodsInfoList []TongToolGoodsInfoList `json:"tongToolGoodsInfoList"`
}

// Order 通途订单
type Order struct {
	ActualTotalPrice          float64        `json:"actualTotalPrice"`          // 实付金额
	AssignStockCompleteTime   string         `json:"assignstockCompleteTime"`   // 配货时间
	BuyerAccountId            string         `json:"buyerAccountId"`            // 买家id
	BuyerCity                 string         `json:"buyerCity"`                 // 买家城市
	BuyerCountry              string         `json:"buyerCountry"`              // 买家国家
	BuyerEmail                string         `json:"buyerEmail"`                // 买家邮箱
	BuyerMobile               string         `json:"buyerMobile"`               // 买家手机
	BuyerName                 string         `json:"buyerName"`                 // 买家名称
	BuyerPassportCode         string         `json:"buyerPassportCode"`         // 收件人识别码（护照等）
	BuyerPhone                string         `json:"buyerPhone"`                // 买家电话
	BuyerState                string         `json:"buyerState"`                // 买家省份
	Carrier                   string         `json:"carrier"`                   // 上传物流的carrier
	CarrierType               string         `json:"carrierType"`               // 物流商类型 ( 0:通途API对接、 1:通途Excel文件导出、 2:通途离线生成跟踪号 3:无对接、 4:自定义Excel对接)
	CarrierURL                string         `json:"carrierUrl"`                // 物流网络地址
	DespatchCompleteTime      string         `json:"despatchCompleteTime"`      // 订单发货完成时间
	DispatchTypeName          string         `json:"dispathTypeName"`           // 邮寄方式名称
	EbayNotes                 string         `json:"ebayNotes"`                 // 订单备注
	EbaySiteEnName            string         `json:"ebaySiteEnName"`            // 站点
	FirstTariff               float64        `json:"firstTariff"`               // 头程运费
	GoodsInfo                 GoodsInfo      `json:"goodsInfo"`                 // 订单商品信息
	InsuranceIncome           float64        `json:"insuranceIncome"`           // 买家所付保费
	InsuranceIncomeCurrency   string         `json:"insuranceIncomeCurrency"`   // 买家所付保费币种
	IsInvalid                 string         `json:"isInvalid"`                 // 是否作废(0,''，null 未作废，1 手工作废 2 订单任务下载永久作废 3 拆分单主单作废 4 拆分单子单作废)
	IsSuspended               string         `json:"isSuspended"`               // 是否需要人工审核 (1需要人工审核,0或null不需要)
	MerchantCarrierShortname  string         `json:"merchantCarrierShortname"`  // 承运人简称
	OrderAmount               float64        `json:"orderAmount"`               // 订单总金额(商品金额+运费+保费)
	OrderAmountCurrency       string         `json:"orderAmountCurrency"`       // 订单金额币种
	OrderDetails              []OrderDetail  `json:"orderDetails"`              // 订单明细
	OrderIdCode               string         `json:"orderIdCode"`               // 通途订单号
	OrderIdKey                string         `json:"orderIdKey"`                // 通途订单id key
	OrderStatus               string         `json:"orderStatus"`               // 订单状态（waitPacking/等待配货 ,waitPrinting/等待打印 ,waitingDespatching/等待发货 ,despatched/已发货）
	PackageInfoList           []OrderPackage `json:"packageInfoList"`           // 订单包裹信息
	PaidTime                  string         `json:"paidTime"`                  // 订单付款完成时间
	ParentOrderId             string         `json:"parentOrderId"`             // 父订单号
	PlatformCode              string         `json:"platformCode"`              // 通途中平台代码
	PlatformFee               float64        `json:"platformFee"`               // 平台手续费
	PostalCode                string         `json:"postalCode"`                // 买家邮编
	PrintCompleteTime         string         `json:"printCompleteTime"`         // 订单打印完成时间
	ProductsTotalCurrency     string         `json:"productsTotalCurrency"`     // 金额小计币种
	ProductsTotalPrice        float64        `json:"productsTotalPrice"`        // 金额小计(只商品金额)
	ReceiveAddress            string         `json:"receiveAddress"`            // 收货地址
	RefundedTime              string         `json:"refundedTime"`              // 退款时间
	SaleAccount               string         `json:"saleAccount"`               // 卖家账号
	SaleTime                  string         `json:"saleTime"`                  // 订单生成时间
	SalesRecordNumber         string         `json:"salesRecordNumber"`         // 平台订单号
	ShippingFee               float64        `json:"shippingFee"`               // 关税
	ShippingFeeIncome         float64        `json:"shippingFeeIncome"`         // 买家所支付的运费
	ShippingFeeIncomeCurrency string         `json:"shippingFeeIncomeCurrency"` // 买家所付运费币种
	ShippingLimitDate         string         `json:"shippingLimiteDate"`        // 发货截止时间
	TaxCurrency               string         `json:"taxCurrency"`               // 税费币种
	TaxIncome                 float64        `json:"taxIncome"`                 // 税费
	WarehouseIdKey            string         `json:"warehouseIdKey"`            // 通途仓库id key
	WarehouseName             string         `json:"warehouseName"`             // 仓库名称
	WebFinalFee               float64        `json:"webFinalFee"`               // 平台佣金
	WebStoreOrderId           string         `json:"webstoreOrderId"`           // 平台交易号
	WebStoreItemSite          string         `json:"webstore_item_site"`        // 平台站点id
}

type OrderQueryParams struct {
	AccountCode      string `json:"accountCode"`
	BuyerEmail       string `json:"buyerEmail,omitempty"`
	MerchantId       string `json:"merchantId"`
	OrderId          string `json:"orderId,omitempty"`
	OrderStatus      string `json:"orderStatus,omitempty"`
	PageNo           int    `json:"pageNo,omitempty"`
	PageSize         int    `json:"pageSize,omitempty"`
	PayDateFrom      string `json:"payDateFrom,omitempty"`
	PayDateTo        string `json:"payDateTo,omitempty"`
	PlatformCode     string `json:"platformCode,omitempty"`
	RefundedDateFrom string `json:"refundedDateFrom,omitempty"`
	RefundedDateTo   string `json:"refundedDateTo,omitempty"`
	SaleDateFrom     string `json:"saleDateFrom,omitempty"`
	SaleDateTo       string `json:"saleDateTo,omitempty"`
	StoreFlag        string `json:"storeFlag"`
	UpdatedDateFrom  string `json:"updatedDateFrom,omitempty"`
	UpdatedDateTo    string `json:"updatedDateTo,omitempty"`
}

// Orders 订单列表
// https://open.tongtool.com/apiDoc.html#/?docId=f4371e5d65c242a588ebe05872c8c4f8
func (s service) Orders(params OrderQueryParams) (items []Order, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	if !in.StringIn(params.StoreFlag, OrderStoreFlagActive, OrderStoreFlagOneYear, OrderStoreFlagArchived) {
		// ”0”查询活跃表，”1”为查询1年表，”2”为查询归档表，默认为”0”
		// 活跃表：3个月内
		// 1年表：3个月到15个月
		// 归档表：15个月以前
		params.StoreFlag = OrderStoreFlagActive
	}
	if params.OrderId != "" {
		params.AccountCode = ""
	}
	items = make([]Order, 0)
	res := struct {
		result
		Datas struct {
			Array    []Order `json:"array"`
			PageNo   int     `json:"pageNo"`
			PageSize int     `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/ordersQuery")
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

func (s service) Order(orderId string) (item Order, err error) {
	if len(orderId) == 0 {
		err = errors.New("orderId params cannot empty")
		return
	}

	params := OrderQueryParams{
		OrderId:  orderId,
		PageNo:   1,
		PageSize: s.tongTool.QueryDefaultValues.PageSize,
	}

	exists := false
	for {
		items := make([]Order, 0)
		isLastPage := false
		items, isLastPage, err = s.Orders(params)
		if err == nil {
			if len(items) == 0 {
				err = tongtool.ErrNotFound
			} else {
				for _, order := range items {
					if strings.EqualFold(order.OrderIdCode, orderId) {
						exists = true
						item = order
						break
					}
				}
				if exists {
					break
				}
			}
		}
		if isLastPage || exists || err != nil {
			break
		}
		params.PageNo++
	}

	if err == nil && !exists {
		err = tongtool.ErrNotFound
	}

	return
}

// OrderBuyer 订单买家
type OrderBuyer struct {
	BuyerAccount     string `json:"buyerAccount"`     // 买家账号
	BuyerAddress1    string `json:"buyerAddress1"`    // 地址1
	BuyerAddress2    string `json:"buyerAddress2"`    // 地址2
	BuyerAddress3    string `json:"buyerAddress3"`    // 地址3
	BuyerCity        string `json:"buyerCity"`        // 城市
	BuyerCountryCode string `json:"buyerCountryCode"` // 国家(名称或代码)
	BuyerEmail       string `json:"buyerEmail"`       // 买家邮箱
	BuyerMobilePhone string `json:"buyerMobilePhone"` // 手机
	BuyerName        string `json:"buyerName"`        // 买家名称
	BuyerPhone       string `json:"buyerPhone"`       // 电话
	BuyerPostalCode  string `json:"buyerPostalCode"`  // 邮编
	BuyerState       string `json:"buyerState"`       // 州
}

// OrderPayment 订单付款信息
type OrderPayment struct {
	OrderAmount           float64 `json:"orderAmount"`           // 订单金额
	OrderAmountCurrency   string  `json:"orderAmountCurrency"`   // 订单金额币种
	PaymentAccount        string  `json:"paymentAccount"`        // 支付账号
	PaymentDate           string  `json:"paymentDate"`           // 付款时间（yyyy-MM-dd HH:mm:ss）
	PaymentMethod         string  `json:"paymentMethod"`         // 付款方式
	PaymentNotes          string  `json:"paymentNotes"`          // 备注
	PaymentTransactionNum string  `json:"paymentTransactionNum"` // 交易流水号
	RecipientAccount      string  `json:"recipientAccount"`      // 收款账号
	URL                   string  `json:"url"`                   // 相关链接
}

// OrderTransaction 订单交易信息
type OrderTransaction struct {
	GoodsDetailId              string  `json:"goodsDetailId"`              // 货品ID,与SKU二传一即可;如果与SKU都传值了，以这个字段值为准
	GoodsDetailRemark          string  `json:"goodsDetailRemark"`          // 货品备注
	ProductsTotalPrice         float64 `json:"productsTotalPrice"`         // 商品总金额
	ProductsTotalPriceCurrency string  `json:"productsTotalPriceCurrency"` // 商品总金额币种
	Quantity                   int     `json:"quantity"`                   // 数量
	ShipType                   string  `json:"shipType"`                   // 买家选择的运输方式
	ShippingFeeIncome          float64 `json:"shippingFeeIncome"`          // 买家所支付的运费
	ShippingFeeIncomeCurrency  string  `json:"shippingFeeIncomeCurrency"`  // 买家所支付的运费币种
	SKU                        string  `json:"sku"`                        // 商品 sku
}

type CreateOrderRequest struct {
	BuyerInfo               OrderBuyer         `json:"buyerInfo"`               // 买家信息
	Currency                string             `json:"currency"`                // 币种
	InsuranceIncome         float64            `json:"insuranceIncome"`         // 买家支付的保险
	InsuranceIncomeCurrency string             `json:"insuranceIncomeCurrency"` // 买家支付的保险币种
	NeedReturnOrderId       string             `json:"needReturnOrderId"`       // 是否需要返回通途订单ID,0-不需要1-需要 默认0不需要;如果需要返回订单ID那么返回结果集是一个Object:{"orderId":"","saleRecordNum":""},否则返回一个字符串，内容是saleRecordNum
	Notes                   string             `json:"notes"`                   // 买家留言
	OrderCurrency           string             `json:"ordercurrency"`           // 订单币种
	PaymentInfos            []OrderPayment     `json:"paymentInfos"`            // 付款信息
	PlatformCode            string             `json:"platformCode"`            // 订单平台代码
	Remarks                 string             `json:"remarks"`                 // 订单备注,只能新增
	SaleRecordNum           string             `json:"saleRecordNum"`           // 订单号
	SellerAccountCode       string             `json:"sellerAccountCode"`       // 卖家账号代码
	ShippingMethodId        string             `json:"shippingMethodId"`        // 渠道ID
	TaxIncome               float64            `json:"taxIncome"`               // 买家支付的税金
	TaxIncomeCurrency       string             `json:"taxIncomeCurrency"`       // 买家支付的税金币种
	TotalPrice              float64            `json:"totalPrice"`              // 订单总额
	TotalPriceCurrency      string             `json:"totalPriceCurrency"`      // 订单总额币种
	Transactions            []OrderTransaction `json:"transactions"`            // 交易信息
	WarehouseId             string             `json:"warehouseId"`             // 仓库ID
}

// CreateOrder 手工创建订单
// https://open.tongtool.com/apiDoc.html#/?docId=908e49d8bf62487aa870335ef6951567
func (s service) CreateOrder(req CreateOrderRequest) (id, number string, err error) {
	req.NeedReturnOrderId = strings.TrimSpace(req.NeedReturnOrderId)
	if req.NeedReturnOrderId == "" || !in.StringIn(req.NeedReturnOrderId, "1", "0") {
		req.NeedReturnOrderId = "0"
	}
	orderReq := struct {
		MerchantId string             `json:"merchantId"` // 商户ID
		Order      CreateOrderRequest `json:"order"`      // 订单信息
	}{
		MerchantId: s.tongTool.MerchantId,
		Order:      req,
	}
	resNoOrderId := struct {
		result
		Datas string `json:"datas"`
	}{}
	resWithOrderId := struct {
		result
		Datas struct {
			OrderId     string `json:"orderId"`
			OrderNumber string `json:"saleRecordNum"`
		} `json:"datas"`
	}{}
	needOrderId := orderReq.Order.NeedReturnOrderId == "1"
	r := s.tongTool.Client.R()
	if needOrderId {
		r.SetResult(&resWithOrderId)
	} else {
		r.SetResult(&resNoOrderId)
	}
	resp, err := r.SetBody(orderReq).Post("/openapi/tongtool/orderImport")
	if err == nil {
		code := 0
		message := ""
		if resp.IsSuccess() {
			if needOrderId {
				code = resWithOrderId.Code
				message = resWithOrderId.Message
			} else {
				code = resNoOrderId.Code
				message = resNoOrderId.Message
			}
			err = tongtool.ErrorWrap(code, message)
			if err == nil {
				if needOrderId {
					id = resWithOrderId.Datas.OrderId
					number = resWithOrderId.Datas.OrderNumber
				} else {
					number = resNoOrderId.Datas
				}
			}
		} else {
			var res interface{}
			if needOrderId {
				res = resWithOrderId
			} else {
				res = resNoOrderId
			}
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				if needOrderId {
					code = resWithOrderId.Code
					message = resWithOrderId.Message
				} else {
					code = resNoOrderId.Code
					message = resNoOrderId.Message
				}
				err = tongtool.ErrorWrap(code, message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	return
}
