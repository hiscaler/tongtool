package erp2

import (
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/inx"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	"strings"
)

// 包裹状态
const (
	PackageStatusWaitPrint   = "waitPrint"   // 等待打印
	PackageStatusWaitDeliver = "waitDeliver" // 等待发货
	PackageStatusDelivered   = "delivered"   // 已发货
	PackageStatusCancel      = "cancel"      // 作废
	PackageStatusOther       = "other"       // 其他
)

type PackageItem struct {
	GoodsSKU string `json:"goodsSKU"` // 通途货品sku
	Quantity int    `json:"quantity"` // 采购数量
}

type Package struct {
	CarrierCurrency       string        `json:"carrierCurrency"`       // 物流商运费币种
	CarrierPostage        float64       `json:"carrierPostage"`        // 物流商运费
	CarrierWeight         float64       `json:"carrierWeight"`         // 物流商称重重量,单位g
	GoodsDetails          []PackageItem `json:"goodsDetails"`          // 包裹商品项目
	IsChecked             string        `json:"isChecked"`             // 包裹是否校验Y/已校验 、 null or N/未校验
	IsCheckedBoolean      bool          `json:"isCheckedBoolean"`      // 包裹是否校验布尔值
	MerchantId            string        `json:"merchantId"`            // 商户id
	PackageId             string        `json:"packageId"`             // 包裹id
	PackageStatus         string        `json:"packageStatus"`         // 包裹状态
	ShippingMethodCode    string        `json:"shippingMethodCode"`    // 邮寄方式代码
	ShippingMethodName    string        `json:"shippingMethodName"`    // 邮寄方式名称
	ThirdPartyPackageCode string        `json:"thirdPartyPackageCode"` // 物流商单号
	TongToolCurrency      string        `json:"tongtoolCurrency"`      // 通途运费币种
	TongToolPostage       float64       `json:"tongtoolPostage"`       // 通途运费
	TongToolWeight        float64       `json:"tongtoolWeight"`        // 通途包裹重量,单位g
	TrackingNumber        string        `json:"trackingNumber"`        // 跟踪号
	UpdatedDate           string        `json:"updatedDate"`           // 包裹更新时间
	UploadCarrier         string        `json:"uploadCarrier"`         // 上传包裹的Carrier
	WarehouseName         string        `json:"warehouseName"`         // 仓库名称
	// 自定义属性
	IsValid bool `json:"isValid"` // 是否有效
}

type PackagesQueryParams struct {
	Paging
	AssignTimeFrom     string `json:"assignTimeFrom,omitempty"`     // 配货开始时间
	AssignTimeTo       string `json:"assignTimeTo,omitempty"`       // 配货结束时间
	DespatchTimeFrom   string `json:"despatchTimeFrom,omitempty"`   // 发货开始时间
	DespatchTimeTo     string `json:"despatchTimeTo,omitempty"`     // 发货结束时间
	MerchantId         string `json:"merchantId"`                   // 商户ID
	OrderNumber        string `json:"orderId,omitempty"`            // 订单号
	PackageStatus      string `json:"packageStatus,omitempty"`      // 包裹状态： waitPrint 等待打印 waitDeliver 等待发货 delivered 已发货 cancel 作废
	ShippingMethodName string `json:"shippingMethodName,omitempty"` // 邮寄方式名称
}

func (m PackagesQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.AssignTimeFrom, validation.When(m.AssignTimeFrom != "", validation.Date(constant.DatetimeFormat).Error("配货开始时间日期格式无效"))),
		validation.Field(&m.AssignTimeTo, validation.When(m.AssignTimeTo != "", validation.Date(constant.DatetimeFormat).Error("配货结束时间日期格式无效"))),
		validation.Field(&m.DespatchTimeFrom, validation.When(m.DespatchTimeFrom != "", validation.Date(constant.DatetimeFormat).Error("发货开始时间日期格式无效"))),
		validation.Field(&m.DespatchTimeTo, validation.When(m.DespatchTimeTo != "", validation.Date(constant.DatetimeFormat).Error("发货结束时间日期格式无效"))),
		validation.Field(&m.PackageStatus, validation.In(PackageStatusWaitPrint, PackageStatusWaitDeliver, PackageStatusDelivered, PackageStatusCancel)),
	)
}

// Packages 包裹列表
// https://open.tongtool.com/apiDoc.html#/?docId=0412c0185dce4a9d88714a9eef44932b
func (s service) Packages(params PackagesQueryParams) (items []Package, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
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
	items = make([]Package, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []Package `json:"array"`
			PageNo   int       `json:"pageNo"`
			PageSize int       `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/packagesQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
			for i, item := range items {
				items[i].IsValid = !inx.StringIn(item.PackageStatus, PackageStatusCancel)
				items[i].IsCheckedBoolean = inx.StringIn(item.IsChecked, "Y")
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
	if err != nil {
		return
	}

	if s.tongTool.EnableCache && len(items) > 0 {
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

// Package 根据订单号和包裹号获取包裹资料
// 调用本函数时，必须同时提供订单号和包裹号（因为一个订单可能存在多个包裹号），仅返回一个有效的包裹信息（包裹状态不是 cancel 的）
// 在使用 item 数据前，您需要先判断 exists 返回值是否为 true 再使用，如果是 false 则需要根据 err 的错误信息确定您后续的业务处理。
//
// 如果需要查询一个订单所有的包裹，请使用 Packages 方法并提供 OrderNumber 参数值
func (s service) Package(orderNumber, packageNumber string) (item Package, exists bool, err error) {
	orderNumber = strings.TrimSpace(orderNumber)
	packageNumber = strings.TrimSpace(packageNumber)
	if orderNumber == "" || packageNumber == "" {
		err = errors.New("订单号和包裹号不能为空")
		return
	}
	params := PackagesQueryParams{
		MerchantId:  s.tongTool.MerchantId,
		OrderNumber: strings.TrimSpace(orderNumber),
	}
	params.PageNo = 1
	for {
		packages := make([]Package, 0)
		isLastPage := false
		packages, isLastPage, err = s.Packages(params)
		if err == nil {
			if len(packages) == 0 {
				err = tongtool.ErrNotFound
			} else {
				for _, p := range packages {
					if p.IsValid && strings.EqualFold(p.PackageId, packageNumber) {
						exists = true
						item = p
						break
					}
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

// 包裹发货处理

type PackageDeliverItemVolume struct {
	Height float64 `json:"height,omitempty"` // 高cm
	Length float64 `json:"length,omitempty"` // 长cm
	Width  float64 `json:"width,omitempty"`  // 宽cm
}

type PackageDeliverItem struct {
	RelatedNo      string                   `json:"relatedNo"`                // 识别号(包裹号、物流跟踪号、物流商处理号、虚拟跟踪号)
	ShipFee        float64                  `json:"shipFee,omitempty"`        // 运费￥
	TrackingNumber string                   `json:"trackingNumber,omitempty"` // 跟踪号
	Volume         PackageDeliverItemVolume `json:"volume,omitempty"`         // 体积cm³
	Weight         float64                  `json:"weight,omitempty"`         // 称重g
}

type PackageDeliverRequest struct {
	DeliverInfos  []PackageDeliverItem `json:"deliverInfos"`  // 发货信息列表
	MerchantId    string               `json:"merchantId"`    // 商户ID
	WarehouseName string               `json:"warehouseName"` // 仓库名称
}

func (m PackageDeliverRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.WarehouseName, validation.Required.Error("仓库名称不能为空")),
		validation.Field(&m.DeliverInfos, validation.Required.Error("发货信息不能为空"), validation.By(func(value interface{}) error {
			items, ok := value.([]PackageDeliverItem)
			if !ok {
				return errors.New("无效的发货信息")
			}
			for i, item := range items {
				if item.RelatedNo == "" {
					return fmt.Errorf("数据 %d 中识别号不能为空", i+1)
				}
			}
			return nil
		})),
	)
}

// PackageDeliver 执行包裹发货
// https://open.tongtool.com/apiDoc.html#/?docId=3493953e628b4f0ca5d32d3f6ac9d545
func (s service) PackageDeliver(req PackageDeliverRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
		Datas struct {
			ErrorList []struct {
				RelatedNo string `json:"relatedNo"`
				Message   string `json:"msg"`
			} `json:"errorList"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/packageDeliver")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			if len(res.Datas.ErrorList) != 0 {
				errorMessageNumbers := make(map[string][]string, 0)
				for _, item := range res.Datas.ErrorList {
					msg := strings.TrimSpace(item.Message)
					if numbers, ok := errorMessageNumbers[msg]; ok {
						numbers = append(numbers, item.RelatedNo)
						errorMessageNumbers[msg] = numbers
					} else {
						errorMessageNumbers[msg] = []string{item.RelatedNo}
					}
				}
				errorMessages := make([]string, 0)
				for msg, numbers := range errorMessageNumbers {
					errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", strings.Join(numbers, ","), msg))
				}
				err = errors.New(strings.Join(errorMessages, "; "))
			}
		}
	} else {
		if e := json.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}
