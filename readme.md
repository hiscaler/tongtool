TongTool API 封装
=================

对 TongTool API 接口的封装，方便统一调用

## 使用方法

```go
import "github.com/hiscaler/tongtool"

ttInstance := tongtool.NewTongTool(AppKey, AppSecret, true)
ttService := NewService(ttInstance)
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