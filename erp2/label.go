package erp2

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

// Label 特性标签
type Label struct {
	LabelName  string `json:"labelName"`  //	标签名称
	ProductNum int    `json:"productNum"` // 正在使用商品数量
}

// 新增特性标签

type CreateLabelRequest struct {
	LabelName  string `json:"labelName"`  // 标签名称
	MerchantId string `json:"merchantId"` // 商户 ID
}

func (m CreateLabelRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.LabelName, validation.Required.Error("标签名称不能为空")),
	)
}

// CreateLabel 创建标签
// https://open.tongtool.com/apiDoc.html#/?docId=81b03b1725ed4306812715a63de9f081
func (s service) CreateLabel(req CreateLabelRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := tongtool.Response{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/createLabel")
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

// 查询标签

type LabelsQueryParams struct {
	Paging
	LabelName  string `json:"labelName,omitempty"` // 标签名称
	MerchantId string `json:"merchantId"`          // 商户 ID
}

func (m LabelsQueryParams) Validate() error {
	return nil
}

// Labels 根据指定参数查询标签列表
// https://open.tongtool.com/apiDoc.html#/?docId=4243c256b768470b9152610744f72764
func (s service) Labels(params LabelsQueryParams) (items []Label, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = keyx.Generate(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = jsoniter.Unmarshal(b, &items); e == nil {
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
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []Label `json:"array"`
			PageNo   int     `json:"pageNo"`
			PageSize int     `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/LabelQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
			isLastPage = len(items) <= params.PageSize
		}
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	if err != nil {
		return
	}

	if s.tongTool.EnableCache && len(items) > 0 {
		if b, e := jsoniter.Marshal(&items); e == nil {
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

// LabelExists 根据 标签名称查询标签是否存在
func (s service) LabelExists(name string) (exists bool, err error) {
	labels, _, err := s.Labels(LabelsQueryParams{LabelName: name})
	if err != nil {
		return
	}

	for _, label := range labels {
		if strings.EqualFold(label.LabelName, name) {
			return true, nil
		}
	}
	return false, nil
}
