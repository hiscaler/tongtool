package erp2

import (
	"strings"
	"testing"
)

// 供应商列表
func TestService_Suppliers(t *testing.T) {
	params := SuppliersQueryParams{}
	params.PageNo = 1
	name := "栀子花开女装店"
	found := false
	for {
		suppliers, isLastPage, err := ttService.Suppliers(params)
		if err == nil {
			for _, supplier := range suppliers {
				if strings.EqualFold(name, supplier.CorporationFullName) {
					found = true
					break
				}
			}
		} else {
			t.Error(err)
			break
		}
		if isLastPage || found {
			break
		}
		params.PageNo++
	}
	if !found {
		t.Errorf("%s not found", name)
	}
}
