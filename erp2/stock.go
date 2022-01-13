package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	"time"
)

// 库存查询

type GoodsShelfStockItem struct {
	AvailableStockQuantity       int    `json:"availableStockQuantity"`       // 可用库存数
	DefectsStockQuantity         int    `json:"defectsStockQuantity"`         // 故障品库存数
	GoodsDetailId                string `json:"goodsDetailId"`                // 通途货品ID
	GoodsShelfCode               string `json:"goodsShelfCode"`               // 货位编号
	GoodsShelfId                 string `json:"goodsShelfId"`                 // 货位ID
	WaitingShipmentStockQuantity int    `json:"waitingShipmentStockQuantity"` // 待发库存数
}

type Stock struct {
	AvailableStockQuantity       int                   `json:"availableStockQuantity"`       // 可用库存数
	CargoSpace                   string                `json:"cargoSpace"`                   // 货位
	DefectsStockQuantity         int                   `json:"defectsStockQuantity"`         // 故障品库存数
	FirstShippingFeeUnit         float64               `json:"firstShippingFeeUnit"`         // 头程运费
	FirstTariff                  float64               `json:"firstTariff"`                  // 头程报关费
	GoodsAvgCost                 float64               `json:"goodsAvgCost"`                 // 货品平均成本
	GoodsCurCost                 float64               `json:"goodsCurCost"`                 // 货品当前成本
	GoodsIdKey                   string                `json:"goodsIdKey"`                   // 通途商品id key
	GoodsShelfStockList          []GoodsShelfStockItem `json:"goodsShelfStockList"`          // 货位库存列表，多货位才会有值
	GoodsSKU                     string                `json:"goodsSku"`                     // 商品sku
	IntransitStockQuantity       int                   `json:"intransitStockQuantity"`       // 在途库存数
	OtherFee                     float64               `json:"otherFee"`                     // 头程其他费用
	SafetyStock                  int                   `json:"safetyStock"`                  // 安全库存数
	WaitingShipmentStockQuantity int                   `json:"waitingShipmentStockQuantity"` // 待发库存数
	WarehouseIdKey               string                `json:"warehouseIdKey"`               // 通途仓库id key
	WarehouseName                string                `json:"warehouseName"`                // 仓库名称
}

type StockQueryParams struct {
	MerchantId      string   `json:"merchantId"`                // 商户ID
	PageNo          int      `json:"pageNo,omitempty"`          // 查询页数
	PageSize        int      `json:"pageSize,omitempty"`        // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
	SKUs            []string `json:"skus,omitempty"`            // SKU 列表
	UpdatedDateFrom string   `json:"updatedDateFrom,omitempty"` // 更新开始时间
	UpdatedDateTo   string   `json:"updatedDateTo,omitempty"`   // 更新结束时间
	WarehouseName   string   `json:"warehouseName"`             // 仓库名称
}

// Stocks 库存列表
// https://open.tongtool.com/apiDoc.html#/?docId=9aaf6b145a014060b3b3f669b0487096
func (s service) Stocks(params StockQueryParams) (items []Stock, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
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
	items = make([]Stock, 0)
	res := struct {
		result
		Datas struct {
			Array    []Stock `json:"array"`
			PageNo   int     `json:"pageNo"`
			PageSize int     `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/stocksQuery")
	if err == nil {
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

// 库存变动日志查询

type StockChangeLog struct {
	AvailableStockQuantity int    `json:"availableStockQuantity"` // 当前可用库存
	ChangeQuantity         int    `json:"changeQuantity"`         // 变动数量；正数增加，负数减少
	GoodsSKU               string `json:"goodsSku"`               // 商品sku
}

type StockChangeLogQueryParams struct {
	MerchantId      string   `json:"merchantId"`         // 商户 ID
	PageNo          int      `json:"pageNo,omitempty"`   // 查询页数
	PageSize        int      `json:"pageSize,omitempty"` // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
	SKUs            []string `json:"skus,omitempty"`     // SKU 列表
	UpdatedDateFrom string   `json:"updatedDateFrom"`    // 变动起始时间；统计此时间以后的库存变动，只能输入距当前时间7天内的值
	WarehouseName   string   `json:"warehouseName"`      // 仓库名称
}

func (m StockChangeLogQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UpdatedDateFrom, validation.Required.Error("变动起始时间不能为空"), validation.Date(constant.DatetimeFormat).Error("变动起始时间格式错误"), validation.By(func(value interface{}) error {
			t, err := time.Parse(constant.DatetimeFormat, value.(string))
			if err != nil {
				return err
			}
			if time.Now().Sub(t).Hours() > 24*7 {
				return errors.New("变动起始时间只能输入距当前时间 7 天内的值")
			}
			return nil
		})),
		validation.Field(&m.WarehouseName, validation.Required.Error("仓库名称不能为空")),
	)
}

// StockChangeLogs 库存变动查询
// https://open.tongtool.com/apiDoc.html#/?docId=bd0971f61f2449eaa9752c7be779afa0
func (s service) StockChangeLogs(params StockChangeLogQueryParams) (items []StockChangeLog, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
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
	items = make([]StockChangeLog, 0)
	res := struct {
		result
		Datas struct {
			Array    []StockChangeLog `json:"array"`
			PageNo   int              `json:"pageNo"`
			PageSize int              `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/stocksChangeQuery")
	if err == nil {
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
