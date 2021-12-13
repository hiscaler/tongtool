package tongtool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"strconv"
	"time"
)

type queryDefaultValues struct {
	PageNo   int
	PageSize int
}

type TongTool struct {
	Client             *resty.Client
	MerchantId         string
	Logger             *log.Logger
	QueryDefaultValues queryDefaultValues
}

func NewTongTool(appKey, appSecret string, debug bool) *TongTool {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	client := resty.New()
	if debug {
		client.SetDebug(true).EnableTrace()
	}
	tokenResponse := struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Datas   string      `json:"datas"`
		Other   interface{} `json:"other"`
	}{}
	_, err := client.R().
		SetResult(&tokenResponse).
		Get(fmt.Sprintf("https://open.tongtool.com/open-platform-service/devApp/appToken?accessKey=%s&secretAccessKey=%s", appKey, appSecret))
	if err != nil || !tokenResponse.Success {
		logger.Panic("Get token failed.")
	}

	timestamp := int(time.Now().UnixNano() / 1e6)
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("app_token%stimestamp%d%s", tokenResponse.Datas, timestamp, appSecret)))
	sign := hex.EncodeToString(h.Sum(nil))
	partnerResponse := struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Datas   []struct {
			PartnerOpenId string `json:"partnerOpenId"`
		} `json:"datas"`
		Other interface{} `json:"other"`
	}{}
	_, err = client.R().
		SetResult(&partnerResponse).
		Get(fmt.Sprintf("https://open.tongtool.com/open-platform-service/partnerOpenInfo/getAppBuyerList?app_token=%s&timestamp=%d&sign=%s",
			tokenResponse.Datas,
			timestamp,
			sign,
		))

	if err != nil || !partnerResponse.Success || len(partnerResponse.Datas) == 0 {
		log.Panicf("Get partnerOpenId failed, error: %s", err.Error())
	}

	merchantId := partnerResponse.Datas[0].PartnerOpenId
	client.SetBaseURL("https://open.tongtool.com/api-service").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"api_version":  "3.0",
		}).
		SetQueryParams(map[string]string{
			"app_token": tokenResponse.Datas,
			"sign":      sign,
			"timestamp": strconv.Itoa(timestamp),
		}).
		SetTimeout(10 * time.Second)

	return &TongTool{
		Client:     client,
		MerchantId: merchantId,
		Logger:     logger,
		QueryDefaultValues: queryDefaultValues{
			PageNo:   1,
			PageSize: 100,
		},
	}
}
