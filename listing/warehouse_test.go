package listing

import (
	"testing"
)

func TestService_Warehouses(t *testing.T) {
	params := WarehousesQueryParams{}
	_, err := ttService.Warehouses(params)
	if err != nil {
		t.Errorf("ttService.Warehouses error: %s", err.Error())
	}
}
