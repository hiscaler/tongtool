package erp2

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/inx"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/json-iterator/go"
)

// PurchaseSuggestion 采购建议
type PurchaseSuggestion struct {
	CalculateDate          string  `json:"caculateDate"`           // 采购建议计算时间
	CurrStockQuantity      int     `json:"currStockQuantity"`      // 可用库存数
	DailySales             float64 `json:"dailySales"`             // 日均销量
	DeliveryDays           int     `json:"devliveryDays"`          // 安全交期
	GoodsIdKey             string  `json:"goodsIdKey"`             // 商品 ID Key
	GoodsSKU               string  `json:"goodsSku"`               // 商品 SKU
	InTransitStockQuantity int     `json:"intransitStockQuantity"` // 在途库存数
	ProposalQuantity       int     `json:"proposalQuantity"`       // 采购建议数量
	SaleAvg15              float64 `json:"saleAvg15"`              // 15 天销量
	SaleAvg30              float64 `json:"saleAvg30"`              // 30 天销量
	SaleAvg7               float64 `json:"saleAvg7"`               // 7 天销量
	UnpickingQuantity      int     `json:"unpickingQuantity"`      // 订单未配货数量
	WarehouseIdKey         string  `json:"warehouseIdKey"`         // 仓库 ID Key
	WarehouseName          string  `json:"warehouseName"`          // 仓库名称
}

type PurchaseSuggestionsQueryParams struct {
	Paging
	MerchantId         string   `json:"merchantId"`               // 商户 ID
	PurchaseTemplateId string   `json:"purchaseTemplateId"`       // 采购建议模板 ID
	WarehouseNames     []string `json:"warehouseNames,omitempty"` // 仓库（扩展）
	SKUs               []string `json:"skus,omitempty"`           // SKU 列表（扩展）
}

func (m PurchaseSuggestionsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PurchaseTemplateId, validation.Required.Error("采购建议模板 ID 不能为空")),
	)
}

// PurchaseSuggestions 采购建议列表
// https://open.tongtool.com/apiDoc.html#/?docId=8e80fde6a4824b288d17bc04be8f4ef6
func (s service) PurchaseSuggestions(params PurchaseSuggestionsQueryParams) (items []PurchaseSuggestion, isLastPage bool, err error) {
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
			Array    []PurchaseSuggestion `json:"array"`
			PageNo   int                  `json:"pageNo"`
			PageSize int                  `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/proposalResultQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			if len(params.SKUs) == 0 && len(params.WarehouseNames) == 0 {
				items = res.Datas.Array
			} else {
				for _, d := range res.Datas.Array {
					if len(params.SKUs) != 0 && !inx.StringIn(d.GoodsSKU, params.SKUs...) ||
						len(params.WarehouseNames) != 0 && !inx.StringIn(d.WarehouseName, params.WarehouseNames...) {
						continue
					}
					items = append(items, d)
				}
			}
			isLastPage = len(res.Datas.Array) < params.PageSize
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
