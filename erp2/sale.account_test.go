package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"testing"
)

func TestService_SaleAccounts(t *testing.T) {
	params := SaleAccountQueryParams{}
	for {
		accounts, isLastPage, err := ttService.SaleAccounts(params)
		if err == nil {
			fmt.Println(jsonx.ToJson(accounts, "[]"))
		} else {
			t.Errorf("ttService.SaleAccounts error: %s", err.Error())
		}
		if err != nil || isLastPage {
			break
		}
		params.PageNo++
	}
}
