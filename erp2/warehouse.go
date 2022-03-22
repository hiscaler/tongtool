package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"strings"
)

type Warehouse struct {
	WarehouseId   string `json:"warehouseId"`   // 仓库id
	WarehouseCode string `json:"warehouseCode"` // 仓库代码
	WarehouseName string `json:"warehouseName"` // 仓库名称
	Status        string `json:"status"`        // 仓库状态：0-失效1-有效
	StatusBoolean bool   `json:"statusBoolean"` // 仓库状态布尔值（返回仓库状态布尔值，方便调用者判断）
}

type WarehousesQueryParams struct {
	Paging
	MerchantId    string `json:"merchantId"`              // 商户ID
	WarehouseName string `json:"warehouseName,omitempty"` // 仓库名称
}

// Warehouses 查询仓库列表
// https://open.tongtool.com/apiDoc.html#/?docId=cdb49c57add3448daf1f4cd0fad40bef
func (s service) Warehouses(params WarehousesQueryParams) (items []Warehouse, isLastPage bool, err error) {
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
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []Warehouse `json:"array"`
			PageNo   int         `json:"pageNo"`
			PageSize int         `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/warehouseQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
			for i, item := range items {
				items[i].StatusBoolean = item.Status == "1"
			}
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

	if s.tongTool.EnableCache {
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

// Warehouse 查询指定仓库
func (s service) Warehouse(id string) (item Warehouse, exists bool, err error) {
	if id == "" {
		err = errors.New("无效的 id 参数值")
		return
	}

	params := WarehousesQueryParams{}
	params.PageNo = 1
	for {
		var items []Warehouse
		isLastPage := false
		items, isLastPage, err = s.Warehouses(params)
		if err == nil {
			if len(items) == 0 {
				err = tongtool.ErrNotFound
			} else {
				for _, warehouse := range items {
					if strings.EqualFold(warehouse.WarehouseId, id) {
						exists = true
						item = warehouse
						break
					}
				}
			}
		}
		if isLastPage || exists || err != nil {
			break
		}
		params.PageNo++
	}

	if err == nil && !exists {
		err = tongtool.ErrNotFound
	}
	return
}

// 物流渠道

type WarehouseShippingMethod struct {
	CarrierName                 string `json:"carrierName"`                 // 物流商简称
	CarrierStatus               string `json:"carrierStatus"`               // 物流商状态（0：失效、1：有效）
	CarrierStatusBoolean        bool   `json:"carrierStatusBoolean"`        // 物流商状态
	ShippingMethodId            string `json:"shippingMethodId"`            // 渠道ID
	ShippingMethodShortname     string `json:"shippingMethodShortname"`     // 渠道名称
	ShippingMethodStatus        string `json:"shippingMethodStatus"`        // 渠道状态（0：失效、1：有效）
	ShippingMethodStatusBoolean bool   `json:"shippingMethodStatusBoolean"` // 渠道状态
	WarehouseId                 string `json:"warehouseId"`                 // 仓库id
	WarehouseName               string `json:"warehouseName"`               // 仓库名称
}

type WarehouseShippingMethodsQueryParams struct {
	Paging
	MerchantId  string `json:"merchantId"`  // 商户ID
	WarehouseId string `json:"warehouseId"` // 仓库id
}

func (m WarehouseShippingMethodsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.WarehouseId, validation.Required.Error("仓库 ID 不能为空")),
	)
}

// WarehouseShippingMethods 仓库物流渠道查询
// https://open.tongtool.com/apiDoc.html#/?docId=9ed7d6c3e7c44e498c0d43329d5a443b
func (s service) WarehouseShippingMethods(params WarehouseShippingMethodsQueryParams) (items []WarehouseShippingMethod, isLastPage bool, err error) {
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
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []WarehouseShippingMethod `json:"array"`
			PageNo   int                       `json:"pageNo"`
			PageSize int                       `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/getShippingMethod")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
			for i, item := range items {
				items[i].ShippingMethodStatusBoolean = item.ShippingMethodStatus == "1"
				items[i].CarrierStatusBoolean = item.CarrierStatus == "1"
			}
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
