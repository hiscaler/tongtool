package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"testing"
)

func TestService_SaleAccounts(t *testing.T) {
	params := SaleAccountsQueryParams{
		PlatformId: PlatformCouPang,
	}
	params.PageNo = 1
	var accounts []SaleAccount
	for {
		pageAccounts, isLastPage, err := ttService.SaleAccounts(params)
		if err == nil {
			accounts = append(accounts, pageAccounts...)
		} else {
			t.Errorf("ttService.SaleAccounts error: %s", err.Error())
		}
		if err != nil || isLastPage {
			break
		}
		params.PageNo++
	}
	fmt.Println(jsonx.ToJson(accounts, "[]"))
}
