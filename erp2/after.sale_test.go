package erp2

import (
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_AfterSales(t *testing.T) {
	_, ttService := newTestTongTool()
	params := AfterSaleQueryParams{
		CreatedDateFrom: "2020-01-01 00:00:00",
		CreatedDateTo:   "2022-01-01 23:59:59",
	}
	items := make([]AfterSale, 0)
	pageItems, _, err := ttService.AfterSales(params)
	if err != nil {
		t.Errorf("ttService.AfterSales error: %s", err.Error())
	} else {
		items = append(items, pageItems...)
	}
	t.Log(cast.ToJson(items))
}
