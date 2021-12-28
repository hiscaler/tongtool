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
	Debug              bool
	Client             *resty.Client
	MerchantId         string
	Logger             *log.Logger
	EnableCache        bool
	Cache              *bigcache.BigCache
	QueryDefaultValues queryDefaultValues
	application        app
}

type app struct {
	TokenId            string  `json:"tokenId"`
	DevId              string  `json:"devId"`
	DevAppId           string  `json:"devAppId"`
	AccessKey          string  `json:"accessKey"`
	AppToken           string  `json:"appToken"`
	AppTokenExpireDate int64   `json:"appTokenExpireDate"`
	PartnerOpenId      string  `json:"partnerOpenId"`
	UserOpenId         string  `json:"userOpenId"`
	UserName           string  `json:"userName"`
	BuyDate            int     `json:"buyDate"`
	Price              float64 `json:"price"`
	CreatedDate        int     `json:"createdDate"`
	CreatedBy          string  `json:"createdBy"`
	UpdatedDate        int     `json:"updatedDate"`
	UpdatedBy          string  `json:"updatedBy"`
	Timestamp          int     `json:"timestamp"`
	Sign               string  `json:"sign"`
	Valid              bool    `json:"valid"`
}

func NewTongTool(appKey, appSecret string, debug bool) *TongTool {
	logger := log.New(os.Stdout, "TongTool", log.LstdFlags)
	ttInstance := &TongTool{
		Debug:  debug,
		Logger: logger,
		QueryDefaultValues: queryDefaultValues{
			PageNo:   1,
			PageSize: 100,
		},
	}
	if application, e := auth(appKey, appSecret, debug); e == nil {
		application.AppTokenExpireDate /= 1000
		ttInstance.application = application
		ttInstance.MerchantId = application.PartnerOpenId
	} else {
		logger.Printf("auth error: %s", e.Error())
	}
	client := resty.New().
		SetBaseURL("https://open.tongtool.com/api-service").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"api_version":  "3.0",
		}).
		SetTimeout(10 * time.Second)
	if debug {
		client.SetDebug(true).EnableTrace()
	}
	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		if !ttInstance.application.Valid || ttInstance.application.AppTokenExpireDate <= time.Now().Unix()-1800 {
			application, e := auth(appKey, appSecret, debug)
			if e != nil {
				logger.Printf("auth error: %s", e.Error())
				return e
			}
			ttInstance.MerchantId = application.PartnerOpenId
			application.AppTokenExpireDate /= 1000
			ttInstance.application = application
		}
		client.SetQueryParams(map[string]string{
			"app_token": ttInstance.application.AppToken,
			"sign":      ttInstance.application.Sign,
			"timestamp": strconv.Itoa(ttInstance.application.Timestamp),
		})
		return nil
	})
	ttInstance.Client = client
	return ttInstance
}

// SwitchCache 激活缓存
func (t *TongTool) SwitchCache(v bool) (err error) {
	if v {
		// Active
		if t.Cache == nil {
			var config bigcache.Config
			if t.Debug {
				config = bigcache.DefaultConfig(10 * time.Minute)
			} else {
				config = bigcache.Config{
					Shards:             1024,
					LifeWindow:         10 * time.Minute,
					CleanWindow:        1 * time.Second,
					MaxEntriesInWindow: 1000 * 10 * 60,
					MaxEntrySize:       500,
					StatsEnabled:       false,
					Verbose:            true,
					Hasher:             nil,
					HardMaxCacheSize:   0,
					Logger:             nil,
				}
			}
			config.Logger = t.Logger
			cache, e := bigcache.NewBigCache(config)
			if e == nil {
				t.EnableCache = true
				t.Cache = cache
			} else {
				t.Logger.Printf("cache: active cache error: %s", e.Error())
				err = e
			}
		} else {
			t.EnableCache = true
		}
	} else {
		// Close
		t.EnableCache = false
	}

	return
}

func auth(appKey, appSecret string, debug bool) (application app, err error) {
	client := resty.New().
		SetBaseURL("https://open.tongtool.com/open-platform-service")
	if debug {
		client.SetDebug(true).EnableTrace()
	}
	application = app{}
	tokenResponse := struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Datas   string      `json:"datas"`
		Other   interface{} `json:"other"`
	}{}
	_, err = client.R().
		SetResult(&tokenResponse).
		Get(fmt.Sprintf("/devApp/appToken?accessKey=%s&secretAccessKey=%s", appKey, appSecret))
	if err != nil {
		return
	}
	if !tokenResponse.Success {
		err = errors.New("get token failed")
	}

	timestamp := int(time.Now().UnixNano() / 1e6)
	appToken := tokenResponse.Datas
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("app_token%stimestamp%d%s", appToken, timestamp, appSecret)))
	sign := hex.EncodeToString(h.Sum(nil))
	appResponse := struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Datas   []app       `json:"datas"`
		Other   interface{} `json:"other"`
	}{}
	_, err = client.R().
		SetResult(&appResponse).
		Get(fmt.Sprintf("/partnerOpenInfo/getAppBuyerList?app_token=%s&timestamp=%d&sign=%s", appToken, timestamp, sign))
	if err != nil {
		return
	}
	if !appResponse.Success || len(appResponse.Datas) == 0 {
		err = errors.New("getAppBuyerList data invalid")
		return
	}

	application = appResponse.Datas[0]
	application.Valid = true
	application.Timestamp = timestamp
	application.Sign = sign
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
