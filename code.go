package tongtool

import (
	"errors"
	"fmt"
)

const (
	SignError            = 519
	TokenExpiredError    = 523
	UnauthorizedError    = 524
	TooManyRequestsError = 526
)

var codeTexts map[int]string

func init() {
	codeTexts = map[int]string{
		SignError:            "签名错误",
		TokenExpiredError:    "Token 已过期",
		UnauthorizedError:    "未授权的请求，请确认应用是否勾选对应接口",
		TooManyRequestsError: "接口请求超请求次数限额\n",
	}
}

func CodeText(code int) string {
	if text, ok := codeTexts[code]; ok {
		return text
	} else {
		return fmt.Sprintf("未知的错误：%d", code)
	}
}

func Successful(code int) (ok bool, err error) {
	if code == 200 {
		ok = true
	} else {
		err = errors.New(CodeText(code))
	}
	return
}
