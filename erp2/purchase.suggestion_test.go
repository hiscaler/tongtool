package erp2

import (
	"testing"
)

// 采购建议查询
func TestService_PurchaseSuggestions(t *testing.T) {
	_, ttService := newTestTongTool()
	number := "0007000007201603230000076503"
	params := PurchaseSuggestionQueryParams{
		PurchaseTemplateId: number,
	}
	suggestions, _, err := ttService.PurchaseSuggestions(params)
	if err == nil {
		if len(suggestions) == 0 {
			t.Errorf("not found suggestion with %s", number)
		}
	} else {
		t.Error(err)
	}
}
