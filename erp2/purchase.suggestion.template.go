package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
)

// 模板类型
const (
	PurchaseOrderTemplateFBAType   = "fba"
	PurchaseOrderTemplateOtherType = "other"
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

type PurchaseSuggestionTemplateQueryParams struct {
	MerchantId string `json:"merchantId"`                   // 商户ID
	PageNo     int    `json:"pageNo,omitempty,omitempty"`   // 查询页数
	PageSize   int    `json:"pageSize,omitempty,omitempty"` // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
}

// PurchaseSuggestionTemplates 采购建议模板列表
// https://open.tongtool.com/apiDoc.html#/?docId=129858303d494c6b90b552eeb5a7514f
func (s service) PurchaseSuggestionTemplates(params PurchaseSuggestionTemplateQueryParams) (items []PurchaseSuggestionTemplate, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	items = make([]PurchaseSuggestionTemplate, 0)
	res := struct {
		result
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
