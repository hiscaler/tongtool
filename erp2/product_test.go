package erp2

import (
	"testing"
)

func TestService_Products(t *testing.T) {
	_, ttService := newTestTongTool()
	params := ProductQueryParams{}
	_, _, err := ttService.Products(params)
	if err != nil {
		t.Errorf("ttService.Products error: %s", err.Error())
	}
}
