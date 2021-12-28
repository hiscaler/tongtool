package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
)

// PurchaseSuggestion 采购建议
type PurchaseSuggestion struct {
	CaculateDate           string `json:"caculateDate"`           // 采购建议计算时间
	CurrStockQuantity      int    `json:"currStockQuantity"`      // 可用库存数
	DailySales             int    `json:"dailySales"`             // 日均销量
	DevliveryDays          int    `json:"devliveryDays"`          // 安全交期
	GoodsIdKey             string `json:"goodsIdKey"`             // 商品id key
	GoodsSku               string `json:"goodsSku"`               // 商品sku
	IntransitStockQuantity int    `json:"intransitStockQuantity"` // 在途库存数
	ProposalQuantity       int    `json:"proposalQuantity"`       // 采购建议数量
	SaleAvg15              int    `json:"saleAvg15"`              // 15天销量
	SaleAvg30              int    `json:"saleAvg30"`              // 30天销量
	SaleAvg7               int    `json:"saleAvg7"`               // 7天销量
	UnpickingQuantity      int    `json:"unpickingQuantity"`      // 订单未配货数量
	WarehouseIdKey         int    `json:"warehouseIdKey"`         // 仓库id key
	WarehouseName          int    `json:"warehouseName"`          // 仓库名称
}

type PurchaseSuggestionQueryParams struct {
	MerchantId         string `json:"merchantId"`                   // 商户ID
	PageNo             int    `json:"pageNo,omitempty,omitempty"`   // 查询页数
	PageSize           int    `json:"pageSize,omitempty,omitempty"` // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
	PurchaseTemplateId string `json:"purchaseTemplateId"`           // 采购建议模板id
}

// PurchaseSuggestions 采购建议列表
// https://open.tongtool.com/apiDoc.html#/?docId=8e80fde6a4824b288d17bc04be8f4ef6
func (s service) PurchaseSuggestions(params PurchaseSuggestionQueryParams) (items []PurchaseSuggestion, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	items = make([]PurchaseSuggestion, 0)
	res := struct {
		result
		Datas struct {
			Array    []PurchaseSuggestion `json:"array"`
			PageNo   int                  `json:"pageNo"`
			PageSize int                  `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/proposalResultQuery")
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
