package erp2

import (
	"errors"
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

type SaleAccountQueryParams struct {
	MerchantId string `json:"merchantId"`         // 商户ID
	PageNo     int    `json:"pageNo,omitempty"`   // 查询页数
	PageSize   int    `json:"pageSize,omitempty"` // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
}

type accountsResult struct {
	result
	Datas struct {
		Array    []SaleAccount `json:"array"`
		PageNo   int           `json:"pageNo"`
		PageSize int           `json:"pageSize"`
	}
}

// SaleAccounts 商户账号列表
// https://open.tongtool.com/apiDoc.html#/?docId=1e81e4bbae0b4d60b5f7777fc629ba2a
func (s service) SaleAccounts(params SaleAccountQueryParams) (items []SaleAccount, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]SaleAccount, 0)
	res := accountsResult{}
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
			err = errors.New(resp.Status())
		}
	}
	return
}
