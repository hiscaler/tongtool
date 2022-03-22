package erp3

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
)

// 保存第三方帐号信息
// https://open.tongtool.com/apiDoc.html#/?docId=59a641205b3f43dd87ba96d515849b73

type ThirdAccount struct {
	Code  string `json:"accountCode"` // Code
	Token string `json:"accessToken"` // Token
}

type UpdateThirdAccountRequest struct {
	MerchantId string         `json:"merchantId"`           // 商户编号
	Accounts   []ThirdAccount `json:"accountCodeTokenList"` // 第三方帐号信息
}

func (m UpdateThirdAccountRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Accounts, validation.Required.Error("第三方帐号信息列表不能为空")),
	)
}

// SaveThirdAccounts 保存第三方帐号信息
func (s service) SaveThirdAccounts(req UpdateThirdAccountRequest) error {
	var err error
	res := struct {
		tongtool.Response
		Datas string `json:"datas"`
	}{}
	req.MerchantId = s.tongTool.MerchantId
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/wmsCommon/saveThirdAccount")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		err = tongtool.ErrorWrap(res.Code, res.Message)
	} else {
		if e := json.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}
