package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
)

// Paypal 付款记录查询
// https://open.tongtool.com/apiDoc.html#/?docId=e4d63dc26dc649ef9d57cdb7da1fdc7d

// PaypalTransaction Paypal 付款单数据
type PaypalTransaction struct {
	Amount              float64 `json:"amount"`              //	总金额
	Currency            string  `json:"currency"`            //	币种
	Fee                 float64 `json:"fee"`                 //	Paypal交易费
	Name                string  `json:"name"`                //	收款人姓名
	PayerPaypalEmail    string  `json:"payerPaypalEmail"`    //	收款人邮箱
	PaymentDate         string  `json:"paymentDate"`         //	付款时间
	PaypalTransactionId string  `json:"paypalTransactionId"` //	Paypal记录ID
}

type PaypalTransactionsQueryParams struct {
	Paging
	MerchantId          string `json:"merchantId"`                    // 商户ID
	PaypalTransactionId string `json:"paypalTransactionId,omitempty"` // Paypal记录ID
	UpdatedDateBegin    string `json:"updatedDateBegin,omitempty"`    // 下载更新起始时间
	UpdatedDateEnd      string `json:"updatedDateEnd,omitempty"`      // 下载更新结束时间
}

func (m PaypalTransactionsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UpdatedDateBegin, validation.When(m.UpdatedDateBegin != "", validation.Date(constant.DatetimeFormat).Error("无效的下载更新起始时间格式"))),
		validation.Field(&m.UpdatedDateEnd, validation.When(m.UpdatedDateEnd != "", validation.Date(constant.DatetimeFormat).Error("无效的下载更新结束时间格式"))),
	)
}

// PaypalTransactions Paypal 付款记录查询
func (s service) PaypalTransactions(params PaypalTransactionsQueryParams) (items []PaypalTransaction, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = keyx.Generate(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = json.Unmarshal(b, &items); e == nil {
				return
			} else {
				s.tongTool.Logger.Printf(`cache data unmarshal error
 DATA: %s
ERROR: %s
`, string(b), e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("get cache %s error: %s", cacheKey, e.Error())
		}
	}
	items = make([]PaypalTransaction, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []PaypalTransaction `json:"array"`
			PageNo   int                 `json:"pageNo"`
			PageSize int                 `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/paypalQueryQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
			isLastPage = len(items) < params.PageSize
		}
	} else {
		if e := json.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	if err != nil {
		return
	}

	if s.tongTool.EnableCache && len(items) > 0 {
		if b, e := json.Marshal(&items); e == nil {
			e = s.tongTool.Cache.Set(cacheKey, b)
			if e != nil {
				s.tongTool.Logger.Printf("set cache %s error: %s", cacheKey, e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("items marshal error: %s", err.Error())
		}
	}
	return
}
