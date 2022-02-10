package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

// SaleAccount 商户账号信息
type SaleAccount struct {
	SaleAccountId string   `json:"saleAccountId"` // 通途账户id
	Account       string   `json:"account"`       // 账户
	AccountCode   string   `json:"accountCode"`   // 账户简码
	PlatformId    string   `json:"platformId"`    // 平台id
	SiteIds       []string `json:"siteIds"`       // 站点id列表
	Status        string   `json:"status"`        // 账号状态 0停用,1 启用
	StatusBoolean bool     `json:"statusBoolean"` // 账号状态布尔值
}

type SaleAccountsQueryParams struct {
	Paging
	MerchantId string `json:"merchantId"` // 商户ID
}

// SaleAccounts 商户账号列表
// https://open.tongtool.com/apiDoc.html#/?docId=1e81e4bbae0b4d60b5f7777fc629ba2a
func (s service) SaleAccounts(params SaleAccountsQueryParams) (items []SaleAccount, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
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
	items = make([]SaleAccount, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []SaleAccount `json:"array"`
			PageNo   int           `json:"pageNo"`
			PageSize int           `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/merchantSaleAccountQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				for i, item := range items {
					items[i].StatusBoolean = item.Status == "1"
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
