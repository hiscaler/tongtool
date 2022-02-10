package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/gox/inx"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

// 模板类型
const (
	PurchaseSuggestionTemplateFBAType   = "fba"
	PurchaseSuggestionTemplateOtherType = "other"
)

// PurchaseSuggestionTemplate 采购建议模板
type PurchaseSuggestionTemplate struct {
	FormulaDaily         string `json:"formulaDaily"`         // 日均销量计算公式
	FrequencyOfPurchase  int    `json:"frequencyOfPurchase"`  // 采购频率(天)
	PurchaseTemplateId   string `json:"purchaseTemplateId"`   // 采购建议模板ID
	PurchaseTemplateName string `json:"purchaseTemplateName"` // 采购建议模板名称
	SuggestionType       string `json:"suggestionType"`       // 模版类型 FBA:FBA模版 other:其他模版
	WarehouseIdKeys      string `json:"warehouseIdKeys"`      // 业务ID 如果业务类型为warehouse则 ALW:所有本地仓库, AOW:所有海外仓库, ATPW:所有第三方仓库 , AFBA:所有FBA仓
	WarehouseName        string `json:"warehouseName"`        // 仓库名称
	WarehouseType        string `json:"warehouseType"`        // 仓库类型 owner/本地仓库、thirdParty/第三方仓库
}

type PurchaseSuggestionTemplatesQueryParams struct {
	Paging
	MerchantId string   `json:"merchantId"`      // 商户ID
	Names      []string `json:"names,omitempty"` // 采购建议模板名称（扩展）
}

// PurchaseSuggestionTemplates 采购建议模板列表
// https://open.tongtool.com/apiDoc.html#/?docId=129858303d494c6b90b552eeb5a7514f
func (s service) PurchaseSuggestionTemplates(params PurchaseSuggestionTemplatesQueryParams) (items []PurchaseSuggestionTemplate, isLastPage bool, err error) {
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
	items = make([]PurchaseSuggestionTemplate, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []PurchaseSuggestionTemplate `json:"array"`
			PageNo   int                          `json:"pageNo"`
			PageSize int                          `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/proposalTemplateQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				if len(params.Names) == 0 {
					items = res.Datas.Array
				} else {
					for _, d := range res.Datas.Array {
						if inx.StringIn(d.PurchaseTemplateName, params.Names...) {
							items = append(items, d)
						}
					}
				}
				isLastPage = len(res.Datas.Array) < params.PageSize
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
