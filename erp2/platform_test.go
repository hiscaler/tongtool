package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"testing"
)

func TestService_Platforms(t *testing.T) {
	platforms, err := ttService.Platforms()
	if err != nil {
		t.Errorf("ttService.Platforms error: %s", err.Error())
	} else {
		fmt.Println(jsonx.ToJson(platforms, "[]"))
	}
}
