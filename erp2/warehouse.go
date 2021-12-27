package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
	"strings"
)

type Warehouse struct {
	WarehouseId   string `json:"warehouseId"`   // 仓库id
	WarehouseCode string `json:"warehouseCode"` // 仓库代码
	WarehouseName string `json:"warehouseName"` // 仓库名称
	Status        string `json:"status"`        // 仓库状态：0-失效1-有效
	StatusBoolean bool   `json:"StatusBoolean"` // 仓库状态布尔值（返回仓库状态布尔值，方便调用者判断）
}

type WarehouseQueryParams struct {
	MerchantId    string `json:"merchantId"`              // 商户ID
	PageNo        int    `json:"pageNo,omitempty"`        // 查询页数
	PageSize      int    `json:"pageSize,omitempty"`      // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
	WarehouseName string `json:"warehouseName,omitempty"` // 仓库名称
	WarehouseId   string `json:"warehouseId,omitempty"`   // 仓库id（通途无此参数）
}

type warehousesResult struct {
	result
	Datas struct {
		Array    []Warehouse `json:"array"`
		PageNo   int         `json:"pageNo"`
		PageSize int         `json:"pageSize"`
	} `json:"datas,omitempty"`
}

// Warehouses 查询仓库列表
// https://open.tongtool.com/apiDoc.html#/?docId=cdb49c57add3448daf1f4cd0fad40bef
func (s service) Warehouses(params WarehouseQueryParams) (items []Warehouse, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = s.tongTool.QueryDefaultValues.PageNo
	}
	if params.PageSize <= 0 {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]Warehouse, 0)
	res := warehousesResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/warehouseQuery")
	if err == nil {
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
	}
	return
}

// Warehouse 查询指定仓库
func (s service) Warehouse(params WarehouseQueryParams) (item Warehouse, err error) {
	const (
		searchById        = "id"
		searchByName      = "name"
		searchByIdAndName = "id.name"
	)
	searchBy := ""
	if params.WarehouseId != "" && params.WarehouseName != "" {
		searchBy = searchByIdAndName
	} else if params.WarehouseId != "" {
		searchBy = searchById
	} else if params.WarehouseName != "" {
		searchBy = searchByName
	}

	if searchBy == "" {
		err = errors.New("invalid query params")
		return
	}

	exists := false
	for {
		items := make([]Warehouse, 0)
		isLastPage := false
		items, isLastPage, err = s.Warehouses(params)
		if err == nil {
			if len(items) == 0 {
				err = errors.New("not found")
			} else {
				for _, warehouse := range items {
					switch searchBy {
					case searchByIdAndName:
						exists = strings.EqualFold(warehouse.WarehouseId, params.WarehouseId) && strings.EqualFold(warehouse.WarehouseName, params.WarehouseName)
					case searchById:
						exists = strings.EqualFold(warehouse.WarehouseId, params.WarehouseId)
					case searchByName:
						exists = strings.EqualFold(warehouse.WarehouseName, params.WarehouseName)
					}
					if exists {
						item = warehouse
						break
					}
				}
				if exists {
					break
				}
			}
		}
		if isLastPage || exists || err != nil {
			break
		}
		params.PageNo++
	}

	if err == nil && !exists {
		err = errors.New("not found")
	}

	return
}

// 物流渠道

type ShippingMethod struct {
	ShippingMethodId            string `json:"shippingMethodId"`            // 渠道ID
	ShippingMethodShortname     string `json:"shippingMethodShortname"`     // 渠道名称
	ShippingMethodStatus        string `json:"shippingMethodStatus"`        // 渠道状态0-失效1-有效
	ShippingMethodStatusBoolean bool   `json:"shippingMethodStatusBoolean"` // 渠道状态
	CarrierName                 string `json:"carrierName"`                 // 物流商简称
	CarrierStatus               string `json:"carrierStatus"`               // 物流商状态0-失效1-有效
	CarrierStatusBoolean        bool   `json:"CarrierStatusBoolean"`        // 物流商状态
	WarehouseId                 string `json:"warehouseId"`                 // 仓库id
	WarehouseName               string `json:"warehouseName"`               // 仓库名称
}

type ShippingMethodQueryParams struct {
	MerchantId  string `json:"merchantId"`         // 商户ID
	PageNo      int    `json:"pageNo,omitempty"`   // 查询页数
	PageSize    int    `json:"pageSize,omitempty"` // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
	WarehouseId string `json:"warehouseId"`        // 仓库id
}

type shippingMethodsResult struct {
	result
	Datas struct {
		Array    []ShippingMethod `json:"array"`
		PageNo   int              `json:"pageNo"`
		PageSize int              `json:"pageSize"`
	} `json:"datas,omitempty"`
}

// ShippingMethods 仓库物流渠道查询
// https://open.tongtool.com/apiDoc.html#/?docId=9ed7d6c3e7c44e498c0d43329d5a443b
func (s service) ShippingMethods(params ShippingMethodQueryParams) (items []ShippingMethod, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]ShippingMethod, 0)
	res := shippingMethodsResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/getShippingMethod")
	if err == nil {
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
	}
	return
}
