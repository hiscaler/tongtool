package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

type SurfaceSheet struct {
	LabelURL       string `json:"labelUrl"`       // 面单地址
	TrackingNumber string `json:"trackingNumber"` // 跟踪号
}

type SurfaceSheetsQueryParams struct {
	MerchantId         string   `json:"merchantId"`         // 商户号
	TrackingNumberList []string `json:"trackingNumberList"` // 跟踪号列表,最多支持 100 个
}

func (m SurfaceSheetsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TrackingNumberList,
			validation.Required.Error("跟踪号列表不能为空"),
			validation.By(func(value interface{}) error {
				numbers, _ := value.([]string)
				if len(numbers) > 100 {
					return errors.New("跟踪号列表最多 100 个")
				}
				return nil
			}),
		),
	)
}

// SurfaceSheets 获取通途 ERP 面单
// https://open.tongtool.com/apiDoc.html#/?docId=f5c43f48665c44179399deb5c153765a
func (s service) SurfaceSheets(params SurfaceSheetsQueryParams) (items []SurfaceSheet, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
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
	items = make([]SurfaceSheet, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array []SurfaceSheet `json:"array"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/getErpLabel")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
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
