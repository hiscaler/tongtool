package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_SaleAccounts(t *testing.T) {
	params := SaleAccountQueryParams{}
	for {
		accounts, isLastPage, err := ttService.SaleAccounts(params)
		if err == nil {
			fmt.Println(cast.ToJson(accounts))
		} else {
			t.Errorf("ttService.SaleAccounts error: %s", err.Error())
		}
		if err != nil || isLastPage {
			break
		}
		params.PageNo++
	}
}
