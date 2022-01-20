package erp2

import (
	"context"
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/inx"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
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
	GoodsMatchedSKU      string  `json:"goodsMatchedSku"`
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
	GoodsSKU           string  `json:"goodsSku"`
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
	// 自定义属性
	IsInvalidBoolean   bool `json:"isInvalidBoolean"`   // 是否作废布尔值
	IsSuspendedBoolean bool `json:"isSuspendedBoolean"` // 是否需要人工审核布尔值
}

type OrderQueryParams struct {
	Paging
	AccountCode      string `json:"accountCode"`                // ERP系统中，基础设置->账号管理 列表中的代码
	BuyerEmail       string `json:"buyerEmail,omitempty"`       // 买家邮箱
	MerchantId       string `json:"merchantId"`                 // 商户ID
	OrderId          string `json:"orderId,omitempty"`          // 订单号
	OrderStatus      string `json:"orderStatus,omitempty"`      // 订单状态 waitPacking/等待配货 ,waitPrinting/等待打印,waitingDespatching/等待发货 ,despatched/已发货,unpaid/未付款,payed/已付款,
	PayDateFrom      string `json:"payDateFrom,omitempty"`      // 付款起始时间
	PayDateTo        string `json:"payDateTo,omitempty"`        // 付款结束时间
	PlatformCode     string `json:"platformCode,omitempty"`     // 通途中平台代码
	RefundedDateFrom string `json:"refundedDateFrom,omitempty"` // 退款起始时间
	RefundedDateTo   string `json:"refundedDateTo,omitempty"`   // 退款结束时间
	SaleDateFrom     string `json:"saleDateFrom,omitempty"`     // 销售起始时间
	SaleDateTo       string `json:"saleDateTo,omitempty"`       // 销售结束时间
	StoreFlag        string `json:"storeFlag"`                  // 是否需要查询1年表和归档表数据（根据时间参数或者全量查询订单的时候使用该参数，”0”查询活跃表，”1”为查询1年表，”2”为查询归档表，默认为”0”）
	UpdatedDateFrom  string `json:"updatedDateFrom,omitempty"`  // 更新开始时间
	UpdatedDateTo    string `json:"updatedDateTo,omitempty"`    // 更新结束时间
}

// StoreCountryCode 获取订单店铺所在国家代码
func (o Order) StoreCountryCode() string {
	code := ""
	if o.WebStoreItemSite != "" {
		switch o.WebStoreItemSite {
		case "100002":
			code = constant.CountryCodeAmerica
		case "100003":
			code = constant.CountryCodeUK
		case "100004":
			code = constant.CountryCodeCanada
		case "100005":
			code = constant.CountryCodeGermany
		case "100006":
			code = constant.CountryCodeSpain
		case "100007":
			code = constant.CountryCodeFrance
		case "100008":
			code = constant.CountryCodeItaly
		case "100009":
			code = constant.CountryCodeJapan
		case "100010":
			code = constant.CountryCodeMexico
		case "100011":
			code = constant.CountryCodeAustralian
		case "100012":
			code = constant.CountryCodeIndia
		case "100013":
			code = constant.CountryCodeUnitedArabEmirates
		case "100014":
			code = constant.CountryCodeTurkey
		case "100015":
			code = constant.CountryCodeSingapore
		case "100016":
			code = constant.CountryCodeNetherlands
		case "100017":
			code = constant.CountryCodeBrazil
		case "100018":
			code = constant.CountryCodeSaudiArabia
		case "100019":
			code = constant.CountryCodeSweden
		case "100020":
			code = constant.CountryCodePoland
		}
	} else if o.BuyerCountry != "" {
		// Todo 美国的买家买的加拿大站点的怎么办？
		country := strings.TrimSpace(o.BuyerCountry)
		if country != "" {
			if inx.StringIn(country, constant.CountryCodeAmerica, constant.CountryCodeCanada, constant.CountryCodeGermany, constant.CountryCodeUK, constant.CountryCodeFrance, constant.CountryCodeSpain, constant.CountryCodeItaly, constant.CountryCodeJapan, constant.CountryCodeMexico, constant.CountryCodeAustralian, constant.CountryCodeIndia, constant.CountryCodeUnitedArabEmirates, constant.CountryCodeTurkey, constant.CountryCodeSingapore, constant.CountryCodeNetherlands, constant.CountryCodeBrazil, constant.CountryCodeSaudiArabia, constant.CountryCodeSweden, constant.CountryCodePoland) {
				code = strings.ToUpper(country)
			}
		}
	}

	return code
}

// Orders 订单列表
// https://open.tongtool.com/apiDoc.html#/?docId=f4371e5d65c242a588ebe05872c8c4f8
func (s service) Orders(params OrderQueryParams) (items []Order, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
	if !inx.StringIn(params.StoreFlag, OrderStoreFlagActive, OrderStoreFlagOneYear, OrderStoreFlagArchived) {
		// ”0”查询活跃表，”1”为查询1年表，”2”为查询归档表，默认为”0”
		// 活跃表：3个月内
		// 1年表：3个月到15个月
		// 归档表：15个月以前
		params.StoreFlag = OrderStoreFlagActive
	}
	if params.OrderId != "" {
		params.AccountCode = ""
	}
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = keyx.Generate(params)
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
				for i, item := range items {
					item.IsInvalidBoolean = !inx.StringIn(item.IsInvalid, "0", "", "null")
					item.IsSuspendedBoolean = item.IsSuspended == "1"
					items[i] = item
				}
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

	if err == nil && s.tongTool.EnableCache && len(items) > 0 {
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

func (s service) Order(orderId string) (item Order, err error) {
	if len(orderId) == 0 {
		err = errors.New("orderId params cannot empty")
		return
	}

	params := OrderQueryParams{
		OrderId: orderId,
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
	Remarks                 []string           `json:"remarks"`                 // 订单备注,只能新增
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

func (m CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.BuyerInfo, validation.Required.Error("买家信息不能为空"), validation.By(func(value interface{}) error {
			buyer, ok := value.(OrderBuyer)
			if !ok {
				return errors.New("无效的买家信息")
			}
			return validation.ValidateStruct(&buyer,
				validation.Field(&buyer.BuyerName, validation.Required.Error("买家名称不能为空")),
				validation.Field(&buyer.BuyerAddress1, validation.Required.Error("买家地址1不能为空")),
				validation.Field(&buyer.BuyerCity, validation.Required.Error("买家城市不能为空")),
				validation.Field(&buyer.BuyerCountryCode, validation.Required.Error("买家国家不能为空")),
			)
		})),
		validation.Field(&m.WarehouseId, validation.Required.Error("仓库 ID 不能为空")),
		validation.Field(&m.NeedReturnOrderId, validation.Required.Error("请填写订单返回值设置"), validation.In("0", "1").Error("无效的订单返回值设置")),
		validation.Field(&m.Transactions, validation.Required.Error("交易信息不能为空"), validation.Each(validation.WithContext(func(ctx context.Context, value interface{}) error {
			item, ok := value.(OrderTransaction)
			if !ok {
				return errors.New("无效的交易信息")
			}
			return validation.ValidateStruct(&item,
				validation.Field(&item.SKU, validation.When(item.GoodsDetailId == "", validation.Required.Error("货品 ID 与 SKU 必传其中一个"))),
				validation.Field(&item.Quantity, validation.Min(1).Error("数量必须大于 1")),
			)
		}))),
	)
}

// CreateOrder 手工创建订单
// https://open.tongtool.com/apiDoc.html#/?docId=908e49d8bf62487aa870335ef6951567
// 在 err 为 nil 的情况下，orderNumber 一定会有值返回，而 orderId 是否有值则取决于是否在查询参数中传递 NeedReturnOrderId 为 1 的值
func (s service) CreateOrder(req CreateOrderRequest) (orderId, orderNumber string, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	orderReq := struct {
		MerchantId string             `json:"merchantId"` // 商户ID
		Order      CreateOrderRequest `json:"order"`      // 订单信息
	}{
		MerchantId: s.tongTool.MerchantId,
		Order:      req,
	}
	res := struct {
		result
		Datas interface{} `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(orderReq).
		SetResult(&res).
		Post("/openapi/tongtool/orderImport")
	if err == nil {
		if resp.IsSuccess() {
			err = tongtool.ErrorWrap(res.Code, res.Message)
			if err == nil {
				if orderReq.Order.NeedReturnOrderId == "1" {
					withOrderIdValue := struct {
						OrderId     string `json:"orderId"`
						OrderNumber string `json:"saleRecordNum"`
					}{}
					var b []byte
					if b, err = json.Marshal(res.Datas); err == nil {
						if err = json.Unmarshal(b, &withOrderIdValue); err == nil {
							orderId = withOrderIdValue.OrderId
							orderNumber = withOrderIdValue.OrderNumber
						}
					}
				} else {
					orderNumber = res.Datas.(string)
				}
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

// 更新订单处理（未配货前可用）

type UpdateOrderTransaction struct {
	GoodsDetailId  string `json:"goodsDetailId"`  // 货品ID与订单详情ID二者必填其一
	OrderDetailsId string `json:"orderDetailsId"` // 订单详情ID,此参数值来自订单查询返回,此参数有值代表是需要更新货品数量或者删除货品(要看quantity参数值)，此参数有值同时会清空原有核查结果，需要重新核查，此参数没有值但goodsDetailId有值代表是需要新增货品
	Quantity       int    `json:"quantity"`       // 数量,等于0表示删除当前货品
}

// UpdateOrderRequest 订单更新请求
type UpdateOrderRequest struct {
	BuyerInfo        OrderBuyer               `json:"buyerInfo"`                  // 买家信息
	Transactions     []UpdateOrderTransaction `json:"transactions"`               // 交易记录信息,删除货品需要传对应的记录并数量传0
	MerchantId       string                   `json:"merchantId"`                 // 商户ID
	OrderId          string                   `json:"orderId"`                    // 通途订单ID
	Remarks          []string                 `json:"remarks,omitempty"`          // 订单备注,只能新增
	ShippingMethodId string                   `json:"shippingMethodId,omitempty"` // 渠道ID
	WarehouseId      string                   `json:"warehouseId,omitempty"`      // 仓库ID
}

func (m UpdateOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderId, validation.Required.Error("订单 ID 不能为空")),
	)
}

// UpdateOrder 更新订单
// https://open.tongtool.com/apiDoc.html#/?docId=3e0d01bfe01441aa8e2071c2c88cc9fb
func (s service) UpdateOrder(req UpdateOrderRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	req.MerchantId = s.tongTool.MerchantId

	res := struct {
		result
		Datas string `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/orderUpdate")
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

// 作废订单处理

type CancelOrderRequest struct {
	MerchantId  string   `json:"merchantId"`  // 商戶 ID
	OrderIdKeys []string `json:"orderIdKeys"` // 通途订单id key
}

func (m CancelOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderIdKeys, validation.Required.Error("订单 ID 不能为空")),
	)
}

type OrderCancelResult struct {
	OrderId string `json:"order_id"` // OrderId
	Result  string `json:"result"`   // 结果
}

// CancelOrder 作废订单
// https://open.tongtool.com/apiDoc.html#/?docId=9ba0ea5da90740f28a0345aa1990c007
func (s service) CancelOrder(req CancelOrderRequest) (results []OrderCancelResult, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		result
		Datas struct {
			Array []struct {
				OrderId string `json:"order_id"`
				Result  string `json:"result"`
			} `json:"array"`
		} `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/orderCancel")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				for _, item := range res.Datas.Array {
					results = append(results, OrderCancelResult{
						OrderId: item.OrderId,
						Result:  strings.TrimSpace(item.Result),
					})
				}
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
