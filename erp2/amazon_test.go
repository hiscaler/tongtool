package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"testing"
)

func TestService_AmazonAccountSites(t *testing.T) {
	params := AmazonAccountSitesQueryParams{
		Account: "a",
	}
	params.PageNo = 1
	logs := make([]string, 0)
	for {
		pageLogs, isLastPage, err := ttService.AmazonAccountSites(params)
		if err != nil {
			t.Errorf("ttService.AmazonAccountSites error: %s", err.Error())
		} else {
			fmt.Println(jsonx.ToJson(pageLogs, "[]"))
			logs = append(logs, pageLogs...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	fmt.Println(fmt.Sprintf("Total found %d logs", len(logs)))
}
