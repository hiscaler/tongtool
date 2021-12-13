package erp2

import (
	"fmt"
	"testing"
)

func TestService_Warehouses(t *testing.T) {
	_, ttService := newTestTongTool()
	params := WarehouseQueryParams{}
	warehouses := make([]Warehouse, 9)
	for {
		pageWarehouses, isLastPage, err := ttService.Warehouses(params)
		if err != nil {
			t.Errorf("ttService.Warehouses error: %s", err.Error())
		} else {
			warehouses = append(warehouses, pageWarehouses...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	fmt.Println(fmt.Sprintf("Total found %d warehouses", len(warehouses)))
}
