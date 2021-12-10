package erp2

import (
	"errors"
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

type PackageQueryParam struct {
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

func (s service) Packages(params PackageQueryParam) (items []Package, isLastPage bool, err error) {
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
			items = res.Datas.Array
			isLastPage = len(items) < params.PageSize
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}

func (s service) Package(orderNumber, packageNumber string) (Package, error) {
	pkg := Package{}
	var err error
	params := PackageQueryParam{
		MerchantId: s.tongTool.MerchantId,
		OrderId:    strings.TrimSpace(orderNumber),
		PageNo:     1,
		PageSize:   s.tongTool.QueryDefaultValues.PageSize,
	}
	if packageNumber != "" {
		packageNumber = strings.TrimSpace(packageNumber)
	}
	for {
		packages := make([]Package, 0)
		isLastPage := false
		packages, isLastPage, err = s.Packages(params)
		if err == nil {
			if len(packages) == 0 {
				err = errors.New("not found")
			} else {
				if packageNumber == "" {
					pkg = packages[len(packages)-1]
				} else {
					for _, p := range packages {
						if strings.EqualFold(p.PackageId, packageNumber) {
							pkg = p
							break
						}
					}
				}
			}
		}
		if err != nil || isLastPage {
			break
		}
		params.PageNo++
	}

	return pkg, err
}
