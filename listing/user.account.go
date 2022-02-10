package listing

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
)

// 保存用户店铺信息
// https://open.tongtool.com/apiDoc.html#/?docId=858f1b7cf16d435bb83ca35fcc4071f1

type UserAccountInfo struct {
	CreatedBy     string `json:"createdBy"`     // 创建人
	CreatedDate   string `json:"createdDate"`   // 创建时间
	IsDelete      bool   `json:"isDelete"`      // 是否删除
	MerchantId    string `json:"merchantId"`    // 商户编号
	SaleAccountId string `json:"saleAccountId"` // 店铺账户编号
	UpdatedBy     string `json:"updatedBy"`     // 修改人
	UpdatedDate   string `json:"updatedDate"`   // 修改时间
	UserId        string `json:"userId"`        // 用户Id
}

type UpsertUserAccountRequest struct {
	IsCovered           bool              `json:"isCovered"`           // 是否覆盖原有店铺权限
	UserAccountInfoList []UserAccountInfo `json:"userAccountInfoList"` // 用户店铺信息列表
	MerchantId          string            `json:"merchantId"`          // 商户编号
	UserId              string            `json:"userId"`              // 用户Id
}

func (m UpsertUserAccountRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserAccountInfoList, validation.Required.Error("店铺信息列表不能为空")),
		validation.Field(&m.UserId, validation.Required.Error("用户Id不能为空")),
	)
}

// SaveUserAccount 保存用户店铺信息
func (s service) SaveUserAccount(req UpsertUserAccountRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
		Datas string `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/saleAccount/saveUserAccount")
	if err == nil {
		if resp.IsSuccess() {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	return err
}
