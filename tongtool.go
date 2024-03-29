package tongtool

import (
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/gox/cryptox"
	"github.com/hiscaler/tongtool/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// 通途返回代码
// https://open.tongtool.com/apiDoc.html#/?docId=72398ba411bd11eab4c00050568e1e6a
const (
	OK                     = 200    // 无错误
	SignError              = 519    // 签名错误
	TokenExpiredError      = 523    // Token 已过期
	UnauthorizedError      = 524    // 未授权的请求，请确认应用是否勾选对应接口
	InvalidParametersError = 525    // 无效的参数
	TooManyRequestsError   = 526    // 接口请求超请求次数限额
	SystemError            = 527    // 系统错误
	AccountExpiredError    = 999999 // 账号已过期
)

func init() {
	extra.RegisterFuzzyDecoders()
}

var ErrNotFound = errors.New("tongtool: not found")

type queryDefaultValues struct {
	PageNo   int // 当前页
	PageSize int // 每页数据量
}

type TongTool struct {
	application        app                // 认证后的应用数据
	assetSaveDir       string             // 资源保存地址
	Debug              bool               // 是否调试模式
	Client             *resty.Client      // HTTP 客户端
	MerchantId         string             // 商户 ID
	Logger             *log.Logger        // 日志
	EnableCache        bool               // 是否激活缓存
	Cache              *bigcache.BigCache // 缓存
	QueryDefaultValues queryDefaultValues // 查询默认值

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
	Timestamp          int64   `json:"timestamp"`
	Sign               string  `json:"sign"`
	Valid              bool    `json:"valid"`
}

func NewTongTool(config config.Config) *TongTool {
	logger := log.New(os.Stdout, "[ TongTool ] ", log.LstdFlags|log.Llongfile)
	ttInstance := &TongTool{
		Debug:  config.Debug,
		Logger: logger,
		QueryDefaultValues: queryDefaultValues{
			PageNo:   1,
			PageSize: 100,
		},
	}
	if application, e := auth(config.AppKey, config.AppSecret, config.Debug); e == nil {
		application.AppTokenExpireDate /= 1000
		ttInstance.application = application
		ttInstance.MerchantId = application.PartnerOpenId
	} else {
		logger.Printf("auth error: %s", e.Error())
	}
	timeoutSeconds := config.Timeout
	if timeoutSeconds <= 0 {
		timeoutSeconds = 10
	}
	retryCount := config.RetryCount
	if retryCount < 0 {
		retryCount = 0
	}
	retryWaitTime := config.RetryWaitTime
	if retryWaitTime <= 0 {
		retryWaitTime = 12
	}
	retryMaxWaitTime := config.RetryMaxWaitTime
	if retryMaxWaitTime <= 0 {
		retryMaxWaitTime = 60
	}
	client := resty.New().
		SetDebug(config.Debug).
		SetBaseURL("https://open.tongtool.com/api-service").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"api_version":  "3.0",
		}).
		SetTimeout(time.Duration(timeoutSeconds) * time.Second).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			if !ttInstance.application.Valid || ttInstance.application.AppTokenExpireDate <= time.Now().Unix()-1800 {
				application, e := auth(config.AppKey, config.AppSecret, config.Debug)
				if e != nil {
					logger.Printf("auth error: %s", e.Error())
					return e
				}
				ttInstance.MerchantId = application.PartnerOpenId
				application.AppTokenExpireDate /= 1000
				ttInstance.application = application
			}
			request.SetQueryParams(map[string]string{
				"app_token": ttInstance.application.AppToken,
				"sign":      ttInstance.application.Sign,
				"timestamp": strconv.FormatInt(ttInstance.application.Timestamp, 10),
			})
			return nil
		}).
		SetRetryCount(retryCount).
		SetRetryWaitTime(time.Duration(retryWaitTime) * time.Second).
		SetRetryMaxWaitTime(time.Duration(retryMaxWaitTime) * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			if response == nil {
				return false
			}

			retry := response.StatusCode() == http.StatusTooManyRequests
			if !retry {
				r := struct{ Code int }{}
				retry = jsoniter.Unmarshal(response.Body(), &r) == nil && r.Code == TooManyRequestsError
			}
			if retry {
				text := response.Request.URL
				if err != nil {
					text += fmt.Sprintf(", error: %s", err.Error())
				}
				logger.Printf("Retry request: %s", text)
			}
			return retry
		})
	if config.ForceWaiting {
		client.SetRetryAfter(func(client *resty.Client, response *resty.Response) (time.Duration, error) {
			seconds := 0
			if response != nil {
				retry := response.StatusCode() == http.StatusTooManyRequests
				if !retry {
					r := struct{ Code int }{}
					retry = jsoniter.Unmarshal(response.Body(), &r) == nil && r.Code == TooManyRequestsError
				}
				if retry {
					seconds = 60 - time.Now().Second()%60 // 最多等待下一分钟到目前的秒数
				}
			}
			if seconds == 0 {
				return 0, nil
			}
			logger.Printf("Waiting %d seconds...", seconds)
			return time.Duration(seconds) * time.Second, nil
		})
	}
	if config.Debug {
		client.EnableTrace()
	}
	if config.EnableCache {
		if err := ttInstance.SwitchCache(true); err != nil {
			logger.Printf("Cache: %s", err.Error())
		}
	}
	jsoniter.RegisterTypeDecoderFunc("float64", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {
		case jsoniter.StringValue:
			var t float64
			v := strings.TrimSpace(iter.ReadString())
			if v != "" {
				var err error
				if t, err = strconv.ParseFloat(v, 64); err != nil {
					iter.Error = err
					return
				}
			}
			*((*float64)(ptr)) = t
		case jsoniter.BoolValue:
			// support bool to float64
			if iter.ReadBool() {
				*((*float64)(ptr)) = 1
			} else {
				*((*float64)(ptr)) = 0
			}
		case jsoniter.NilValue:
			iter.Skip()
			*((*float64)(ptr)) = 0
		default:
			*((*float64)(ptr)) = iter.ReadFloat64()
		}
	})
	jsoniter.RegisterTypeDecoderFunc("bool", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {
		case jsoniter.StringValue:
			var t bool
			v := strings.TrimSpace(iter.ReadString())
			if v != "" {
				var err error
				if t, err = strconv.ParseBool(strings.ToLower(v)); err != nil {
					iter.Error = err
					return
				}
			}
			*((*bool)(ptr)) = t
		case jsoniter.NumberValue:
			if v, err := iter.ReadNumber().Int64(); err != nil {
				iter.Error = err
				return
			} else {
				*((*bool)(ptr)) = v > 0
			}
		case jsoniter.NilValue:
			iter.Skip()
			*((*bool)(ptr)) = false
		default:
			*((*bool)(ptr)) = iter.ReadBool()
		}
	})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	ttInstance.Client = client
	return ttInstance
}

// SwitchCache 激活缓存
func (t *TongTool) SwitchCache(v bool) (err error) {
	if v {
		// Active
		if t.Cache == nil {
			var cacheConfig bigcache.Config
			if t.Debug {
				cacheConfig = bigcache.DefaultConfig(10 * time.Minute)
			} else {
				cacheConfig = bigcache.Config{
					Shards:             1024,
					LifeWindow:         10 * time.Minute,
					CleanWindow:        1 * time.Second,
					MaxEntriesInWindow: 1000 * 10 * 60,
					MaxEntrySize:       500,
					StatsEnabled:       false,
					Verbose:            false,
					Hasher:             nil,
					HardMaxCacheSize:   0,
					Logger:             nil,
				}
			}
			cacheConfig.Logger = t.Logger
			if cache, e := bigcache.NewBigCache(cacheConfig); e == nil {
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

func (t *TongTool) SetAssetSaveDir(dir string) *TongTool {
	t.assetSaveDir = dir
	return t
}

func (t *TongTool) GetAssetSaveDir() string {
	return t.assetSaveDir
}

func auth(appKey, appSecret string, debug bool) (application app, err error) {
	client := resty.New().SetDebug(debug).SetBaseURL("https://open.tongtool.com/open-platform-service")
	if debug {
		client.EnableTrace()
	}
	application = app{}
	tokenResponse := struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Datas   string      `json:"datas"`
		Other   interface{} `json:"other"`
	}{}
	resp, err := client.R().
		SetResult(&tokenResponse).
		Get(fmt.Sprintf("/devApp/appToken?accessKey=%s&secretAccessKey=%s", appKey, appSecret))
	if err != nil {
		return
	}
	if resp.IsError() {
		return application, ErrorWrap(resp.StatusCode(), resp.String())
	} else if !tokenResponse.Success {
		return application, ErrorWrap(tokenResponse.Code, tokenResponse.Message)
	}

	timestamp := time.Now().Unix()
	appToken := tokenResponse.Datas
	sign := cryptox.Md5(fmt.Sprintf("app_token%stimestamp%d%s", appToken, timestamp, appSecret))
	appResponse := struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Datas   []app       `json:"datas"`
		Other   interface{} `json:"other"`
	}{}
	resp, err = client.R().
		SetResult(&appResponse).
		Get(fmt.Sprintf("/partnerOpenInfo/getAppBuyerList?app_token=%s&timestamp=%d&sign=%s", appToken, timestamp, sign))
	if err != nil {
		return
	}

	if resp.IsError() {
		return application, ErrorWrap(resp.StatusCode(), resp.String())
	} else if !appResponse.Success || len(appResponse.Datas) == 0 {
		return application, ErrorWrap(appResponse.Code, appResponse.Message)
	}

	application = appResponse.Datas[0]
	application.Valid = true
	application.Timestamp = timestamp
	application.Sign = sign
	return
}

// Response Normal API response
type Response struct {
	Code    int         `json:"code"`
	Datas   interface{} `json:"datas"`
	Message string      `json:"message"`
	Others  interface{} `json:"others"`
}

// ErrorWrap 错误包装
func ErrorWrap(code int, message string) error {
	if code == OK {
		return nil
	}

	message = strings.TrimSpace(message)
	if message == "" {
		switch code {
		case SignError:
			message = "签名错误"
		case TokenExpiredError:
			message = "Token 已过期"
		case UnauthorizedError:
			message = "未授权的请求，请确认应用是否勾选对应接口"
		case SystemError:
			message = "系统错误"
		case InvalidParametersError:
			message = "无效的参数"
		case TooManyRequestsError:
			message = "接口请求超请求次数限额"
		case AccountExpiredError:
			message = "账号已过期"
		default:
			message = "未知错误"
		}
	}

	return fmt.Errorf("%d: %s", code, message)
}
