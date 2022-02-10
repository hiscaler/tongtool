package erp3

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/config"
	"os"
	"testing"
)

var ttInstance *tongtool.TongTool
var ttService Service

func TestMain(m *testing.M) {
	b, err := os.ReadFile("../config/config_test.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var c config.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	ttInstance = tongtool.NewTongTool(c)
	ttService = NewService(ttInstance)
	m.Run()
}

func TestService_Products(t *testing.T) {
	params := ProductsQueryParams{}
	_, _, _, err := ttService.Products(params)
	if err != nil {
		t.Errorf("ttService.Products error: %s", err.Error())
	}
}
