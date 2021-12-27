package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
)

// Platform 平台
type Platform struct {
	PlatformId            string         `json:"platformId"`            // 平台id
	PlatformName          string         `json:"platformName"`          // 平台名称
	PlatformSites         []PlatformSite `json:"platformSites"`         // 平台对应站点
	PlatformStatus        string         `json:"platformStatus"`        // 平台状态 :0、有效 1、失效
	PlatformStatusBoolean bool           `json:"platformStatusBoolean"` // 平台状态
}

// PlatformSite 平台站点
type PlatformSite struct {
	CountryCode string `json:"countryCode"` // 站点对应国家简码
	SiteId      int    `json:"siteId"`      // 站点id
	SiteName    string `json:"siteName"`    // 站点名称
	TimeZone    string `json:"timeZone"`    // 站点时区
}

// Platforms 平台及站点信息
// https://open.tongtool.com/apiDoc.html#/?docId=3c5d0c2f549e4ebfb21d01c9e4cf5449
func (s service) Platforms() (items []Platform, err error) {
	items = make([]Platform, 0)
	res := struct {
		result
		Datas struct {
			Array []Platform `json:"array"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(map[string]string{"merchantId": s.tongTool.MerchantId}).
		SetResult(&res).
		Post("/openapi/tongtool/merchantPlatformQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				for i, item := range items {
					items[i].PlatformStatusBoolean = item.PlatformStatus == "0"
				}
			}
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}
	return
}
