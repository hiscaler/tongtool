package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_Packages(t *testing.T) {
	_, ttService := newTestTongTool()
	params := PackageQueryParams{
		AssignTimeFrom: "2021-12-01 00:00:00",
		AssignTimeTo:   "2021-12-11 23:59:59",
	}
	packages, _, err := ttService.Packages(params)
	if err == nil {
		fmt.Println(cast.ToJson(packages))
	} else {
		t.Error(err)
	}
}

func TestService_Package(t *testing.T) {
	_, ttService := newTestTongTool()
	orderNumber := "L-M20211221152430918"
	packageNumber := "P02914669"
	_, err := ttService.Package(orderNumber, packageNumber)
	if err != nil {
		t.Error(err)
	}
}

func TestService_PackageWithCache(t *testing.T) {
	tt, ttService := newTestTongTool()
	tt.WithCache(true)
	fmt.Println(fmt.Sprintf("1111 %#v", tt))
	for i := 0; i <= 400; i++ {
		orderNumber := "L-M20211221152430918"
		packageNumber := "P02914669"
		_, err := ttService.Package(orderNumber, packageNumber)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("ok")
		}
	}
}
