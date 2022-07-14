package erp2

import (
	"github.com/hiscaler/gox/jsonx"
	"strings"
	"testing"
)

// 供应商列表
func TestService_Suppliers(t *testing.T) {
	params := SuppliersQueryParams{}
	params.PageNo = 1
	name := "栀子花开女装店"
	found := false
	var supplier Supplier
	for {
		suppliers, isLastPage, err := ttService.Suppliers(params)
		if err != nil {
			t.Errorf("suppliers error: %s", err.Error())
		}

		for _, v := range suppliers {
			if strings.EqualFold(name, v.CorporationFullName) {
				found = true
				supplier = v
				break
			}
		}

		if isLastPage || found {
			break
		}
		params.PageNo++
	}
	if !found {
		t.Errorf("%s not found", name)
	}
	t.Log(jsonx.ToJson(supplier, "{}"))
}
