package listing

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

const (
	WarehouseStatusInvalid = "0"
	WarehouseStatusValid   = "1"
)

type Warehouse struct {
	AutoSyncTime        string `json:"autoSyncTime"`        // 自动同步时间如,12:25
	City                string `json:"city"`                // 城市
	Country             string `json:"country"`             // 国家
	CreatedBy           string `json:"createdBy"`           // 创建人
	CreatedDate         string `json:"createdDate"`         // 创建时间
	IsSyncErpPercentage string `json:"isSyncErpPercentage"` // 是否同步ERP仓库库存数的xx%;Y,N
	MerchantId          string `json:"merchantId"`          // 商户编号
	PostalCode          string `json:"postalCode"`          // 邮编
	Status              string `json:"status"`              // 状态（1生效，0失效）
	SyncErpPercentage   int    `json:"syncErpPercentage"`   // 同步ERP仓库库存数的xx
	SyncErpType         string `json:"syncErpType"`         // A-只同步可用库存数,B-同步可用库存数+在途库存
	SyncTime            string `json:"syncTime"`            // 上次同步ERP库存时间
	UpdatedBy           string `json:"updatedBy"`           // 修改人
	UpdatedDate         string `json:"updatedDate"`         // 修改时间
	WarehouseId         string `json:"warehouseId"`         // 仓库编号
	WarehouseName       string `json:"warehouseName"`       // 仓库名称
}

// 查询仓库列表
// https://open.tongtool.com/apiDoc.html#/?docId=9e7d98f617d84f2e93754bf78cfeac1c

type WarehousesQueryParams struct {
	MerchantId    string `json:"merchantId"`              // 商户编号
	Status        string `json:"status,omitempty"`        // 状态（1生效，0失效）
	WarehouseName string `json:"warehouseName,omitempty"` // 仓库名称
}

func (m WarehousesQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.When(m.Status != "", validation.In(WarehouseStatusInvalid, WarehouseStatusValid).Error("无效的仓库状态"))),
	)
}

func (s service) Warehouses(params WarehousesQueryParams) (items []Warehouse, err error) {
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
	items = make([]Warehouse, 0)
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
		Post("/openapi/tongtool/listing/warehouse/getWarehouse")
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
