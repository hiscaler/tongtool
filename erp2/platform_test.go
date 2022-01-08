package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_Platforms(t *testing.T) {
	platforms, err := ttService.Platforms()
	if err != nil {
		t.Errorf("ttService.Platforms error: %s", err.Error())
	} else {
		fmt.Println(cast.ToJson(platforms))
	}
}
