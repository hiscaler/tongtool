package erp2

import (
	"github.com/hiscaler/gox/jsonx"
	"testing"
)

// 采购建议查询
func TestService_PurchaseSuggestions(t *testing.T) {
	number := "0007000007201603230000076503"
	params := PurchaseSuggestionsQueryParams{
		PurchaseTemplateId: number,
		SKUs:               []string{"abc"},
	}
	suggestions := make([]PurchaseSuggestion, 0)
	for {
		pageSuggestions, isLastPage, err := ttService.PurchaseSuggestions(params)
		if err != nil {
			t.Errorf("ttService.PurchaseSuggestions error: %s", err.Error())
		} else {
			suggestions = append(suggestions, pageSuggestions...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	t.Log(jsonx.ToJson(suggestions, "[]"))
}
