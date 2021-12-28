package erp2

import (
	"fmt"
	"strings"
	"testing"
)

// 采购建议模板查询
func TestService_PurchaseSuggestionTemplates(t *testing.T) {
	_, ttService := newTestTongTool()
	templateId := "6014000007201703150000118544"
	params := PurchaseSuggestionTemplateQueryParams{}
	templates, _, err := ttService.PurchaseSuggestionTemplates(params)
	if err == nil {
		exists := false
		for _, template := range templates {
			if strings.EqualFold(templateId, template.PurchaseTemplateId) {
				exists = true
				break
			}
		}
		if exists {
			fmt.Println(fmt.Sprintf("Purchase suggestion templates: %#v", templates))
		} else {
			t.Errorf("not found %s", templateId)
		}
	} else {
		t.Error(err)
	}
}
