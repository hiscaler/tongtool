package erp3

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
)

// 仓库分区关系
// https://open.tongtool.com/apiDoc.html#/?docId=3fe29cb2f9e04eacad5842f33fde78bd

type WarehouseArea struct {
	CreatedBy                 string `json:"createdBy"`                 // 创建人
	CreatedTime               string `json:"createdTime"`               // 创建时间
	MerchantId                string `json:"merchantId"`                // 商户id
	Name                      string `json:"name"`                      // 名称
	Remark                    string `json:"remark"`                    // 备注
	Status                    int    `json:"status"`                    // 状态（0：停用、1：启用、默认启用）
	UpdatedBy                 string `json:"updatedBy"`                 // 更新人
	UpdatedTime               string `json:"updatedTime"`               // 更新时间
	WmsWarehouseAreaProgramId string `json:"wmsWarehouseAreaProgramId"` // 仓库分区方案流水号
}

type WarehouseAreasQueryParams struct {
	MerchantId   string   `json:"merchantId"`   // 商户号
	WarehouseIds []string `json:"warehouseIds"` // 仓库 IDs
}

func (m WarehouseAreasQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.WarehouseIds, validation.Required.Error("仓库 ID 列表不能为空")),
	)
}

// WarehouseAreas 仓库分区关系
func (s service) WarehouseAreas(params WarehouseAreasQueryParams) (items []WarehouseArea, err error) {
	if err = params.Validate(); err != nil {
		return
	}

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
	res := struct {
		tongtool.Response
		Datas []WarehouseArea `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/wmsWarehouseAreaRelated/query")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas
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
