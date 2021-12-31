package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_AmazonAccountSites(t *testing.T) {
	_, ttService := newTestTongTool()
	params := AmazonAccountSiteQueryParams{}
	logs := make([]string, 0)
	for {
		pageLogs, isLastPage, err := ttService.AmazonAccountSites(params)
		if err != nil {
			t.Errorf("ttService.AmazonAccountSites error: %s", err.Error())
		} else {
			fmt.Println(cast.ToJson(pageLogs))
			logs = append(logs, pageLogs...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	fmt.Println(fmt.Sprintf("Total found %d logs", len(logs)))
}
