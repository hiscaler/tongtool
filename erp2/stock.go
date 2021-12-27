package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
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
	DefectsStockQuantity         int                   `json:"cargoSpace"`                   // 故障品库存数
	FirstShippingFeeUnit         float64               `json:"firstShippingFeeUnit"`         // 头程运费
	FirstTariff                  float64               `json:"firstTariff"`                  // 头程报关费
	GoodsAvgCost                 float64               `json:"goodsAvgCost"`                 // 货品平均成本
	GoodsCurCost                 float64               `json:"goodsCurCost"`                 // 货品当前成本
	GoodsIdKey                   float64               `json:"goodsIdKey"`                   // 通途商品id key
	goodsShelfStockList          []GoodsShelfStockItem `json:"goodsShelfStockList"`          // 货位库存列表，多货位才会有值
	GoodsSKU                     string                `json:"goodsSku"`                     // 商品sku
	IntransitStockQuantity       string                `json:"intransitStockQuantity"`       // 在途库存数
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
	return
}
