package erp2

import (
	"github.com/hiscaler/gox/jsonx"
	"testing"
)

func TestService_SurfaceSheets(t *testing.T) {
	params := SurfaceSheetsQueryParams{
		TrackingNumberList: []string{""},
	}
	items, err := ttService.SurfaceSheets(params)
	if err != nil {
		t.Errorf("ttService.SurfaceSheets error: %s", err.Error())
	} else {
		t.Log(jsonx.ToJson(items, "[]"))
	}
}
