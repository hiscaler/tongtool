package erp2

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hiscaler/gox/filex"
	"github.com/hiscaler/gox/inx"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 订单状态
const (
	OrderStatusWaitPacking        = "waitPacking"        // 等待配货
	OrderStatusWaitPrinting       = "waitPrinting"       // 等待打印
	OrderStatusWaitingDespatching = "waitingDespatching" // 等待发货
	OrderStatusDespatched         = "despatched"         // 已发货
	OrderStatusUnpaid             = "unpaid"             // 未付款
	OrderStatusPaid               = "payed"              // 已付款
)

// 查询对象
const (
	OrderStoreFlagActive   = "0" // 活跃表（3 个月内）
	OrderStoreFlagOneYear  = "1" // 一年表（3 个月到 15 个月）
	OrderStoreFlagArchived = "2" // 归档表（15 个月以前）
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
	TrackingNumberStatus string `json:"trackingNumberStatus"` // 物流跟踪号获取状态(00：未就绪、01：就绪、02：处理中、03：处理成功、04：处理失败)
	TrackingNumberTime   int    `json:"trackingNumberTime"`   // 物流跟踪号获取时间
}

// PlatformGoodsInfo 平台商品信息
type PlatformGoodsInfo struct {
	ProductId             string `json:"product_id"`          // 商品顺序号
	Quantity              int    `json:"quantity"`            // 原始 SKU 数量
	WebTransactionId      string `json:"webTransactionId"`    // 平台订单产品交易号
	WebStoreCustomLabel   string `json:"webstoreCustomLabel"` // 原始 SKU
	WebStoreItemId        string `json:"webstoreItemId"`      // 平台订单产品 ItemId
	WebStoreSKU           string `json:"webstoreSku"`         // 通途 SKU
	CustomizedURL         string `json:"customizedUrl"`       // 定制信息下载地址
	CustomizedInformation struct {
		Ok            bool              `json:"ok"`            // 是否处理完毕
		Error         string            `json:"error"`         // 处理时所产生的错误信息
		SnapshotImage string            `json:"snapshotImage"` // Image is base64 format
		Text          string            `json:"text"`          // All text use \n split
		Images        map[string]string `json:"images"`        // Image is base64 format
	} `json:"customizedInformation"` // 定制信息
}

// TongToolGoodsInfo 通途商品信息
type TongToolGoodsInfo struct {
	GoodsAverageCost     float64 `json:"goodsAverageCost"`     // 货品平均成本
	GoodsCurrentCost     float64 `json:"goodsCurrentCost"`     // 货品成本（最新成本）
	GoodsImageGroupId    string  `json:"goodsImageGroupId"`    // 商品图片
	GoodsPackagingCost   float64 `json:"goodsPackagingCost"`   // 货品包装成本
	GoodsPackagingWeight float64 `json:"goodsPackagingWeight"` // 货品包装重量（克）
	GoodsSKU             string  `json:"goodsSku"`             // 货品 SKU
	GoodsTitle           string  `json:"goodsTitle"`           // 商品规格
	GoodsWeight          float64 `json:"goodsWeight"`          // 货品重量（克）
	PackageHeight        float64 `json:"packageHeight"`        // 包裹尺寸（高cm）
	PackageLength        float64 `json:"packageLength"`        // 包裹尺寸（长cm）
	PackageWidth         float64 `json:"packageWidth"`         // 包裹尺寸（宽cm）
	PackagingCost        float64 `json:"packagingCost"`        // 货品包装成本
	PackagingWeight      float64 `json:"packagingWeight"`      // 商品包装重量（克）
	ProductAverageCost   float64 `json:"productAverageCost"`   // 商品平均成本
	ProductCurrentCost   float64 `json:"productCurrentCost"`   // 商品成本
	ProductHeight        float64 `json:"productHeight"`        // 商品尺寸高（cm）
	ProductLength        float64 `json:"productLength"`        // 商品尺寸长（cm）
	ProductName          string  `json:"productName"`          // 商品名称
	ProductWeight        float64 `json:"productWeight"`        // 商品重量（克）
	ProductWidth         float64 `json:"productWidth"`         // 商品尺寸宽（cm）
	Quantity             int     `json:"quantity"`             // 货品数量
	WotId                string  `json:"wotId"`                // 平台交易编号
}

type GoodsInfo struct {
	PlatformGoodsInfoList []PlatformGoodsInfo `json:"platformGoodsInfoList"` // 平台商品信息列表
	TongToolGoodsInfoList []TongToolGoodsInfo `json:"tongToolGoodsInfoList"` // 通途商品信息列表
}

// Order 通途订单
type Order struct {
	ActualTotalPrice          float64        `json:"actualTotalPrice"`          // 实付金额
	AssignStockCompleteTime   string         `json:"assignstockCompleteTime"`   // 配货时间
	BuyerAccountId            string         `json:"buyerAccountId"`            // 买家 ID
	BuyerCity                 string         `json:"buyerCity"`                 // 买家城市
	BuyerCountry              string         `json:"buyerCountry"`              // 买家国家
	BuyerEmail                string         `json:"buyerEmail"`                // 买家邮箱
	BuyerMobile               string         `json:"buyerMobile"`               // 买家手机
	BuyerName                 string         `json:"buyerName"`                 // 买家名称
	BuyerPassportCode         string         `json:"buyerPassportCode"`         // 收件人识别码（护照等）
	BuyerPhone                string         `json:"buyerPhone"`                // 买家电话
	BuyerState                string         `json:"buyerState"`                // 买家省份
	Carrier                   string         `json:"carrier"`                   // 上传物流的carrier
	CarrierType               string         `json:"carrierType"`               // 物流商类型（0：通途API对接、1：通途Excel文件导出、2：通途离线生成跟踪号、3：无对接、4：自定义Excel对接）
	CarrierURL                string         `json:"carrierUrl"`                // 物流网络地址
	DespatchCompleteTime      string         `json:"despatchCompleteTime"`      // 订单发货完成时间
	DispatchTypeName          string         `json:"dispathTypeName"`           // 邮寄方式名称
	DownloadTime              string         `json:"downloadTime"`              // 订单首次下载或创建时间
	EarliestDeliveryDate      string         `json:"earliestDeliveryDate"`      // 亚马逊订单预计最早送达时间
	EbayNotes                 string         `json:"ebayNotes"`                 // 订单备注
	EbaySiteEnName            string         `json:"ebaySiteEnName"`            // 站点
	FirstTariff               float64        `json:"firstTariff"`               // 头程运费
	GoodsInfo                 GoodsInfo      `json:"goodsInfo"`                 // 订单商品信息
	InsuranceIncome           float64        `json:"insuranceIncome"`           // 买家所付保费
	InsuranceIncomeCurrency   string         `json:"insuranceIncomeCurrency"`   // 买家所付保费币种
	IsInvalid                 string         `json:"isInvalid"`                 // 是否作废（0,'',null：未作废、1：手工作废、2：订单任务下载永久作废、3：拆分单主单作废、4：拆分单子单作废）
	IsSuspended               string         `json:"isSuspended"`               // 是否需要人工审核（1：需要人工审核、0或null：不需要）
	MerchantCarrierShortname  string         `json:"merchantCarrierShortname"`  // 承运人简称
	OrderAmount               float64        `json:"orderAmount"`               // 订单总金额（商品金额+运费+保费）
	OrderAmountCurrency       string         `json:"orderAmountCurrency"`       // 订单金额币种
	OrderDetails              []OrderDetail  `json:"orderDetails"`              // 订单明细
	OrderIdCode               string         `json:"orderIdCode"`               // 通途订单号
	OrderIdKey                string         `json:"orderIdKey"`                // 通途订单 ID Key
	OrderStatus               string         `json:"orderStatus"`               // 订单状态（waitPacking：等待配货、waitPrinting：等待打印、waitingDespatching：等待发货、despatched：已发货）
	PackageInfoList           []OrderPackage `json:"packageInfoList"`           // 订单包裹信息
	PaidTime                  string         `json:"paidTime"`                  // 订单付款完成时间
	ParentOrderId             string         `json:"parentOrderId"`             // 父订单号
	PlatformCode              string         `json:"platformCode"`              // 通途中平台代码
	PlatformFee               float64        `json:"platformFee"`               // 平台手续费
	PostalCode                string         `json:"postalCode"`                // 买家邮编
	PrintCompleteTime         string         `json:"printCompleteTime"`         // 订单打印完成时间
	ProductsTotalCurrency     string         `json:"productsTotalCurrency"`     // 金额小计币种
	ProductsTotalPrice        float64        `json:"productsTotalPrice"`        // 金额小计（只商品金额）
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
	UpdatedTime               string         `json:"updatedTime"`               // 订单更新时间
	WarehouseIdKey            string         `json:"warehouseIdKey"`            // 通途仓库 ID Key
	WarehouseName             string         `json:"warehouseName"`             // 仓库名称
	WebFinalFee               float64        `json:"webFinalFee"`               // 平台佣金
	WebStoreOrderId           string         `json:"webstoreOrderId"`           // 平台交易号
	WebStoreItemSite          string         `json:"webstore_item_site"`        // 平台站点 ID
	// 自定义属性
	IsInvalidBoolean   bool `json:"isInvalidBoolean"`   // 是否作废布尔值
	IsSuspendedBoolean bool `json:"isSuspendedBoolean"` // 是否需要人工审核布尔值
}

// StoreCountryCode 获取订单店铺所在国家代码
func (o Order) StoreCountryCode() string {
	code := getSiteCountryCodeById(o.WebStoreItemSite)
	if code == "" && o.BuyerCountry != "" {
		// Todo 美国的买家买的加拿大站点的怎么办？或者国际站的也会判断不正确
		country := strings.TrimSpace(o.BuyerCountry)
		if country != "" {
			if inx.StringIn(country, constant.CountryCodeAmerica, constant.CountryCodeCanada, constant.CountryCodeGermany, constant.CountryCodeUnitedKingdom, constant.CountryCodeFrance, constant.CountryCodeSpain, constant.CountryCodeItaly, constant.CountryCodeJapan, constant.CountryCodeMexico, constant.CountryCodeAustralian, constant.CountryCodeIndia, constant.CountryCodeUnitedArabEmirates, constant.CountryCodeTurkey, constant.CountryCodeSingapore, constant.CountryCodeNetherlands, constant.CountryCodeBrazil, constant.CountryCodeSaudiArabia, constant.CountryCodeSweden, constant.CountryCodePoland, constant.CountryCodeChina) {
				code = strings.ToUpper(country)
			}
		}
	}
	return code
}

// Amount 获取订单金额数据
func (o Order) Amount(exchangeRates map[string]float64, precision int32, shippingFee, otherFee float64) *OrderAmount {
	oa := NewOrderAmount(o, exchangeRates, precision, shippingFee, otherFee)
	return oa
}

type OrdersQueryParams struct {
	Paging
	AccountCode      string `json:"accountCode"`                // ERP系统中，基础设置->账号管理 列表中的代码
	BuyerEmail       string `json:"buyerEmail,omitempty"`       // 买家邮箱
	MerchantId       string `json:"merchantId"`                 // 商户 ID
	OrderId          string `json:"orderId,omitempty"`          // 订单号
	OrderStatus      string `json:"orderStatus,omitempty"`      // 订单状态（waitPacking：等待配货、waitPrinting：等待打印、waitingDespatching：等待发货、despatched：已发货、unpaid：未付款、payed：已付款）
	PayDateFrom      string `json:"payDateFrom,omitempty"`      // 付款起始时间
	PayDateTo        string `json:"payDateTo,omitempty"`        // 付款结束时间
	PlatformCode     string `json:"platformCode,omitempty"`     // 通途中平台代码
	RefundedDateFrom string `json:"refundedDateFrom,omitempty"` // 退款起始时间
	RefundedDateTo   string `json:"refundedDateTo,omitempty"`   // 退款结束时间
	SaleDateFrom     string `json:"saleDateFrom,omitempty"`     // 销售起始时间
	SaleDateTo       string `json:"saleDateTo,omitempty"`       // 销售结束时间
	StoreFlag        string `json:"storeFlag"`                  // 是否需要查询 1 年表或归档表数据（根据时间参数或者全量查询订单的时候使用该参数，0：活跃表、1：一年表、2：归档表，默认为 0）
	UpdatedDateFrom  string `json:"updatedDateFrom,omitempty"`  // 更新开始时间
	UpdatedDateTo    string `json:"updatedDateTo,omitempty"`    // 更新结束时间
}

func (m OrdersQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.BuyerEmail, validation.When(m.BuyerEmail != "", is.EmailFormat.Error("无效的邮箱格式"))),
		validation.Field(&m.StoreFlag,
			validation.When(m.StoreFlag != "", validation.In(OrderStoreFlagActive, OrderStoreFlagOneYear, OrderStoreFlagArchived).Error("无效的查询范围")),
		),
		validation.Field(&m.OrderStatus, validation.When(m.OrderStatus != "", validation.In(OrderStatusWaitPacking, OrderStatusWaitPrinting, OrderStatusWaitingDespatching, OrderStatusDespatched, OrderStatusUnpaid, OrderStatusPaid).Error("无效的订单状态"))),
		validation.Field(&m.PayDateFrom, validation.When(m.PayDateTo != "", validation.Date(constant.DatetimeFormat).Error("无效的付款起始时间"))),
		validation.Field(&m.PayDateTo, validation.When(m.PayDateTo != "",
			validation.Date(constant.DatetimeFormat).Error("无效的付款结束时间"),
			validation.By(func(value interface{}) error {
				var err error
				var fromDate, toDate time.Time
				t, _ := value.(string)
				if toDate, err = time.Parse(constant.DatetimeFormat, t); err != nil {
					return err
				}
				if fromDate, err = time.Parse(constant.DatetimeFormat, m.PayDateFrom); err != nil {
					return err
				}
				if toDate.Before(fromDate) {
					return fmt.Errorf("付款结束时间 %s 不能小于开始时间 %s", m.PayDateTo, m.PayDateFrom)
				}
				return nil
			}),
		)),
		validation.Field(&m.RefundedDateFrom, validation.When(m.RefundedDateFrom != "", validation.Date(constant.DatetimeFormat).Error("无效的退款起始时间"))),
		validation.Field(&m.RefundedDateTo, validation.When(m.RefundedDateTo != "",
			validation.Date(constant.DatetimeFormat).Error("无效的退款结束时间"),
			validation.By(func(value interface{}) error {
				var err error
				var fromDate, toDate time.Time
				t, _ := value.(string)
				if toDate, err = time.Parse(constant.DatetimeFormat, t); err != nil {
					return err
				}
				if fromDate, err = time.Parse(constant.DatetimeFormat, m.RefundedDateFrom); err != nil {
					return err
				}
				if toDate.Before(fromDate) {
					return fmt.Errorf("退款结束时间 %s 不能小于开始时间 %s", m.RefundedDateTo, m.RefundedDateFrom)
				}
				return nil
			}),
		)),
		validation.Field(&m.SaleDateFrom, validation.When(m.SaleDateFrom != "", validation.Date(constant.DatetimeFormat).Error("无效的销售起始时间"))),
		validation.Field(&m.SaleDateTo, validation.When(m.SaleDateTo != "",
			validation.Date(constant.DatetimeFormat).Error("无效的销售结束时间"),
			validation.By(func(value interface{}) error {
				var err error
				var fromDate, toDate time.Time
				t, _ := value.(string)
				if toDate, err = time.Parse(constant.DatetimeFormat, t); err != nil {
					return err
				}
				if fromDate, err = time.Parse(constant.DatetimeFormat, m.SaleDateFrom); err != nil {
					return err
				}
				if toDate.Before(fromDate) {
					return fmt.Errorf("销售结束时间 %s 不能小于开始时间 %s", m.SaleDateTo, m.SaleDateFrom)
				}
				return nil
			}),
		)),
		validation.Field(&m.UpdatedDateFrom, validation.When(m.UpdatedDateFrom != "", validation.Date(constant.DatetimeFormat).Error("无效的更新起始时间"))),
		validation.Field(&m.UpdatedDateTo, validation.When(m.UpdatedDateTo != "",
			validation.Date(constant.DatetimeFormat).Error("无效的更新结束时间"),
			validation.By(func(value interface{}) error {
				var err error
				var fromDate, toDate time.Time
				t, _ := value.(string)
				if toDate, err = time.Parse(constant.DatetimeFormat, t); err != nil {
					return err
				}
				if fromDate, err = time.Parse(constant.DatetimeFormat, m.UpdatedDateFrom); err != nil {
					return err
				}
				if toDate.Before(fromDate) {
					return fmt.Errorf("更新结束时间 %s 不能小于开始时间 %s", m.UpdatedDateTo, m.UpdatedDateFrom)
				}
				return nil
			}),
		)),
	)
}

func download(url, filename, dir string) (path string, err error) {
	path = filepath.Join(dir, filename+".zip")
	if filex.Exists(path) {
		// Return it if zip file exists.
		return
	}

	file, err := os.Create(path)
	if err != nil {
		return
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	defer file.Close()
	return
}

// Orders 订单列表
// https://open.tongtool.com/apiDoc.html#/?docId=f4371e5d65c242a588ebe05872c8c4f8
func (s service) Orders(params OrdersQueryParams) (items []Order, isLastPage bool, err error) {
	if params.StoreFlag == "" {
		params.StoreFlag = OrderStoreFlagActive
	}
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
	if params.OrderId != "" {
		params.AccountCode = ""
	}
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = keyx.Generate(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = jsoniter.Unmarshal(b, &items); e == nil {
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
	res := struct {
		tongtool.Response
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
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
			parser := NewAmazonCustomizationInformationParser()
			for i := range items {
				for _, detail := range items[i].OrderDetails {
					for ii, gf := range items[i].GoodsInfo.PlatformGoodsInfoList {
						if gf.WebStoreItemId == detail.WebStoreItemId {
							if gf.CustomizedURL != "" {
								var zipFile string
								zipFile, err = download(gf.CustomizedURL, fmt.Sprintf("%s_%s", items[i].OrderIdCode, detail.WebStoreItemId), s.tongTool.GetAssetSaveDir())
								if err != nil {
									items[i].GoodsInfo.PlatformGoodsInfoList[ii].CustomizedInformation.Error = err.Error()
									return
								}
								_, err = parser.Reset().SetZipFile(zipFile).Parse()
								if err != nil {
									items[i].GoodsInfo.PlatformGoodsInfoList[ii].CustomizedInformation.Error = err.Error()
									return
								}
								items[i].GoodsInfo.PlatformGoodsInfoList[ii].CustomizedInformation.Ok = true
								items[i].GoodsInfo.PlatformGoodsInfoList[ii].CustomizedInformation.SnapshotImage = parser.SnapshotImage
								items[i].GoodsInfo.PlatformGoodsInfoList[ii].CustomizedInformation.Text = parser.Text
								items[i].GoodsInfo.PlatformGoodsInfoList[ii].CustomizedInformation.Images = parser.Images
							}
							break
						}
					}
				}

				items[i].IsInvalidBoolean = !inx.StringIn(items[i].IsInvalid, "0", "", "null")
				items[i].IsSuspendedBoolean = items[i].IsSuspended == "1"
			}
			isLastPage = len(items) < params.PageSize
		}
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	if err != nil {
		return
	}

	if s.tongTool.EnableCache && len(items) > 0 {
		if b, e := jsoniter.Marshal(&items); e == nil {
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

// Order 根据订单号获取订单信息
func (s service) Order(orderId string) (item Order, exists bool, err error) {
	err = validation.Validate(orderId, validation.Required.Error("orderId 参数不能为空"))
	if err != nil {
		return
	}

	params := OrdersQueryParams{OrderId: orderId}
	params.PageNo = 1
	for {
		var items []Order
		isLastPage := false
		items, isLastPage, err = s.Orders(params)
		if err == nil {
			if len(items) == 0 {
				err = tongtool.ErrNotFound
			} else {
				for i := range items {
					if strings.EqualFold(items[i].OrderIdCode, orderId) {
						exists = true
						item = items[i]
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
	BuyerCountryCode string `json:"buyerCountryCode"` // 国家（名称或代码）
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
	GoodsDetailId              string  `json:"goodsDetailId"`              // 货品 ID（与 SKU 二传一即可，如果与 SKU 都传值了，以这个字段值为准）
	GoodsDetailRemark          string  `json:"goodsDetailRemark"`          // 货品备注
	ProductsTotalPrice         float64 `json:"productsTotalPrice"`         // 商品总金额
	ProductsTotalPriceCurrency string  `json:"productsTotalPriceCurrency"` // 商品总金额币种
	Quantity                   int     `json:"quantity"`                   // 数量
	ShipType                   string  `json:"shipType"`                   // 买家选择的运输方式
	ShippingFeeIncome          float64 `json:"shippingFeeIncome"`          // 买家所支付的运费
	ShippingFeeIncomeCurrency  string  `json:"shippingFeeIncomeCurrency"`  // 买家所支付的运费币种
	SKU                        string  `json:"sku"`                        // 商品 SKU
}

type CreateOrderRequest struct {
	BuyerInfo               OrderBuyer         `json:"buyerInfo"`               // 买家信息
	Currency                string             `json:"currency"`                // 币种
	InsuranceIncome         float64            `json:"insuranceIncome"`         // 买家支付的保险
	InsuranceIncomeCurrency string             `json:"insuranceIncomeCurrency"` // 买家支付的保险币种
	NeedReturnOrderId       string             `json:"needReturnOrderId"`       // 是否需要返回通途订单ID（0：不需要、1：需要）默认 0 不需要；如果需要返回订单 ID 那么返回结果集是一个Object:{"orderId":"","saleRecordNum":""},否则返回一个字符串，内容是saleRecordNum
	Notes                   string             `json:"notes"`                   // 买家留言
	OrderCurrency           string             `json:"ordercurrency"`           // 订单币种
	PaymentInfos            []OrderPayment     `json:"paymentInfos"`            // 付款信息
	PlatformCode            string             `json:"platformCode"`            // 订单平台代码
	Remarks                 []string           `json:"remarks"`                 // 订单备注,只能新增
	SaleRecordNum           string             `json:"saleRecordNum"`           // 订单号
	SellerAccountCode       string             `json:"sellerAccountCode"`       // 卖家账号代码
	ShippingMethodId        string             `json:"shippingMethodId"`        // 渠道 ID
	TaxIncome               float64            `json:"taxIncome"`               // 买家支付的税金
	TaxIncomeCurrency       string             `json:"taxIncomeCurrency"`       // 买家支付的税金币种
	TotalPrice              float64            `json:"totalPrice"`              // 订单总额
	TotalPriceCurrency      string             `json:"totalPriceCurrency"`      // 订单总额币种
	Transactions            []OrderTransaction `json:"transactions"`            // 交易信息
	WarehouseId             string             `json:"warehouseId"`             // 仓库 ID
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
		validation.Field(&m.NeedReturnOrderId,
			validation.Required.Error("请填写订单返回值设置"),
			validation.In("0", "1").Error("无效的订单返回值设置"),
		),
		validation.Field(&m.Transactions,
			validation.Required.Error("交易信息不能为空"),
			validation.Each(validation.WithContext(func(ctx context.Context, value interface{}) error {
				item, ok := value.(OrderTransaction)
				if !ok {
					return errors.New("无效的交易信息")
				}
				return validation.ValidateStruct(&item,
					validation.Field(&item.SKU, validation.When(item.GoodsDetailId == "", validation.Required.Error("货品 ID 与 SKU 必传其中一个"))),
					validation.Field(&item.GoodsDetailRemark, validation.When(item.SKU == "", validation.Required.Error("货品 ID 与 SKU 必传其中一个"))),
					validation.Field(&item.Quantity, validation.Min(1).Error("数量不能小于 {{.threshold}}")),
				)
			})),
		),
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
		MerchantId string             `json:"merchantId"` // 商户 ID
		Order      CreateOrderRequest `json:"order"`      // 订单信息
	}{
		MerchantId: s.tongTool.MerchantId,
		Order:      req,
	}
	res := struct {
		tongtool.Response
		Datas interface{} `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(orderReq).
		SetResult(&res).
		Post("/openapi/tongtool/orderImport")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			if orderReq.Order.NeedReturnOrderId == "1" {
				withOrderIdValue := struct {
					OrderId     string `json:"orderId"`
					OrderNumber string `json:"saleRecordNum"`
				}{}
				var b []byte
				if b, err = jsoniter.Marshal(res.Datas); err == nil {
					if err = jsoniter.Unmarshal(b, &withOrderIdValue); err == nil {
						orderId = withOrderIdValue.OrderId
						orderNumber = withOrderIdValue.OrderNumber
					}
				}
			} else {
				orderNumber = res.Datas.(string)
			}
		}
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}

// 更新订单处理（未配货前可用）

type UpdateOrderTransaction struct {
	GoodsDetailId  string `json:"goodsDetailId"`  // 货品 ID 与订单详情 ID 二者必填其一
	OrderDetailsId string `json:"orderDetailsId"` // 订单详情 ID，此参数值来自订单查询返回，此参数有值代表是需要更新货品数量或者删除货品（要看quantity 参数值），此参数有值同时会清空原有核查结果，需要重新核查，此参数没有值但 goodsDetailId 有值代表是需要新增货品
	Quantity       int    `json:"quantity"`       // 数量（等于 0 表示删除当前货品）
}

// UpdateOrderRequest 订单更新请求
type UpdateOrderRequest struct {
	BuyerInfo        OrderBuyer               `json:"buyerInfo,omitempty"`        // 买家信息
	Transactions     []UpdateOrderTransaction `json:"transactions,omitempty"`     // 交易记录信息，删除货品需要传对应的记录并数量传 0
	MerchantId       string                   `json:"merchantId"`                 // 商户 ID
	OrderId          string                   `json:"orderId"`                    // 通途订单 ID
	Remarks          []string                 `json:"remarks,omitempty"`          // 订单备注（只能新增）
	ShippingMethodId string                   `json:"shippingMethodId,omitempty"` // 渠道 ID
	WarehouseId      string                   `json:"warehouseId,omitempty"`      // 仓库 ID
}

func (m UpdateOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderId, validation.Required.Error("订单 ID 不能为空")),
		validation.Field(&m.Transactions, validation.When(len(m.Transactions) > 0, validation.Each(validation.WithContext(func(ctx context.Context, value interface{}) error {
			if transaction, ok := value.(UpdateOrderTransaction); !ok {
				return errors.New("无效的交易记录信息")
			} else {
				return validation.ValidateStruct(&transaction,
					validation.Field(&transaction.GoodsDetailId, validation.When(transaction.OrderDetailsId == "", validation.Required.Error("货品 ID 与订单详情 ID 二者必填其一"))),
					validation.Field(&transaction.OrderDetailsId, validation.When(transaction.GoodsDetailId == "", validation.Required.Error("货品 ID 与订单详情 ID 二者必填其一"))),
					validation.Field(&transaction.Quantity, validation.Min(0).Error("数量不能小于 {{.threshold}}")),
				)
			}
		})))),
	)
}

// UpdateOrder 更新订单（未配货前可用）
// https://open.tongtool.com/apiDoc.html#/?docId=3e0d01bfe01441aa8e2071c2c88cc9fb
// 注意：当前接口仅支持手工订单和 API 导入的订单进行更新，平台订单不支持更新！
func (s service) UpdateOrder(req UpdateOrderRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId

	res := struct {
		tongtool.Response
		Datas string `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/orderUpdate")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		err = tongtool.ErrorWrap(res.Code, res.Message)
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}

// 作废订单处理

type CancelOrderRequest struct {
	MerchantId  string   `json:"merchantId"`  // 商戶 ID
	OrderIdKeys []string `json:"orderIdKeys"` // 通途订单 ID Key
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
		tongtool.Response
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
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			results = make([]OrderCancelResult, len(res.Datas.Array))
			for i := range res.Datas.Array {
				results[i] = OrderCancelResult{
					OrderId: res.Datas.Array[i].OrderId,
					Result:  strings.TrimSpace(res.Datas.Array[i].Result),
				}
			}
		}
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}

// 订单配对
// 订单未配货前可用；目前支持亚马逊、速卖通和 shopify 平台订单添加货品

type OrderPairTransaction struct {
	GoodsDetailId  string `json:"goodsDetailId"`  // 货品 ID
	OrderDetailsId string `json:"orderDetailsId"` // 订单详情 ID
	Quantity       int    `json:"quantity"`       // 数量
}

type OrderPairRequest struct {
	MerchantId   string                 `json:"merchantId"`   // 商户 ID
	Transactions []OrderPairTransaction `json:"transactions"` // 订单交易信息
	OrderId      string                 `json:"orderId"`      // 通途订单 ID
}

func (m OrderPairRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Transactions, validation.Required.Error("订单交易信息不能为空")),
		validation.Field(&m.Transactions, validation.When(len(m.Transactions) > 0, validation.Each(validation.WithContext(func(ctx context.Context, value interface{}) error {
			if transaction, ok := value.(OrderPairTransaction); !ok {
				return errors.New("无效的配对信息")
			} else {
				return validation.ValidateStruct(&transaction,
					validation.Field(&transaction.GoodsDetailId, validation.Required.Error("货品 ID 不能为空")),
					validation.Field(&transaction.OrderDetailsId, validation.Required.Error("订单详情 ID 不能为空")),
					validation.Field(&transaction.Quantity, validation.Required.Error("数量不能为空"), validation.Min(1).Error("数量不能小于 {{.threshold}}")),
				)
			}
		})))),
		validation.Field(&m.OrderId, validation.Required.Error("订单 ID 不能为空")),
	)
}

func (s service) OrderPair(req OrderPairRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/orderAddProduct")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		err = tongtool.ErrorWrap(res.Code, res.Message)
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}
