package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
	"strings"
)

type Warehouse struct {
	WarehouseId   string `json:"warehouseId"`   // 仓库id
	WarehouseCode string `json:"warehouseCode"` // 仓库代码
	WarehouseName string `json:"warehouseName"` // 仓库名称
	Status        string `json:"status"`        // 仓库状态：0-失效1-有效
	TTEnabled     bool   `json:"tt_enabled"`    // 激活（返回仓库状态布尔值，方便调用者判断）
}

type WarehouseQueryParams struct {
	MerchantId    string `json:"merchantId"`              // 商户ID
	PageNo        int    `json:"pageNo,omitempty"`        // 查询页数
	PageSize      int    `json:"pageSize,omitempty"`      // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
	WarehouseName string `json:"warehouseName,omitempty"` // 仓库名称
	WarehouseId   string `json:"warehouseId,omitempty"`   // 仓库名称
}

type warehousesResult struct {
	result
	Datas struct {
		Array    []Warehouse `json:"array"`
		PageNo   int         `json:"pageNo"`
		PageSize int         `json:"pageSize"`
	}
}

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
			if err = tongtool.HasError(res.Code); err == nil {
				items = res.Datas.Array
				for i, item := range items {
					items[i].TTEnabled = item.Status == "1"
				}
				isLastPage = len(items) < params.PageSize
			}
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}

func (s service) Warehouse(params WarehouseQueryParams) (item Warehouse, err error) {
	byId := params.WarehouseId != ""
	byName := params.WarehouseName != ""
	if !byId && !byName {
		err = errors.New("invalid query params")
		return
	}

	for {
		items := make([]Warehouse, 0)
		isLastPage := false
		items, isLastPage, err = s.Warehouses(params)
		if err == nil {
			if len(items) == 0 {
				err = errors.New("not found")
			} else {
				exists := false
				for _, warehouse := range items {
					if byId {
						if strings.EqualFold(warehouse.WarehouseId, params.WarehouseId) {
							item = warehouse
							exists = true
							break
						}
					} else if byName {
						if strings.EqualFold(warehouse.WarehouseName, params.WarehouseName) {
							item = warehouse
							exists = true
							break
						}
					}
					if exists {
						break
					}
				}
				if exists {
					break
				}
			}
		}
		if err != nil || isLastPage {
			break
		}
		params.PageNo++
	}
	return
}
