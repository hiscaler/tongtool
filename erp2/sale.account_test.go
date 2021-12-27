package erp2

import (
	"fmt"
	"testing"
)

func TestService_SaleAccounts(t *testing.T) {
	_, ttService := newTestTongTool()
	params := SaleAccountQueryParams{}
	for {
		accounts, isLastPage, err := ttService.SaleAccounts(params)
		if err == nil {
			fmt.Println(fmt.Sprintf("%#v", accounts))
		} else {
			t.Errorf("ttService.SaleAccounts error: %s", err.Error())
		}
		if err != nil || isLastPage {
			break
		}
		params.PageNo++
	}
}
