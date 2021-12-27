package erp2

// https://open.tongtool.com/apiDoc.html#/?docId=8e80fde6a4824b288d17bc04be8f4ef6

// PurchaseSuggestion 采购建议
type PurchaseSuggestion struct {
	CaculateDate           string `json:"caculateDate"`           // 采购建议计算时间
	CurrStockQuantity      int    `json:"currStockQuantity"`      // 可用库存数
	DailySales             int    `json:"dailySales"`             // 日均销量
	DevliveryDays          int    `json:"devliveryDays"`          // 安全交期
	GoodsIdKey             string `json:"goodsIdKey"`             // 商品id key
	GoodsSku               string `json:"goodsSku"`               // 商品sku
	IntransitStockQuantity int    `json:"intransitStockQuantity"` // 在途库存数
	ProposalQuantity       int    `json:"proposalQuantity"`       // 采购建议数量
	SaleAvg15              int    `json:"saleAvg15"`              // 15天销量
	SaleAvg30              int    `json:"saleAvg30"`              // 30天销量
	SaleAvg7               int    `json:"saleAvg7"`               // 7天销量
	UnpickingQuantity      int    `json:"unpickingQuantity"`      // 订单未配货数量
	WarehouseIdKey         int    `json:"warehouseIdKey"`         // 仓库id key
	WarehouseName          int    `json:"warehouseName"`          // 仓库名称
}
