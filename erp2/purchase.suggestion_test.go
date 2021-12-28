package erp2

import (
	"fmt"
	"testing"
)

// 采购建议查询
func TestService_PurchaseSuggestions(t *testing.T) {
	_, ttService := newTestTongTool()
	number := "6014000007201703150000118544"
	params := PurchaseSuggestionQueryParams{
		PurchaseTemplateId: number,
	}
	suggestions, _, err := ttService.PurchaseSuggestions(params)
	if err == nil {
		exists := false
		//for _, order := range suggestions {
		//	if strings.EqualFold(number, order.PoNum) {
		//		exists = true
		//		break
		//	}
		//}
		if exists {
			fmt.Println(fmt.Sprintf("Orders: %#v", suggestions))
		} else {
			t.Errorf("not found %s", number)
		}
	} else {
		t.Error(err)
	}
}
