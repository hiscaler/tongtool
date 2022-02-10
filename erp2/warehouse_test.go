package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"testing"
)

func TestService_Warehouses(t *testing.T) {
	params := WarehousesQueryParams{}
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
		fmt.Println(jsonx.ToJson(warehouse, "[]"))
	}
}

func TestService_ShippingMethods(t *testing.T) {
	params := ShippingMethodsQueryParams{WarehouseId: "8151050530202008250000047045"}
	items, _, err := ttService.ShippingMethods(params)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		fmt.Println(jsonx.ToJson(items, "[]"))
	}
}
