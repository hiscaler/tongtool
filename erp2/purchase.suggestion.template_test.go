package erp2

import (
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

// 采购建议模板查询
func TestService_PurchaseSuggestionTemplates(t *testing.T) {
	_, ttService := newTestTongTool()
	params := PurchaseSuggestionTemplateQueryParams{
		Names: []string{"test"},
	}
	templates := make([]PurchaseSuggestionTemplate, 0)
	for {
		pageTemplates, isLastPage, err := ttService.PurchaseSuggestionTemplates(params)
		if err != nil {
			t.Errorf("ttService.PurchaseSuggestionTemplates error: %s", err.Error())
		} else {
			templates = append(templates, pageTemplates...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	t.Log(cast.ToJson(templates))
}
