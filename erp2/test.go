package erp2

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/tongtool"
	"os"
)

var tt *tongtool.TongTool
var ttService Service

func init() {
	type config struct {
		Debug     bool   `json:"debug"`
		AppKey    string `json:"appKey"`
		AppSecret string `json:"appSecret"`
	}
	b, err := os.ReadFile("../config/config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	tt = tongtool.NewTongTool(c.AppKey, c.AppSecret, c.Debug)
	ttService = NewService(*tt)
}
