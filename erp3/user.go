package erp3

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
	"strings"
)

// 根据 ticket 获取员工信息
// https://open.tongtool.com/apiDoc.html#/?docId=6b99851aaad0420699e300ff7889709d
// post /openapi/userInfo/userByTicket

type User struct {
	ActivationDate string `json:"activationDate"` // 邮件激活时间
	CreateTime     string `json:"createTime"`     // 生成日期
	CreatedBy      string `json:"createdBy"`      // 创建人
	CreatedDate    string `json:"createdDate"`    // 创建时间
	Email          string `json:"email"`          // email
	EmailBind      int    `json:"emailBind"`      // 是否已绑定邮箱（null or 0：未绑定、1：已绑定）
	Extension      string `json:"extension"`      // 分机
	IsAdmin        string `json:"isAdmin"`        // 是否商户管理员
	LastLoginTime  string `json:"lastLoginTime"`  // 最后登录时间
	MerchantId     string `json:"merchantId"`     // 商户编号
	Mobile         string `json:"mobile"`         // 手机
	MobileBind     int    `json:"mobileBind"`     // 是否已绑定手机（null or 0：未绑定、1：已绑定）
	Openid         string `json:"openid"`         // 微信用户唯一标识
	RoleName       string `json:"roleName"`       // 角色名称
	Status         string `json:"status"`         // 状态(是否生效)
	Telephone      string `json:"telephone"`      // 电话号码
	UpdatedBy      string `json:"updatedBy"`      // 修改人
	UpdatedDate    string `json:"updatedDate"`    // 修改时间
	UserId         string `json:"userId"`         // 用户代码
	UserName       string `json:"userName"`       // 姓名
	WxBind         int    `json:"wxBind"`         // 是否已绑定微信（null or 0：未绑定、1：已绑定）
	ZipCode        string `json:"zipCode"`        // 电话区号
}

type UserTicket struct {
	Expire        int    `json:"expireln"`      // refreshTicket 失效时间
	RefreshTicket string `json:"refreshTicket"` // 访问令牌
	UserInfo      User   `json:"userInfo"`      // 用户信息
}

type UserQueryParams struct {
	MerchantId string `json:"merchantId"` // 商户号
	Ticket     string `json:"ticket"`     // ticket
}

func (m UserQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Ticket, validation.Required.Error("ticket 不能为空")),
	)
}

func (s service) UserTicket(ticket string) (u User, refreshTicket string, expire int, err error) {
	params := UserQueryParams{
		MerchantId: s.tongTool.MerchantId,
		Ticket:     strings.TrimSpace(ticket),
	}
	res := struct {
		tongtool.Response
		Datas UserTicket `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/userInfo/userByTicket")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			u = res.Datas.UserInfo
			refreshTicket = res.Datas.RefreshTicket
			expire = res.Datas.Expire
		}
	} else {
		if e := json.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}
