package erp2

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	AmazonAccountSites(params AmazonAccountSitesQueryParams) (items []string, isLastPage bool, err error)                                       // 查询亚马逊账号对应的站点
	FBAOrders(params FBAOrdersQueryParams) (items []FBAOrder, isLastPage bool, err error)                                                       // FBA 订单列表
	ShopifyOrders(params ShopifyOrdersQueryParams) (items []ShopifyOrder, isLastPage bool, err error)                                           // Shopify 订单列表
	CreateOrder(req CreateOrderRequest) (orderId, orderNumber string, err error)                                                                // 手工创建订单
	UpdateOrder(req UpdateOrderRequest) error                                                                                                   // 更新订单
	Orders(params OrdersQueryParams) (items []Order, isLastPage bool, err error)                                                                // 订单列表
	Order(id string) (item Order, exists bool, err error)                                                                                       // 单个订单
	CancelOrder(req CancelOrderRequest) (results []OrderCancelResult, err error)                                                                // 作废订单
	Products(params ProductsQueryParams) (items []Product, isLastPage bool, err error)                                                          // 商品列表
	Product(typ string, sku string, isAlias bool) (item Product, exists bool, err error)                                                        // 单个商品
	ProductExists(typ string, sku string, isAlias bool) (exists bool, err error)                                                                // 商品是否存在
	CreateProduct(req CreateProductRequest) error                                                                                               // 创建商品
	UpdateProduct(req UpdateProductRequest) error                                                                                               // 更新商品
	Packages(params PackagesQueryParams) (items []Package, isLastPage bool, err error)                                                          // 包裹列表
	Package(orderNumber, packageNumber string) (item Package, exists bool, err error)                                                           // 单个包裹
	PackageDeliver(req PackageDeliverRequest) error                                                                                             // 执行包裹发货
	Suppliers(params SuppliersQueryParams) (items []Supplier, isLastPage bool, err error)                                                       // 供应商列表
	PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error)                                        // 采购单列表
	CreatePurchaseOrder(params CreatePurchaseOrderRequest) (number string, err error)                                                           // 创建采购单
	PurchaseOrderArrival(req PurchaseOrderArrivalRequest) error                                                                                 // 采购单到货
	PurchaseOrderStockIn(req PurchaseOrderStockInRequest) error                                                                                 // 采购单入库
	PurchaseOrderStockInLogs(params PurchaseOrderLogsQueryParams) (items []PurchaseOrderLog, isLastPage bool, err error)                        // 采购单入库查询
	SaleAccounts(params SaleAccountsQueryParams) (items []SaleAccount, isLastPage bool, err error)                                              // 商户账号列表
	Stocks(params StocksQueryParams) (items []Stock, isLastPage bool, err error)                                                                // 库存列表
	StockChangeLogs(params StockChangeLogsQueryParams) (items []StockChangeLog, isLastPage bool, err error)                                     //  库存变动查询
	Warehouses(params WarehousesQueryParams) (items []Warehouse, isLastPage bool, err error)                                                    // 仓库列表
	Warehouse(id string) (item Warehouse, err error)                                                                                            // 仓库列表
	ShippingMethods(params ShippingMethodsQueryParams) (items []ShippingMethod, isLastPage bool, err error)                                     // 仓库物流渠道列表
	TrackingNumbers(params TrackingNumbersQueryParams) (items []TrackingNumber, isLastPage bool, err error)                                     // 订单物流单号列表
	Platforms() (items []Platform, err error)                                                                                                   // 平台及站点信息
	PurchaseSuggestionTemplates(params PurchaseSuggestionTemplatesQueryParams) (items []PurchaseSuggestionTemplate, isLastPage bool, err error) // 采购建议模板列表
	PurchaseSuggestions(params PurchaseSuggestionsQueryParams) (items []PurchaseSuggestion, isLastPage bool, err error)                         // 采购建议列表
	QuotePrices(params QuotedPricesQueryParams) (items []QuotedPrice, isLastPage bool, err error)                                               // 供应商报价查询
	AfterSales(params AfterSalesQueryParams) (items []AfterSale, isLastPage bool, err error)                                                    // 售后单信息查询
	PaypalTransaction(params PaypalTransactionsQueryParams) (items []PaypalTransaction, isLastPage bool, err error)                                         // Paypal 付款记录查询
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
