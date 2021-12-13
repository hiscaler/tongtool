package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
	"strings"
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

func (s service) Packages(params PackageQueryParams) (items []Package, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = s.tongTool.QueryDefaultValues.PageNo
	}
	if params.PageSize <= 0 {
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

func (s service) Package(orderNumber, packageNumber string) (item Package, err error) {
	params := PackageQueryParams{
		MerchantId: s.tongTool.MerchantId,
		OrderId:    strings.TrimSpace(orderNumber),
		PageNo:     1,
		PageSize:   s.tongTool.QueryDefaultValues.PageSize,
	}
	if packageNumber != "" {
		packageNumber = strings.TrimSpace(packageNumber)
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
				if packageNumber == "" {
					exists = true
					item = packages[len(packages)-1]
				} else {
					for _, p := range packages {
						if strings.EqualFold(p.PackageId, packageNumber) {
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
