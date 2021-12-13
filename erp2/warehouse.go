package erp2

import (
	"errors"
	"github.com/hiscaler/tongtool"
)

type Warehouse struct {
	WarehouseId   string `json:"warehouseId"`   // 仓库id
	WarehouseCode string `json:"warehouseCode"` // 仓库代码
	WarehouseName string `json:"warehouseName"` // 仓库名称
	Status        string `json:"status"`        // 仓库状态：0-失效1-有效
}

type WarehouseQueryParams struct {
	MerchantId    string `json:"merchantId"`              // 商户ID
	PageNo        int    `json:"pageNo,omitempty"`        // 查询页数
	PageSize      int    `json:"pageSize,omitempty"`      // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
	WarehouseName string `json:"warehouseName,omitempty"` // 仓库名称
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
				isLastPage = len(items) < params.PageSize
			}
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}
