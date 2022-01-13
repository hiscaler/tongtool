package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/pkg/cache"
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

type TrackingNumberQueryParams struct {
	MerchantId string   `json:"merchantId"`         // 商户ID
	OrderIds   []string `json:"orderIds,omitempty"` // orderNumber集合
	PageNo     int      `json:"pageNo,omitempty"`   // 查询页数
	PageSize   int      `json:"pageSize,omitempty"` // 每页数量
}

// TrackingNumbers 订单物流单号列表
// 需要注意的是该封装总是返回包含所有查询订单集合的数据，无论是否有物流数据
// https://open.tongtool.com/apiDoc.html#/?docId=3b3cceec8fe04e6db44da17ec4b38f08
func (s service) TrackingNumbers(params TrackingNumberQueryParams) (items []TrackingNumber, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	items = make([]TrackingNumber, 0)
	if len(params.OrderIds) == 0 {
		return
	}
	for _, orderId := range params.OrderIds {
		items = append(items, TrackingNumber{OrderId: orderId})
	}
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = cache.GenerateKey(params)
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
	res := struct {
		result
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
	if err == nil {
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
