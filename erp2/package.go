package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/pkg/cache"
	"github.com/hiscaler/tongtool/pkg/in"
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
	carrierCurrency       string        `json:"carrierCurrency"`       // 物流商运费币种
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
	TongtoolCurrency      string        `json:"tongtoolCurrency"`      // 通途运费币种
	TongtoolPostage       float64       `json:"tongtoolPostage"`       // 通途运费
	TongtoolWeight        float64       `json:"tongtoolWeight"`        // 通途包裹重量,单位g
	TrackingNumber        string        `json:"trackingNumber"`        // 跟踪号
	UpdatedDate           string        `json:"updatedDate"`           // 包裹更新时间
	UploadCarrier         string        `json:"uploadCarrier"`         // 上传包裹的Carrier
	WarehouseName         string        `json:"warehouseName"`         // 仓库名称
}

type PackageQueryParams struct {
	AssignTimeFrom     string `json:"assignTimeFrom,omitempty"`
	AssignTimeTo       string `json:"assignTimeTo,omitempty"`
	DespatchTimeFrom   string `json:"despatchTimeFrom,omitempty"`
	DespatchTimeTo     string `json:"despatchTimeTo,omitempty"`
	MerchantId         string `json:"merchantId"`
	OrderId            string `json:"orderId,omitempty"`
	PackageStatus      string `json:"packageStatus,omitempty"`
	PageNo             int    `json:"pageNo"`
	PageSize           int    `json:"pageSize"`
	ShippingMethodName string `json:"shippingMethodName,omitempty"`
}

// Packages 包裹列表
// https://open.tongtool.com/apiDoc.html#/?docId=0412c0185dce4a9d88714a9eef44932b
func (s service) Packages(params PackageQueryParams) (items []Package, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	var cacheKey string
	if s.tongTool.ActivateCache {
		cacheKey = cache.GenerateKey(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = json.Unmarshal(b, &items); e == nil {
				return
			}
		}
	}
	items = make([]Package, 0)
	res := struct {
		result
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
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				for i, item := range items {
					items[i].IsCheckedBoolean = in.StringIn(item.IsChecked, "Y")
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
	if err == nil && s.tongTool.ActivateCache {
		if b, e := json.Marshal(&items); e == nil {
			s.tongTool.Cache.Set(cacheKey, b)
		}
	}
	return
}

func (s service) Package(orderId, packageId string) (item Package, err error) {
	params := PackageQueryParams{
		MerchantId: s.tongTool.MerchantId,
		OrderId:    strings.TrimSpace(orderId),
		PageNo:     1,
		PageSize:   s.tongTool.QueryDefaultValues.PageSize,
	}
	if packageId != "" {
		packageId = strings.TrimSpace(packageId)
	}

	var cacheKey string
	if s.tongTool.ActivateCache {
		cacheKey = cache.GenerateKey(params, packageId)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = json.Unmarshal(b, &item); e == nil {
				return
			}
		}
	}

	exists := false
	for {
		packages := make([]Package, 0)
		isLastPage := false
		packages, isLastPage, err = s.Packages(params)
		if err == nil {
			if len(packages) == 0 {
				err = tongtool.ErrNotFound
			} else {
				for _, p := range packages {
					if p.PackageStatus != PackageStatusCancel {
						if packageId != "" {
							if strings.EqualFold(p.PackageId, packageId) {
								exists = true
								item = p
							}
						} else {
							exists = true
							item = p
						}
						if exists {
							break
						}
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
	if err == nil && s.tongTool.ActivateCache {
		if b, e := json.Marshal(&item); e == nil {
			s.tongTool.Cache.Set(cacheKey, b)
		}
	}

	return
}
