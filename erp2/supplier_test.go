package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"strings"
	"testing"
)

// 供应商列表
func TestService_Suppliers(t *testing.T) {
	_, ttService := newTestTongTool()
	params := SuppliersQueryParams{}
	findName := "栀子花开女装店"
	for {
		suppliers, isLastPage, err := ttService.Suppliers(params)
		if err == nil {
			for _, supplier := range suppliers {
				if strings.EqualFold(findName, supplier.CorporationFullName) {
					fmt.Println(cast.ToJson(supplier))
				}
			}
		} else {
			t.Error(err)
			break
		}
		if isLastPage {
			break
		}
		params.PageNo++
	}
}
