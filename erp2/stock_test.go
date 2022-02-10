package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"testing"
)

func TestService_Stocks(t *testing.T) {
	params := StocksQueryParams{}
	stocks := make([]Stock, 0)
	for {
		pageItems, isLastPage, err := ttService.Stocks(params)
		if err != nil {
			t.Errorf("ttService.Stocks error: %s", err.Error())
		} else {
			fmt.Println(jsonx.ToJson(pageItems, "[]"))
			stocks = append(stocks, pageItems...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	fmt.Println(fmt.Sprintf("Total found %d stocks", len(stocks)))
}

func TestService_StockChangeLogs(t *testing.T) {
	params := StockChangeLogsQueryParams{
		UpdatedDateFrom: "2018-01-01 00:00:00",
		WarehouseName:   "万邑通美国西岸仓",
	}
	logs := make([]StockChangeLog, 0)
	for {
		pageLogs, isLastPage, err := ttService.StockChangeLogs(params)
		if err != nil {
			t.Errorf("ttService.StockChangeLogs error: %s", err.Error())
		} else {
			logs = append(logs, pageLogs...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	fmt.Println(fmt.Sprintf("Total found %d logs", len(logs)))
}
