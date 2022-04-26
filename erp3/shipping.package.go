package erp3

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
	jsoniter "github.com/json-iterator/go"
)

// 出库单交运
// https://open.tongtool.com/apiDoc.html#/?docId=d137b0c765104292afb2cc042cb2f961

type ShippingPackage struct {
	ErrorCode             int    `json:"errorCode"`             // 错误代码
	ErrorMsg              string `json:"errorMsg"`              // 错误信息
	LabelPath             string `json:"labelPath"`             // 物流商地址标签路径
	PackageCode           string `json:"packageCode"`           // 出库单号
	ThirdPartyNo          string `json:"thirdPartyNo"`          // 物流商单号
	TrackingNumber        string `json:"trackingNumber"`        // 承运商运单号
	VirtualTrackingNumber string `json:"virtualTrackingNumber"` // 承运商虚拟运单号（部分物流商返回）
}

type AddShippingPackageRequest struct {
	MerchantId      string   `json:"merchantId"`      // 商户号
	PackageCodeList []string `json:"packageCodeList"` // 出库单编号列表
}

func (m AddShippingPackageRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PackageCodeList, validation.Required.Error("出库单编号列表不能为空")),
	)
}

func (s service) AddShippingPackage(req AddShippingPackageRequest) (packages []ShippingPackage, err error) {
	req.MerchantId = s.tongTool.MerchantId
	if err = req.Validate(); err != nil {
		return
	}

	res := struct {
		tongtool.Response
		Datas []ShippingPackage `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/packageInfo/addShippingPackage")
	if err != nil {
		return
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
	return
}
