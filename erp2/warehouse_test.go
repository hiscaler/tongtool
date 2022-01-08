package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_Warehouses(t *testing.T) {
	params := WarehouseQueryParams{}
	warehouses := make([]Warehouse, 0)
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

	enabledCount := 0
	for _, warehouse := range warehouses {
		if warehouse.StatusBoolean {
			enabledCount++
		}
	}
	fmt.Println(fmt.Sprintf("Total found %d warehouses, enabled warehouses: %d", len(warehouses), enabledCount))
}

func TestService_Warehouse(t *testing.T) {
	id := "a"
	warehouse, err := ttService.Warehouse(id)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		fmt.Println(cast.ToJson(warehouse))
	}
}

func TestService_ShippingMethods(t *testing.T) {
	params := ShippingMethodQueryParams{WarehouseId: "8151050530202008250000047045"}
	items, _, err := ttService.ShippingMethods(params)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		fmt.Println(cast.ToJson(items))
	}
}
