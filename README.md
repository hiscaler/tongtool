通途 ERP 开放平台 SDK
=======================

针对 TongTool API 接口的封装，方便开发者调用，使用者无需关注接口认证、接口限制等繁琐的细节，提供 appKey 和 appSecret 即可使用。

接口返回具体格式和数据请参考 [通途接口文档](https://open.tongtool.com/apiDoc.html)

针对通途 API 接口返回的数据和格式未做任何改动，所以具体数据和格式您可以以通途开发文档为准。

## 支持的方法

### ERP2.0

- AmazonAccountSites(params AmazonAccountSitesQueryParams) (items []string, isLastPage bool, err error)                                       // 查询亚马逊账号对应的站点
- FBAOrders(params FBAOrdersQueryParams) (items []FBAOrder, isLastPage bool, err error)                                                       // FBA 订单列表
- ShopifyOrders(params ShopifyOrdersQueryParams) (items []ShopifyOrder, isLastPage bool, err error)                                           // Shopify 订单列表
- CreateOrder(req CreateOrderRequest) (orderId, orderNumber string, err error)                                                                // 手工创建订单
- UpdateOrder(req UpdateOrderRequest) error                                                                                                   // 更新订单
- Orders(params OrdersQueryParams) (items []Order, isLastPage bool, err error)                                                                // 订单列表
- Order(id string) (item Order, exists bool, err error)                                                                                       // 单个订单
- CancelOrder(req CancelOrderRequest) (results []OrderCancelResult, err error)                                                                // 作废订单
- OrderPair(req OrderPairRequest) error                                                                                                       // 订单配对
- Products(params ProductsQueryParams) (items []Product, isLastPage bool, err error)                                                          // 商品列表
- Product(typ string, sku string, isAlias bool) (item Product, exists bool, err error)                                                        // 单个商品
- ProductExists(typ string, sku string, isAlias bool) (exists bool, err error)                                                                // 商品是否存在
- CreateProduct(req CreateProductRequest) error                                                                                               // 创建商品
- UpdateProduct(req UpdateProductRequest) error                                                                                               // 更新商品
- Packages(params PackagesQueryParams) (items []Package, isLastPage bool, err error)                                                          // 包裹列表
- Package(orderNumber, packageNumber string) (item Package, exists bool, err error)                                                           // 单个包裹
- PackageDeliver(req PackageDeliverRequest) error                                                                                             // 执行包裹发货
- Suppliers(params SuppliersQueryParams) (items []Supplier, isLastPage bool, err error)                                                       // 供应商列表
- PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error)                                        // 采购单列表
- CreatePurchaseOrder(params CreatePurchaseOrderRequest) (number string, err error)                                                           // 创建采购单
- PurchaseOrderArrival(req PurchaseOrderArrivalRequest) error                                                                                 // 采购单到货
- PurchaseOrderStockIn(req PurchaseOrderStockInRequest) error                                                                                 // 采购单入库
- PurchaseOrderStockInLogs(params PurchaseOrderLogsQueryParams) (items []PurchaseOrderLog, isLastPage bool, err error)                        // 采购单入库查询
- SaleAccounts(params SaleAccountsQueryParams) (items []SaleAccount, isLastPage bool, err error)                                              // 商户账号列表
- Stocks(params StocksQueryParams) (items []Stock, isLastPage bool, err error)                                                                // 库存列表
- StockChangeLogs(params StockChangeLogsQueryParams) (items []StockChangeLog, isLastPage bool, err error)                                     // 库存变动查询
- Warehouses(params WarehousesQueryParams) (items []Warehouse, isLastPage bool, err error)                                                    // 仓库列表
- Warehouse(id string) (item Warehouse, exists bool, err error)                                                                               // 仓库列表
- WarehouseShippingMethods(params ShippingMethodsQueryParams) (items []WarehouseShippingMethod, isLastPage bool, err error)                   // 仓库物流渠道列表
- TrackingNumbers(params TrackingNumbersQueryParams) (items []TrackingNumber, isLastPage bool, err error)                                     // 订单物流单号列表
- Platforms() (items []Platform, err error)                                                                                                   // 平台及站点信息
- PurchaseSuggestionTemplates(params PurchaseSuggestionTemplatesQueryParams) (items []PurchaseSuggestionTemplate, isLastPage bool, err error) // 采购建议模板列表
- PurchaseSuggestions(params PurchaseSuggestionsQueryParams) (items []PurchaseSuggestion, isLastPage bool, err error)                         // 采购建议列表
- QuotePrices(params QuotedPricesQueryParams) (items []QuotedPrice, isLastPage bool, err error)                                               // 供应商报价查询
- AfterSales(params AfterSalesQueryParams) (items []AfterSale, isLastPage bool, err error)                                                    // 售后单信息查询
- PaypalTransactions(params PaypalTransactionsQueryParams) (items []PaypalTransaction, isLastPage bool, err error)                            // Paypal 付款记录查询
- SurfaceSheets(params SurfaceSheetsQueryParams) (items []SurfaceSheet, err error)                                                            // 通途 ERP 面单

### 刊登

- Categories(params CategoriesQueryParams) (items []Category, err error)    // 类目列表
- CreateCategory(req CreateCategoryRequest) error                           // 添加类目
- UpdateCategory(req UpdateCategoryRequest) error                           // 更新类目
- DeleteCategory(req DeleteCategoryRequest) error                           // 删除类目
- Tags(req TagsQueryParams) (items []Tag, err error)                        // 标签列表
- CreateTag(req CreateTagRequest) error                                     // 添加标签
- UpdateTag(req UpdateTagRequest) error                                     // 更新标签
- DeleteTag(req DeleteTagRequest) error                                     // 删除标签
- Warehouses(params WarehousesQueryParams) (items []Warehouse, err error)   // 仓库列表
- UpsertStockProduct(req UpsertStockProductRequest) error                   // 保存库存产品资料
- UpsertSaleAccount(req UpsertSaleAccountRequest) error                     // 保存店铺信息
- UpsertUser(req UpsertUserRequest) error                                   // 保存用户信息
- SaveUserAccount(req UpsertUserAccountRequest) error                       // 保存用户店铺信息
- Products(params ProductsQueryParams) (items []Product, err error)         // 批量获取售卖详情
- Product(params ProductQueryParams) (item Product, exists bool, err error) // 获取售卖基本资料
- UpdateProduct(req UpdateProductRequest) error                             // 修改售卖资料

### ERP3.0

- Products(params ProductsQueryParams) (items []Product, nextToken string, isLastPage bool, err error)                // 商品列表
- UserTicket(ticket string) (u User, refreshTicket string, expire int, err error)                                     // 根据 ticket 获取员工信息
- Suppliers(params SuppliersQueryParams) (items []Supplier, nextToken string, isLastPage bool, err error)             // 供应商列表
- WarehouseAreas(params WarehouseAreasQueryParams) (items []WarehouseArea, err error)                                 // 仓库分区关系
- SaveThirdAccounts(req UpdateThirdAccountRequest) error                                                              // 保存第三方帐号信息
- StockInSheets(params StockInSheetsQueryParams) (items []StockInSheet, nextToken string, isLastPage bool, err error) // 入库单列表
- AddShippingPackage(req AddShippingPackageRequest) (packages []ShippingPackage, err error)                           // 出库单交运

### 物流
- Packages(params PackagesQueryParams) (items []Package, nextToken string, isLastPage bool, err error) // 获取包裹信息
- WriteBackPackageProcessingResult(req PackageWriteBackRequest) error                                  // 回写包裹处理结果
- WriteBackPackageDeliveryInformation(req PackageWriteBackRequest) error                               // 回写包裹发货信息

## 配置

创建连接实例时，您需要提供一个配置参数。参数具体说明如下：

- Debug

  是否为调试模式，开启的情况下会输出接口请求数据，在开发模式下建议开启，方便调试，生产系统上则建议关闭。
- AppKey

  通途 APP Key（从通途开放平台的应用管理中获取）
- RetryCount

  HTTP 请求失败的情况下重试次数
- RetryWaitTime

  重试等待时间
- RetryMaxWaitTime

  最大重试等待时间
- ForceWaiting

  强制等待，如果设置为 true 的话，会总是等待接口端返回数据。
- AppSecret

  通途 APP Secret（从通途开放平台的应用管理中获取）
- EnableCache

  是否激活缓存，激活的情况下，10 分钟内多次发起的请求第二次起都会从缓存中获取后直接返回，不会再走请求接口，如果您的数据变化比较频繁，建议关闭，以免获取不到最新的数据。同时也需要注意的是通途有接口请求次数限制，一分钟内最多发起 5 次接口请求，所以在应用端需要做相应的处理。同时支持开启 forceWaiting 选项，如果设置为 true 的话，会总是等待接口端返回数据，您可以根据自己的需求开启或者关闭，默认情况下该选项是关闭的。

## 使用方法

```go
package main

import (
  "github.com/hiscaler/tongtool"
  ttConfig "github.com/hiscaler/tongtool/config"
  "github.com/hiscaler/tongtool/erp2"
  "fmt"
)

func main() {
  ttInstance := tongtool.NewTongTool(ttConfig.Config{
    Debug:       true,
    Timeout: 120,
    RetryCount: 2,
    RetryWaitTime: 12,
    RetryMaxWaitTime: 60,
    ForceWaiting: true,
    AppKey:      "",
    AppSecret:   "",
    EnableCache: false,
  })
  ttService := erp2.NewService(ttInstance)
  params := erp2.OrdersQueryParams{
    SaleDateFrom: "2021-12-01 00:00:00",
    SaleDateTo:   "2021-12-31 23:59:59",
  }
  params.PageNo = 1
  orders := make([]erp2.Order, 0)
  for {
    pageOrders, isLastPage, err := ttService.Orders(params)
    if err != nil {
      fmt.Println(fmt.Sprintf("ttService.Orders error: %s", err.Error()))
    } else {
      orders = append(orders, pageOrders...)
    }
    if isLastPage || err != nil {
      break
    }
    params.PageNo++
  }
  fmt.Println(fmt.Sprintf("%#v", orders))
}
```

## 注意事项

### 判断数据是否存在

针对单项数据的查询，通常返回数据格式为：

```go
FuncName(req Request) (item DataType, exists bool, err error)
```

如果你需要判断返回的数据是否有效，在使用数据前请先判断 exists 返回值是否为 true。为 true 的情况下数据一定是存在且有效的，而如果为 false 则需要继续判断 err 返回值是否为 nil，为 nil 则表示数据确实不存在，而非 nil 则是在查询过程中出现问题导致不能正常获取到相应的数据，需要您根据业务需求来确定后续的代码逻辑。

### 数据扩展

通途的返回格式比较混乱，比如布尔值的返回有多种（Y, 1, null, ""），为了减少开发者负担，针对这种情况做了部分处理，增加的属性为原属性名称增加 Boolean 后缀，返回值类型为布尔值。

### 调用速率限制处理

所有接口调用频率为一分钟 5 次，需要调用端做好频率控制。但是通途接口并没有在返回数据中告知剩余的可访问次数，所以不能做到精细控制。

目前如果遇到调用速率限制，请调整配置参数中的 `RetryCount`、`RetryWaitTime`、`RetryMaxWaitTime`、`ForceWaiting` 参数。同时也建议在生产环境中开启缓存，进一步地避免该问题。