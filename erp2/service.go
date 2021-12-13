package erp2

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	FBAOrders(params FBAOrderQueryParams) (items []FBAOrder, isLastPage bool, err error)                 // FBA 订单列表
	Orders(params OrderQueryParams) (items []Order, isLastPage bool, err error)                          // 订单列表
	Order(id string) (Order, error)                                                                      // 单个订单
	Products(params ProductQueryParams) (items []Product, isLastPage bool, err error)                    // 商品列表
	Product(typ int, skus []string, isAlias bool) (product Product, err error)                           // 单个商品
	ProductExists(typ int, skus []string, isAlias bool) bool                                             // 商品是否存在
	CreateProduct(req CreateProductRequest) error                                                        // 创建商品
	UpdateProduct(req UpdateProductRequest) error                                                        // 更新商品
	Packages(params PackageQueryParams) (items []Package, isLastPage bool, err error)                    // 包裹列表
	Package(orderNumber, packageNumber string) (pkg Package, err error)                                  // 单个包裹
	Suppliers(params SuppliersQueryParams) (items []Supplier, isLastPage bool, err error)                // 供应商列表
	PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error) // 采购单列表
	CreatePurchaseOrder(params CreatePurchaseOrderRequest) (number string, err error)                    // 创建采购单
	SaleAccounts(params SaleAccountQueryParams) (items []SaleAccount, isLastPage bool, err error)        // 创建采购单
	Stocks(params StockQueryParams) (items []Stock, isLastPage bool, err error)                          // 创建采购单
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
