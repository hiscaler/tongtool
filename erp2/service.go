package erp2

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	AmazonAccountSites(params AmazonAccountSiteQueryParams) (items []string, isLastPage bool, err error)                                       // 查询亚马逊账号对应的站点
	FBAOrders(params FBAOrderQueryParams) (items []FBAOrder, isLastPage bool, err error)                                                       // FBA 订单列表
	ShopifyOrders(params ShopifyOrderQueryParams) (items []ShopifyOrder, isLastPage bool, err error)                                           // Shopify 订单列表
	CreateOrder(req CreateOrderRequest) (orderId, orderNumber string, err error)                                                               // 手工创建订单
	Orders(params OrderQueryParams) (items []Order, isLastPage bool, err error)                                                                // 订单列表
	Order(id string) (item Order, err error)                                                                                                   // 单个订单
	CancelOrder(req CancelOrderRequest) (results []OrderCancelResult, err error)                                                               // 作废订单
	Products(params ProductQueryParams) (items []Product, isLastPage bool, err error)                                                          // 商品列表
	Product(typ string, sku string, isAlias bool) (item Product, err error)                                                                    // 单个商品
	ProductExists(typ string, sku string, isAlias bool) bool                                                                                   // 商品是否存在
	CreateProduct(req CreateProductRequest) error                                                                                              // 创建商品
	UpdateProduct(req UpdateProductRequest) error                                                                                              // 更新商品
	Packages(params PackageQueryParams) (items []Package, isLastPage bool, err error)                                                          // 包裹列表
	Package(orderNumber, packageNumber string) (item Package, err error)                                                                       // 单个包裹
	PackageDeliver(req PackageDeliverRequest) error                                                                                            // 执行包裹发货
	Suppliers(params SuppliersQueryParams) (items []Supplier, isLastPage bool, err error)                                                      // 供应商列表
	PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error)                                       // 采购单列表
	CreatePurchaseOrder(params CreatePurchaseOrderRequest) (number string, err error)                                                          // 创建采购单
	PurchaseOrderArrival(req PurchaseOrderArrivalRequest) error                                                                                // 采购单到货
	PurchaseOrderStockIn(req PurchaseOrderStockInRequest) error                                                                                // 采购单入库
	PurchaseOrderStockInLogs(params PurchaseOrderLogQueryParams) (items []PurchaseOrderLog, isLastPage bool, err error)                        // 采购单入库查询
	SaleAccounts(params SaleAccountQueryParams) (items []SaleAccount, isLastPage bool, err error)                                              // 商户账号列表
	Stocks(params StockQueryParams) (items []Stock, isLastPage bool, err error)                                                                // 库存列表
	StockChangeLogs(params StockChangeLogQueryParams) (items []StockChangeLog, isLastPage bool, err error)                                     //  库存变动查询
	Warehouses(params WarehouseQueryParams) (items []Warehouse, isLastPage bool, err error)                                                    // 仓库列表
	Warehouse(id string) (item Warehouse, err error)                                                                                           // 仓库列表
	ShippingMethods(params ShippingMethodQueryParams) (items []ShippingMethod, isLastPage bool, err error)                                     // 仓库物流渠道列表
	TrackingNumbers(params TrackingNumberQueryParams) (items []TrackingNumber, isLastPage bool, err error)                                     // 订单物流单号列表
	Platforms() (items []Platform, err error)                                                                                                  // 平台及站点信息
	PurchaseSuggestionTemplates(params PurchaseSuggestionTemplateQueryParams) (items []PurchaseSuggestionTemplate, isLastPage bool, err error) // 采购建议模板列表
	PurchaseSuggestions(params PurchaseSuggestionQueryParams) (items []PurchaseSuggestion, isLastPage bool, err error)                         // 采购建议列表
	QuotePrices(params QuotedPriceQueryParams) (items []QuotedPrice, isLastPage bool, err error)                                               // 供应商报价查询
	AfterSales(params AfterSaleQueryParams) (items []AfterSale, isLastPage bool, err error)                                                    // 售后单信息查询
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
