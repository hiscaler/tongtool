package erp3

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/config"
	"os"
)

func newTestTongTool() (*tongtool.TongTool, Service) {
	b, err := os.ReadFile("../config/config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var c config.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	instance := tongtool.NewTongTool(c)
	return instance, NewService(instance)
}
