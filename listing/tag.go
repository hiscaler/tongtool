package listing

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

const (
	ProductTag                  = "product"           // 产品资料标签
	ProductEbayDraftTag         = "ebayDraft"         // eBay草稿标签
	ProductEbayListingTag       = "ebayListing"       // eBay在线标签
	ProductAliexpressDraftTag   = "aliexpressDraft"   // 速卖通草稿标签
	ProductAliexpressListingTag = "aliexpressListing" // 速卖通在线标签
)

type Tag struct {
	LabelId   string `json:"labelId"`   // 标签ID
	LabelName string `json:"labelName"` // 标签名称
	LabelType string `json:"labelType"` // 标签类别
}

// 标签列表
// https://open.tongtool.com/apiDoc.html#/?docId=f22b1937adf04312974634495a9bbb6e

type TagsQueryParams struct {
	LabelId    string `json:"labelId,omitempty"`   // 标签ID
	LabelName  string `json:"labelName,omitempty"` // 标签名称
	LabelType  string `json:"labelType,omitempty"` // 标签类别
	MerchantId string `json:"merchantId"`          // 商户号
}

func (m TagsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.LabelType, validation.When(m.LabelType != "", validation.In(ProductTag, ProductEbayDraftTag, ProductEbayListingTag, ProductAliexpressDraftTag, ProductAliexpressListingTag).Error("无效的标签类别"))),
	)
}

func (s service) Tags(params TagsQueryParams) (items []Tag, err error) {
	params.MerchantId = s.tongTool.MerchantId
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
	items = make([]Tag, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []Tag `json:"array"`
			PageNo   int   `json:"pageNo"`
			PageSize int   `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/listing/productTag/getProductTag")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
			}
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	if err == nil && s.tongTool.EnableCache && len(items) > 0 {
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

// 添加标签
// https://open.tongtool.com/apiDoc.html#/?docId=2daeddee878b457ea736a4de7b67717d

type CreateTagRequest struct {
	LabelId    string `json:"labelId"`    // 标签ID
	LabelName  string `json:"labelName"`  // 标签名称
	LabelType  string `json:"labelType"`  // 标签类别
	MerchantId string `json:"merchantId"` // 商户号
}

func (m CreateTagRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.LabelName, validation.Required.Error("标签名称不能为空")),
		validation.Field(&m.LabelType, validation.Required.Error("标签类别不能为空")),
	)
}

func (s service) CreateTag(req CreateTagRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/productTag/createProductTag")
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

// 替换标签库标签
// https://open.tongtool.com/apiDoc.html#/?docId=453ecba3472f41da85eb76bbef08da4a

type UpdateTagRequest struct {
	DestinationLabelName string `json:"destinationLabelName"` // 目的标签名称
	LabelType            string `json:"labelType"`            // 标签类型
	MerchantId           string `json:"merchantId"`           // 商户号
	OriginalLabelName    string `json:"originalLabelName"`    // 原标签名称
}

func (m UpdateTagRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DestinationLabelName, validation.Required.Error("目的标签名称不能为空")),
		validation.Field(&m.LabelType, validation.Required.Error("标签类别不能为空")),
		validation.Field(&m.OriginalLabelName, validation.Required.Error("原标签名称不能为空")),
	)
}

func (s service) UpdateTag(req UpdateTagRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/productTag/replaceLabelLibrary")
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

// 删除标签
// https://open.tongtool.com/apiDoc.html#/?docId=9ebbd46fa32642c58e56cb9e7383050a

type DeleteTagRequest struct {
	LabelId    string `json:"labelId"`    // 标签ID
	LabelType  string `json:"labelType"`  // 标签类别
	MerchantId string `json:"merchantId"` // 商户号
}

func (m DeleteTagRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.LabelId, validation.Required.Error("标签 ID 不能为空")),
	)
}

func (s service) DeleteTag(req DeleteTagRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/productTag/removeProductTag")
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
