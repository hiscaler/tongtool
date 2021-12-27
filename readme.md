TongTool API 封装
=================

对 TongTool API 接口的封装，方便统一调用

## 支持的方法
### ERP2
- FBAOrders(params FBAOrderQueryParams) (items []FBAOrder, isLastPage bool, err error)                   // FBA 订单列表
- ShopifyOrders(params ShopifyOrderQueryParams) (items []ShopifyOrder, isLastPage bool, err error)       // Shopify 订单列表
- Orders(params OrderQueryParams) (items []Order, isLastPage bool, err error)                            // 订单列表
- Order(id string) (item Order, err error)                                                               // 单个订单
- Products(params ProductQueryParams) (items []Product, isLastPage bool, err error)                      // 商品列表
- Product(typ string, sku string, isAlias bool) (item Product, err error)                                // 单个商品
- ProductExists(typ string, sku string, isAlias bool) bool                                               // 商品是否存在
- CreateProduct(req CreateProductRequest) error                                                          // 创建商品
- UpdateProduct(req UpdateProductRequest) error                                                          // 更新商品
- Packages(params PackageQueryParams) (items []Package, isLastPage bool, err error)                      // 包裹列表
- Package(orderNumber, packageNumber string) (item Package, err error)                                   // 单个包裹
- Suppliers(params SuppliersQueryParams) (items []Supplier, isLastPage bool, err error)                  // 供应商列表
- PurchaseOrders(params PurchaseOrdersQueryParams) (items []PurchaseOrder, isLastPage bool, err error)   // 采购单列表
- CreatePurchaseOrder(params CreatePurchaseOrderRequest) (number string, err error)                      // 创建采购单
- SaleAccounts(params SaleAccountQueryParams) (items []SaleAccount, isLastPage bool, err error)          // 商户账号列表
- Stocks(params StockQueryParams) (items []Stock, isLastPage bool, err error)                            // 库存列表
- Warehouses(params WarehouseQueryParams) (items []Warehouse, isLastPage bool, err error)                // 仓库列表
- Warehouse(params WarehouseQueryParams) (item Warehouse, err error)                                     // 仓库列表
- ShippingMethods(params ShippingMethodQueryParams) (items []ShippingMethod, isLastPage bool, err error) // 仓库物流渠道列表
- TrackingNumbers(params TrackingNumberQueryParams) (items []TrackingNumber, isLastPage bool, err error) // 订单物流单号列表

## 使用方法

```go
import "github.com/hiscaler/tongtool"

ttInstance := tongtool.NewTongTool(AppKey, AppSecret, true)
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