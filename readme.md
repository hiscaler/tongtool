通途 ERP 开放平台 API 封装
=======================

对 TongTool API 接口的封装，方便开发者调用，使用者无需关注接口认证、接口限制等繁琐的细节，提供 appKey 和 appSecret 即可使用。

接口返回具体格式和数据请参考 [通途接口文档](https://open.tongtool.com/apiDoc.html#/?docId=43a41f3680e04756a122d8671f2fc0ca)

针对通途返回的数据未做任何改动，所以具体格式您可以以通途开发文档为准，有部分接口为了方便开发者添加了一些扩展属性。

## 支持的方法

### ERP2

- AmazonAccountSites(params AmazonAccountSiteQueryParams) (items []string, isLastPage bool, err error)                                       // 查询亚马逊账号对应的站点
- FBAOrders(params FBAOrderQueryParams) (items []FBAOrder, isLastPage bool, err error)                                                       // FBA 订单列表
- ShopifyOrders(params ShopifyOrderQueryParams) (items []ShopifyOrder, isLastPage bool, err error)                                           // Shopify 订单列表
- CreateOrder(req CreateOrderRequest) (orderId, orderNumber string, err error)                                                               // 手工创建订单
- UpdateOrder(req UpdateOrderRequest) error                                                                                                  // 更新订单
- Orders(params OrderQueryParams) (items []Order, isLastPage bool, err error)                                                                // 订单列表
- Order(id string) (item Order, err error)                                                                                                   // 单个订单
- CancelOrder(req CancelOrderRequest) (results []OrderCancelResult, err error)                                                               // 作废订单
- Products(params ProductQueryParams) (items []Product, isLastPage bool, err error)                                                          // 商品列表
- Product(typ string, sku string, isAlias bool) (item Product, err error)                                                                    // 单个商品
- ProductExists(typ string, sku string, isAlias bool) (exists bool, err error)                                                               // 商品是否存在
- CreateProduct(req CreateProductRequest) error                                                                                              // 创建商品
- UpdateProduct(req UpdateProductRequest) error                                                                                              // 更新商品
- Packages(params PackageQueryParams) (items []Package, isLastPage bool, err error)                                                          // 包裹列表
- Package(orderNumber, packageNumber string) (item Package, err error)                                                                       // 单个包裹
- PackageDeliver(req PackageDeliverRequest) error                                                                                            // 执行包裹发货
- Suppliers(params SuppliersQueryParams) (items []Supplier, isLastPage bool, err error)                                                      // 供应商列表
- PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error)                                       // 采购单列表
- CreatePurchaseOrder(params CreatePurchaseOrderRequest) (number string, err error)                                                          // 创建采购单
- PurchaseOrderArrival(req PurchaseOrderArrivalRequest) error                                                                                // 采购单到货
- PurchaseOrderStockIn(req PurchaseOrderStockInRequest) error                                                                                // 采购单入库
- PurchaseOrderStockInLogs(params PurchaseOrderLogQueryParams) (items []PurchaseOrderLog, isLastPage bool, err error)                        // 采购单入库查询
- SaleAccounts(params SaleAccountQueryParams) (items []SaleAccount, isLastPage bool, err error)                                              // 商户账号列表
- Stocks(params StockQueryParams) (items []Stock, isLastPage bool, err error)                                                                // 库存列表
- StockChangeLogs(params StockChangeLogQueryParams) (items []StockChangeLog, isLastPage bool, err error)                                     //  库存变动查询
- Warehouses(params WarehouseQueryParams) (items []Warehouse, isLastPage bool, err error)                                                    // 仓库列表
- Warehouse(id string) (item Warehouse, err error)                                                                                           // 仓库列表
- ShippingMethods(params ShippingMethodQueryParams) (items []ShippingMethod, isLastPage bool, err error)                                     // 仓库物流渠道列表
- TrackingNumbers(params TrackingNumberQueryParams) (items []TrackingNumber, isLastPage bool, err error)                                     // 订单物流单号列表
- Platforms() (items []Platform, err error)                                                                                                  // 平台及站点信息
- PurchaseSuggestionTemplates(params PurchaseSuggestionTemplateQueryParams) (items []PurchaseSuggestionTemplate, isLastPage bool, err error) // 采购建议模板列表
- PurchaseSuggestions(params PurchaseSuggestionQueryParams) (items []PurchaseSuggestion, isLastPage bool, err error)                         // 采购建议列表
- QuotePrices(params QuotedPriceQueryParams) (items []QuotedPrice, isLastPage bool, err error)                                               // 供应商报价查询
- AfterSales(params AfterSaleQueryParams) (items []AfterSale, isLastPage bool, err error)                                                    // 售后单信息查询

## 配置
创建连接实例时，您需要提供一个配置参数。具体说明如下：
- Debug 是否为调试模式，开启的情况下会输入接口请求数据，在开发模式下建议开启，方便调试，生产系统上则建议关闭。
- AppKey 通途 APP Key
- AppSecret 通途 APP Secret
- EnableCache 是否激活缓存，激活的情况下，10 分钟内多次发起的请求第二次起都会从缓存中获取后直接返回，不会再走请求接口，如果您的数据变化比较频繁，建议关闭，以免获取不到最新的数据。同时也需要注意的是通途有接口请求次数限制，一分钟内最多发起 5 次接口请求，所以在应用端需要做相应的处理。

## 使用方法

```go
import (
    "github.com/hiscaler/tongtool"
    ttConfig "github.com/hiscaler/tongtool/config"
    "github.com/hiscaler/tongtool/erp2"
)

ttInstance := tongtool.NewTongTool(ttConfig.Config{
    Debug:       false,
    AppKey:      "",
    AppSecret:   "",
    EnableCache: true,
})
ttService := erp2.NewService(ttInstance)
params := OrderQueryParams{
    SaleDateFrom: "2021-12-01 00:00:00",
    SaleDateTo:   "2021-12-31 23:59:59",
}
orders := make([]Order, 0)
for {
    pageOrders, isLastPage, err := ttService.Orders(params)
    if err != nil {
        t.Errorf("ttService.Orders error: %s", err.Error())
    } else {
        orders = append(orders, pageOrders...)
    }
    if isLastPage || err != nil {
        break
    }
    params.PageNo++
}
fmt.Println(fmt.Sprintf("%#v", orders))
```

## 通途接口调用速率限制处理
所有接口调用频率为一分钟 5 次，需要调用端做好频率控制。但是通途接口并没有在返回数据中告知剩余的可访问次数，所以不能做到精细控制。

目前如果遇到调用速率限制，会间隔 5 秒后再次发起请求，最多发起两次。同时也建议在生产环境中开启缓存，进一步地避免该问题。
