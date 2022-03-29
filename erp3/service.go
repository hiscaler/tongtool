package erp3

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	Products(params ProductsQueryParams) (items []Product, nextToken string, isLastPage bool, err error)                // 商品列表
	UserTicket(ticket string) (u User, refreshTicket string, expire int, err error)                                     // 根据 ticket 获取员工信息
	Suppliers(params SuppliersQueryParams) (items []Supplier, nextToken string, isLastPage bool, err error)             // 供应商列表
	WarehouseAreas(params WarehouseAreasQueryParams) (items []WarehouseArea, err error)                                 // 仓库分区关系
	SaveThirdAccounts(req UpdateThirdAccountRequest) error                                                              // 保存第三方帐号信息
	StockInSheets(params StockInSheetsQueryParams) (items []StockInSheet, nextToken string, isLastPage bool, err error) // 入库单列表
	AddShippingPackage(req AddShippingPackageRequest) (packages []ShippingPackage, err error)                           // 出库单交运
}

func NewService(tt *tongtool.TongTool) Service {
	tt.QueryDefaultValues.PageSize = 500
	return service{tt}
}
