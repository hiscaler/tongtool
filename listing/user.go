package listing

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hiscaler/tongtool"
)

const (
	UserOperatingTypeAdd    = "add"
	UserOperatingTypeEdit   = "edit"
	UserOperatingTypeUpdate = "update"
)

type UpsertUserRequest struct {
	Email         string `json:"email"`         // email (新增时必填)
	ListingStatus string `json:"listingStatus"` // 刊登系统状态(是否生效) 1或0,启用、停用时必填
	MerchantId    string `json:"merchantId"`    // 商户编号
	Mobile        string `json:"mobile"`        // 手机
	OperatingType string `json:"operatingType"` // 操作类型（add新增，edit编辑，update启用/停用）
	Password      string `json:"password"`      // 密码（新增时必填）
	UserId        string `json:"userId"`        // 用户Id（编辑、修改时必填）
	UserName      string `json:"userName"`      // 姓名（新增时必填）
}

func (m UpsertUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.When(m.Email != "", is.Email.Error("无效的邮箱地址"))),
		validation.Field(&m.ListingStatus, validation.When(m.OperatingType == UserOperatingTypeUpdate, validation.Required.Error("刊登系统状态不能为空"), validation.In("0", "1").Error("无效的刊登系统状态"))),
		validation.Field(&m.OperatingType, validation.In(UserOperatingTypeAdd, UserOperatingTypeEdit, UserOperatingTypeUpdate).Error("无效的操作类型")),
	)
}

// UpsertUser 保存用户信息
func (s service) UpsertUser(req UpsertUserRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		result
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/user/saveUserInfo")
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