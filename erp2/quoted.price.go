package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
)

// 商品供应商报价

type QuotedPrice struct {
	Currency        string  `json:"currency"`        // 币种
	GoodsSKU        string  `json:"currency"`        // 商品sku
	Price           float64 `json:"price"`           // 供应商最新报价
	PurchaseLink    string  `json:"purchaseLink"`    // 采购链接
	QuotedPriceDate string  `json:"quotedPriceDate"` // 报价时间
	SupplierName    string  `json:"supplierName"`    // 供应商名称
}

type QuotedPriceQueryParams struct {
	MerchantId           string `json:"merchantId"` // 商家 ID
	PageNo               int    `json:"pageNo"`
	PageSize             int    `json:"pageSize"`
	QuotedPriceDateBegin string `json:"quotedPriceDateBegin"` // 报价起始时间
	QuotedPriceDateEnd   string `json:"quotedPriceDateEnd"`   // 报价结束时间
	SKU                  string `json:"sku"`                  // 商品 SKU
}

// QuotePrices 供应商报价查询
// https://open.tongtool.com/apiDoc.html#/?docId=0a508970886f4c7596b064f3b37987c9
func (s service) QuotePrices(params QuotedPriceQueryParams) (items []QuotedPrice, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	items = make([]QuotedPrice, 0)
	res := struct {
		result
		Datas struct {
			Array    []QuotedPrice `json:"array"`
			PageNo   int           `json:"pageNo"`
			PageSize int           `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/goodsPriceQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
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
	return
}
