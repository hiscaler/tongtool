package listing

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hiscaler/tongtool"
	jsoniter "github.com/json-iterator/go"
)

const (
	UserOperationTypeAdd    = "add"    // 新增
	UserOperationTypeEdit   = "edit"   // 编辑
	UserOperationTypeUpdate = "update" // 启用/停用
)

// 保存用户信息
// https://open.tongtool.com/apiDoc.html#/?docId=7e44eb4fd4d647919fbca632cbae1638

type UpsertUserRequest struct {
	Email         string `json:"email,omitempty"`         // email (新增时必填)
	ListingStatus string `json:"listingStatus,omitempty"` // 刊登系统状态(是否生效) 1或0,启用、停用时必填
	MerchantId    string `json:"merchantId"`              // 商户编号
	Mobile        string `json:"mobile,omitempty"`        // 手机
	OperatingType string `json:"operatingType,omitempty"` // 操作类型（add：新增、edit：编辑、update：启用/停用）
	Password      string `json:"password,omitempty"`      // 密码（新增时必填）
	UserId        string `json:"userId,omitempty"`        // 用户Id（编辑、修改时必填）
	UserName      string `json:"userName,omitempty"`      // 姓名（新增时必填）
}

func (m UpsertUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email,
			validation.When(m.Email != "", is.EmailFormat.Error("无效的邮箱格式")),
			validation.When(m.OperatingType == UserOperationTypeAdd, validation.Required.Error("邮箱地址不能为空")),
		),
		validation.Field(&m.ListingStatus, validation.When(m.OperatingType == UserOperationTypeUpdate,
			validation.Required.Error("刊登系统状态不能为空"),
			validation.In("0", "1").Error("无效的刊登系统状态")),
		),
		validation.Field(&m.OperatingType, validation.In(UserOperationTypeAdd, UserOperationTypeEdit, UserOperationTypeUpdate).Error("无效的操作类型")),
		validation.Field(&m.Password, validation.When(m.OperatingType == UserOperationTypeAdd, validation.Required.Error("密码不能为空"))),
		validation.Field(&m.UserId, validation.When(m.OperatingType == UserOperationTypeEdit, validation.Required.Error("用户 ID 不能为空"))),
		validation.Field(&m.UserName, validation.When(m.OperatingType == UserOperationTypeAdd, validation.Required.Error("姓名不能为空"))),
	)
}

// UpsertUser 保存用户信息
func (s service) UpsertUser(req UpsertUserRequest) error {
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
		Post("/openapi/tongtool/listing/user/saveUserInfo")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		err = tongtool.ErrorWrap(res.Code, res.Message)
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}
