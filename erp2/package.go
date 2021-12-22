package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
	"strings"
)

// 包裹状态
const (
	PackageStatusWaitPrint   = "waitPrint"   // 等待打印
	PackageStatusWaitDeliver = "waitDeliver" // 等待发货
	PackageStatusDelivered   = "delivered"   // 已发货
	PackageStatusCancel      = "cancel"      // 作废
)

type Package struct {
	PackageId             string  `json:"packageId"`             // 包裹id
	TrackingNumber        string  `json:"trackingNumber"`        // 跟踪号
	PackageStatus         string  `json:"packageStatus"`         // 包裹状态
	ThirdPartyPackageCode string  `json:"thirdPartyPackageCode"` // 物流商单号
	TongtoolWeight        float64 `json:"tongtoolWeight"`        // 通途包裹重量,单位g
	TongtoolPostage       float64 `json:"tongtoolPostage"`       // 通途运费
	CarrierWeight         float64 `json:"carrierWeight"`         // 物流商称重重量,单位g
	CarrierPostage        float64 `json:"carrierPostage"`        // 物流商运费
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

type packageResult struct {
	result
	Datas struct {
		Array    []Package `json:"array"`
		PageNo   int       `json:"pageNo"`
		PageSize int       `json:"pageSize"`
	}
}

// Packages 包裹列表
// https://open.tongtool.com/apiDoc.html#/?docId=0412c0185dce4a9d88714a9eef44932b
func (s service) Packages(params PackageQueryParams) (items []Package, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]Package, 0)
	res := packageResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/packagesQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				isLastPage = len(items) < params.PageSize
			}
		} else {
			err = errors.New(resp.Status())
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

	exists := false
	for {
		packages := make([]Package, 0)
		isLastPage := false
		packages, isLastPage, err = s.Packages(params)
		if err == nil {
			if len(packages) == 0 {
				err = errors.New("not found")
			} else {
				if packageId == "" {
					exists = true
					item = packages[len(packages)-1] // last package
				} else {
					for _, p := range packages {
						if strings.EqualFold(p.PackageId, packageId) {
							exists = true
							item = p
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
		err = errors.New("not found")
	}

	return
}
