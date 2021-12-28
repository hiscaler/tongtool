package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"strings"
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
	orderId := "L-M20211221152430918"
	packageId := "P02914669"
	_, err := ttService.Package(orderId, packageId)
	if err != nil {
		t.Error(err)
	}
}

func TestService_PackageWithCache(t *testing.T) {
	tt, ttService := newTestTongTool()
	tt.SwitchCache(true)
	for i := 0; i <= 400; i++ {
		orderId := "L-M20211221152430918"
		packageId := "P02914669"
		p, err := ttService.Package(orderId, packageId)
		if err != nil {
			t.Errorf("ttService.Package error: %s", err.Error())
		} else if !strings.EqualFold(p.PackageId, packageId) {
			t.Errorf("package.package id %s not equal %s", p.PackageId, packageId)
		} else {
			fmt.Println("ok")
		}
	}
}
