package erp2

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	jsoniter "github.com/json-iterator/go"
)

// 亚马逊账号对应的站点

type AmazonAccountSitesQueryParams struct {
	Paging
	Account    string `json:"account"`    // 账号
	MerchantId string `json:"merchantId"` // 商户 ID
}

func (m AmazonAccountSitesQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Account, validation.Required.Error("帐号不能为空")),
	)
}

// AmazonAccountSites 查询亚马逊账号对应的站点
// https://open.tongtool.com/apiDoc.html#/?docId=4dd54cb61d6c4719860bec1d875f48af
func (s service) AmazonAccountSites(params AmazonAccountSitesQueryParams) (items []string, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
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
			Array []struct {
				SiteId string `json:"siteId"`
			} `json:"array"`
			PageNo   int `json:"pageNo"`
			PageSize int `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/queryAmazonAccountSiteId")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = make([]string, len(res.Datas.Array))
			for i := range res.Datas.Array {
				items[i] = res.Datas.Array[i].SiteId
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
