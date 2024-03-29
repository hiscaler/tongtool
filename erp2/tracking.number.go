package erp2

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type TrackingNumber struct {
	CarrierCode        string `json:"carrierCode"`        // 物流代码
	CarrierName        string `json:"carrierName"`        // 物流名称
	OrderId            string `json:"orderId"`            // 订单号
	ShippingMethodCode string `json:"shippingMethodCode"` // 邮寄方式代码
	ShippingMethodName string `json:"shippingMethodName"` // 邮寄方式名称
	TrackingNumber     string `json:"trackingNumber"`     // 跟踪号
	// 扩展属性
	IsMatched bool `json:"isMatched"` // 是否匹配
}

type TrackingNumbersQueryParams struct {
	Paging
	MerchantId string   `json:"merchantId"` // 商户 ID
	OrderIds   []string `json:"orderIds"`   // orderNumber 集合
}

func (m TrackingNumbersQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderIds, validation.Required.Error("订单号不能为空")),
	)
}

// TrackingNumbers 订单物流单号列表
// 需要注意的是该封装总是返回包含所有查询订单集合的数据，无论是否有物流数据
// https://open.tongtool.com/apiDoc.html#/?docId=3b3cceec8fe04e6db44da17ec4b38f08
func (s service) TrackingNumbers(params TrackingNumbersQueryParams) (items []TrackingNumber, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if err = params.Validate(); err != nil {
		return
	}

	items = make([]TrackingNumber, len(params.OrderIds))
	for i, orderId := range params.OrderIds {
		items[i] = TrackingNumber{OrderId: orderId}
	}
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
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
			Array    []TrackingNumber `json:"array"`
			PageNo   int              `json:"pageNo"`
			PageSize int              `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/trackingNumberQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			for _, d := range res.Datas.Array {
				for i, item := range items {
					if strings.EqualFold(d.OrderId, item.OrderId) {
						d.IsMatched = true
						items[i] = d
					}
				}
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
