package erp2

import (
	"tongtool"
)

type service struct {
	tongTool tongtool.TongTool
}

type result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Service interface {
	FBAOrders(param FBAOrderQueryParam) (items []FBAOrder, isLastPage bool, err error)                      // FBA 订单列表
	Orders(param OrderQueryParam) (items []Order, isLastPage bool, err error)                               // 订单列表
	Order(id string) (Order, error)                                                                         // 单个订单
	Products(param ProductQueryParam) (items []Product, isLastPage bool, err error)                         // 商品列表
	Product(typ int, skus []string, isAlias bool) (product Product, err error)                              // 单个商品
	ProductExists(typ int, skus []string, isAlias bool) bool                                                // 商品是否存在
	CreateProduct(req CreateProductRequest) error                                                           // 创建商品
	UpdateProduct(req UpdateProductRequest) error                                                           // 更新商品
	Packages(params PackageQueryParam) (items []Package, isLastPage bool, err error)                        // 包裹列表
	Package(orderNumber, packageNumber string) (pkg Package, err error)                                     // 单个包裹
	Suppliers(params SuppliersQueryParam) (suppliers []Supplier, isLastPage bool, err error)                // 供应商列表
	PurchaseOrders(params PurchaseOrdersQueryParam) (suppliers []PurchaseOrder, isLastPage bool, err error) // 采购单列表
	CreatePurchaseOrder(params CreatePurchaseOrderRequest) (number string, err error)                       // 创建采购单
}

func NewService(tt tongtool.TongTool) Service {
	return service{tt}
}
