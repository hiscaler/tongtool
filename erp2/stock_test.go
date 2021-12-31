package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_StockChangeLogs(t *testing.T) {
	_, ttService := newTestTongTool()
	params := StockChangeLogQueryParams{
		UpdatedDateFrom: "2018-01-01 00:00:00",
		WarehouseName:   "万邑通美国西岸仓",
	}
	logs := make([]StockChangeLog, 0)
	for {
		pageLogs, isLastPage, err := ttService.StockChangeLogs(params)
		if err != nil {
			t.Errorf("ttService.StockChangeLogs error: %s", err.Error())
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
