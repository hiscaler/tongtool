package tongtool

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"strconv"
	"time"
)

// 通途返回代码
const (
	OK                   = 200
	SignError            = 519
	TokenExpiredError    = 523
	UnauthorizedError    = 524
	TooManyRequestsError = 526
	AccountExpiredError  = 999999
)

var ErrNotFound = errors.New("tongtool: not found")

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
	logger := log.New(os.Stdout, "TongTool", log.LstdFlags)
	client := resty.New().SetBaseURL("https://open.tongtool.com/open-platform-service")
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
		Get(fmt.Sprintf("/devApp/appToken?accessKey=%s&secretAccessKey=%s", appKey, appSecret))
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
		Get(fmt.Sprintf("/partnerOpenInfo/getAppBuyerList?app_token=%s&timestamp=%d&sign=%s",
			tokenResponse.Datas,
			timestamp,
			sign,
		))

	if err != nil || !partnerResponse.Success || len(partnerResponse.Datas) == 0 {
		log.Panicf("Get partnerOpenId failed, error: %s", err.Error())
	}

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
		MerchantId: partnerResponse.Datas[0].PartnerOpenId,
		Logger:     logger,
		QueryDefaultValues: queryDefaultValues{
			PageNo:   1,
			PageSize: 100,
		},
	}
}

// WithCache 激活缓存
func (t *TongTool) WithCache(v bool) (err error) {
	if v {
		// Active
		if t.Cache == nil {
			cache, e := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
			if e == nil {
				t.ActivateCache = true
				t.Cache = cache
			} else {
				t.Logger.Printf("cache: active cache error: %s", e.Error())
				err = e
			}
		} else {
			t.ActivateCache = true
		}
	} else {
		// Close
		t.ActivateCache = false
	}

	return
}

// ErrorWrap 错误包装
func ErrorWrap(code int, defaultMessage string) error {
	if code == OK {
		return nil
	}

	msg := ""
	switch code {
	case SignError:
		msg = "签名错误"
	case TokenExpiredError:
		msg = "Token 已过期"
	case UnauthorizedError:
		msg = "未授权的请求，请确认应用是否勾选对应接口"
	case TooManyRequestsError:
		msg = "接口请求超请求次数限额"
	case AccountExpiredError:
		msg = "账号已过期"
	default:
		if defaultMessage == "" {
			msg = fmt.Sprintf("未知的错误代码：%d", code)
		} else {
			msg = fmt.Sprintf("%d: %s", code, defaultMessage)
		}
	}
	return errors.New(msg)
}
