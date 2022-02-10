package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"strings"
	"testing"
)

func TestService_Packages(t *testing.T) {
	params := PackagesQueryParams{
		AssignTimeFrom: "2021-10-01 00:00:00",
		AssignTimeTo:   "2021-12-30 23:59:59",
		PackageStatus:  PackageStatusWaitDeliver,
	}
	packages, _, err := ttService.Packages(params)
	if err == nil {
		fmt.Println(jsonx.ToJson(packages, "[]"))
	} else {
		t.Error(err)
	}
}

func TestService_Package(t *testing.T) {
	orderId := "L-M20211221152430918"
	packageId := "P02914669"
	_, err := ttService.Package(orderId, packageId)
	if err != nil {
		t.Error(err)
	}
}

func TestService_PackageWithCache(t *testing.T) {
	ttInstance.SwitchCache(true)
	times := 400
	n := 0
	for i := 0; i < times; i++ {
		orderId := "L-M20211221152430918"
		packageId := "P02914669"
		p, err := ttService.Package(orderId, packageId)
		if err != nil {
			t.Errorf("ttService.Package error: %s", err.Error())
		} else if !strings.EqualFold(p.PackageId, packageId) {
			t.Errorf("package.PackageId %s not equal %s", p.PackageId, packageId)
		} else {
			n++
		}
	}
	if n != times {
		t.Errorf("except %d times, actual %d times", times, n)
	}
}

func TestService_PackageDeliver(t *testing.T) {
	req := PackageDeliverRequest{
		DeliverInfos: []PackageDeliverItem{
			{
				RelatedNo: "P02912767",
				Volume:    PackageDeliverItemVolume{1, 2, 3},
			},
			{
				RelatedNo: "P02913843",
				Volume:    PackageDeliverItemVolume{4, 5, 6},
			},
		},
		WarehouseName: "test",
	}
	err := ttService.PackageDeliver(req)
	if err != nil {
		t.Error(err)
	}
}
