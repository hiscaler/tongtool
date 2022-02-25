package listing

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

// 获取产品类目
// https://open.tongtool.com/apiDoc.html#/?docId=11a5118bb70642f198a7acca0c0b56a2

type CategoriesQueryParams struct {
	CategoryId       string `json:"categoryId,omitempty"`       // 类目编号
	CategoryName     string `json:"categoryName,omitempty"`     // 类目名称
	MerchantId       string `json:"merchantId"`                 // 商户编号
	ParentCategoryId string `json:"parentCategoryId,omitempty"` // 父类目编号
}

type Category struct {
	CategoryCode     string     `json:"categoryCode"`     // 类目 CODE
	CategoryId       string     `json:"categoryId"`       // 类目编号
	CategoryName     string     `json:"categoryName"`     // 类目名称
	ChildList        []Category `json:"childList"`        // 子类目集合
	IsRoot           string     `json:"isRoot"`           // 是否根类目
	ParentCategoryId string     `json:"parentCategoryId"` // 父类目编号
}

// Categories 根据指定参数查询商品列表
// https://open.tongtool.com/apiDoc.html#/?docId=919e8fff6c8047deb77661f4d8c92a3a
func (s service) Categories(params CategoriesQueryParams) (items []Category, err error) {
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
	items = make([]Category, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []Category `json:"array"`
			PageNo   int        `json:"pageNo"`
			PageSize int        `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/listing/productCategory/getProductCategory")
	if err != nil {
		return
	}

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

// CUDCategory Create/Update/Delete Category
type CUDCategory struct {
	CategoryId       string `json:"categoryId"`       // 类目编号
	CategoryName     string `json:"categoryName"`     // 类目名称
	MerchantId       string `json:"merchantId"`       // 商户编号
	ParentCategoryId string `json:"parentCategoryId"` // 父类目编号
}

// 添加产品类目
// https://open.tongtool.com/apiDoc.html#/?docId=94ef3350cd064550a7cdb7b88f008b54

type CreateCategoryRequest struct {
	CUDCategory
}

func (m CreateCategoryRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CategoryId, validation.Required.Error("类目编号不能为空")),
		validation.Field(&m.CategoryName, validation.Required.Error("类目名称不能为空")),
	)
}

// CreateCategory 添加产品类目
func (s service) CreateCategory(req CreateCategoryRequest) error {
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
		Post("/openapi/tongtool/listing/productCategory/createProductCategory")
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

// 修改产品类目
// https://open.tongtool.com/apiDoc.html#/?docId=afb5b4128bc94fb291aec0f5e9310f83

type UpdateCategoryRequest struct {
	CUDCategory
}

func (m UpdateCategoryRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CategoryId, validation.Required.Error("类目编号不能为空")),
		validation.Field(&m.CategoryName, validation.Required.Error("类目名称不能为空")),
	)
}

// UpdateCategory 修改产品类目
func (s service) UpdateCategory(req UpdateCategoryRequest) error {
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
		Post("/openapi/tongtool/listing/productCategory/changeProductCategory")
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

// 删除类目
// https://open.tongtool.com/apiDoc.html#/?docId=76b970744cf64824b5093f59ecb8f3ba

type DeleteCategoryRequest struct {
	CUDCategory
}

func (m DeleteCategoryRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CategoryId, validation.Required.Error("类目编号不能为空")),
		validation.Field(&m.CategoryName, validation.Required.Error("类目名称不能为空")),
	)
}

// DeleteCategory 删除产品类目
func (s service) DeleteCategory(req DeleteCategoryRequest) error {
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
		Post("/openapi/tongtool/listing/productCategory/delProductCategory")
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
