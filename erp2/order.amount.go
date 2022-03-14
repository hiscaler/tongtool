package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/constant"
	"github.com/shopspring/decimal"
)

// 订单金额计算
// 计算后的值统一为人民币，如果需要使用其他币种，需要自行转换

var decimal1, decimal100 decimal.Decimal

func init() {
	decimal1 = decimal.NewFromInt(1)
	decimal100 = decimal.NewFromInt(100)
}

// orderIncome 订单收入
type orderIncome struct {
	Product   float64 `json:"product"`   // 商品金额
	Shipping  float64 `json:"shipping"`  // 运费（买家支付）
	Insurance float64 `json:"insurance"` // 保费
}

// orderExpenditure 订单支出
type orderExpenditure struct {
	Product  float64 `json:"product"`   // 商品成本
	Platform float64 `json:"platform "` // 平台佣金
	VAT      float64 `json:"vat"`       // 增值税（只有欧洲才有）
	Package  float64 `json:"package"`   // 包装
	Shipping float64 `json:"shipping"`  // 运费（卖家支付）
	Other    float64 `json:"other"`     // 其他（比如后端费用等）
}

type orderAmountConfig struct {
	rates     map[string]float64 // 汇率转换
	precision int32              // 保留精度
}

type orderItemExpenditure struct {
	Platform float64 `json:"platform "` // 平台佣金
	VAT      float64 `json:"vat"`       // 增值税（只有欧洲才有）
	Package  float64 `json:"package"`   // 包装
	Shipping float64 `json:"shipping"`  // 运费（卖家支付）
	Other    float64 `json:"other"`     // 其他（比如后端费用等）
}

// 订单详情项
type orderItem struct {
	StoreSKU    string               `json:"store_sku"`   // 平台 SKU
	SKU         string               `json:"sku"`         // 系统 SKU
	Price       float64              `json:"price"`       // 单价
	Quantity    int                  `json:"quantity"`    // 数量
	Amount      float64              `json:"amount"`      // 合计
	Expenditure orderItemExpenditure `json:"expenditure"` // 支出
}

// 订单收支
type orderIncomeExpenditure struct {
	Income      orderIncome      `json:"income"`      // 收入
	Expenditure orderExpenditure `json:"expenditure"` // 支出
}

// 汇总
type orderSummary struct {
	Income      float64 `json:"income"`      // 收入
	Expenditure float64 `json:"expenditure"` // 支出
	Profit      float64 `json:"profit"`      // 利润
}

// 占比（百分比）
type proportion struct {
	Product  float64 `json:"product"`  // 商品占比
	Shipping float64 `json:"shipping"` // 运费占比
	Platform float64 `json:"platform"` // 平台佣金占比
	Profit   float64 `json:"profit"`   // 利润占比
	Other    float64 `json:"other"`    // 其他占比
}

type OrderAmount struct {
	config            orderAmountConfig      // 设置
	Number            string                 `json:"number"`             // 订单号
	Currency          string                 `json:"currency"`           // 货币
	TotalQuantity     int                    `json:"total_quantity"`     // 商品总数量
	Items             []orderItem            `json:"items"`              // 详情
	IncomeExpenditure orderIncomeExpenditure `json:"income_expenditure"` // 收入支出
	Summary           orderSummary           `json:"summary"`            // 汇总
	Proportion        proportion             `json:"proportion"`         // 占比
}

// 货币金额转换
func currencyExchange(value float64, exchangeRate map[string]float64, currency string) decimal.Decimal {
	decimalValue := decimal.NewFromFloat(value)
	if rate, ok := exchangeRate[currency]; ok {
		decimalValue = decimalValue.Div(decimal1.Div(decimal.NewFromFloat(rate)))
	}
	return decimalValue
}

// NewOrderAmount
//
// exchangeRates 以人民币为基准，比如美元兑人民币 1:6.3，对应的设置为：
// map[string]float64{"USD": 6.3}
// 默认情况下，转换后的币种为人民币，如果需要转换为其他币种，请使用 ExchangeTo 函数获取
// 在传参过程中请注意，shippingFee, otherFee 均为人民币
func NewOrderAmount(order Order, exchangeRates map[string]float64, precision int32, shippingFee, otherFee float64) *OrderAmount {
	oa := &OrderAmount{
		config:        orderAmountConfig{rates: exchangeRates, precision: precision},
		Number:        order.OrderIdCode,
		Currency:      CNY,
		TotalQuantity: 0,
		IncomeExpenditure: orderIncomeExpenditure{
			Income: orderIncome{
				Product:   0,
				Shipping:  0,
				Insurance: 0,
			},
			Expenditure: orderExpenditure{
				Product:  0,
				Platform: 0,
				VAT:      0,
				Package:  0,
				Shipping: shippingFee,
				Other:    otherFee,
			},
		},
		Summary: orderSummary{
			Income:      0,
			Expenditure: 0,
		},
	}
	items := make([]orderItem, len(order.OrderDetails))
	// 收入
	incomeProduct := decimal.NewFromFloat(0)
	for i, detail := range order.OrderDetails {
		items[i] = orderItem{
			StoreSKU: detail.WebStoreSKU,
			SKU:      detail.GoodsMatchedSKU,
			Quantity: detail.Quantity,
		}
		quantity := decimal.NewFromInt(int64(detail.Quantity))
		price := currencyExchange(detail.TransactionPrice, exchangeRates, order.OrderAmountCurrency)
		amount := price.Mul(quantity)
		items[i].Price, _ = price.Round(precision).Float64()
		items[i].Amount, _ = amount.Round(precision).Float64()
		incomeProduct = incomeProduct.Add(amount)
		oa.TotalQuantity += detail.Quantity
	}
	oa.IncomeExpenditure.Income.Product, _ = incomeProduct.Round(precision).Float64() // 商品收入
	incomeShipping := currencyExchange(order.ShippingFeeIncome, exchangeRates, order.ShippingFeeIncomeCurrency)
	oa.IncomeExpenditure.Income.Shipping, _ = incomeShipping.Round(precision).Float64()
	incomeInsurance := currencyExchange(order.InsuranceIncome, exchangeRates, order.InsuranceIncomeCurrency)
	oa.IncomeExpenditure.Income.Insurance, _ = incomeInsurance.Round(precision).Float64()
	totalIncomeAmount := incomeProduct.Add(incomeShipping).Add(incomeInsurance) //
	oa.Summary.Income, _ = totalIncomeAmount.Round(precision).Float64()

	// 支出
	if order.PlatformCode == PlatformAmazon {
		oa.Proportion.Platform = 0.15 // 亚马逊固定 15%
		if order.StoreCountryCode() == constant.CountryCodeUnitedKingdom {
			// ((商品金额 + 客户支付的运费) / 1.2 * 0.2) 简化后为 ((商品金额 + 客户支付的运费) / 6)
			oa.IncomeExpenditure.Expenditure.VAT, _ = incomeProduct.Add(incomeShipping).
				Div(decimal.NewFromInt(6)).
				Round(precision).
				Float64()
		}
	} else {
		// todo
	}
	if oa.Proportion.Platform > 0 {
		v := decimal.NewFromFloat(oa.Proportion.Platform)
		oa.Proportion.Platform, _ = v.Mul(decimal100).Round(precision).Float64() // 转为百分比
		oa.IncomeExpenditure.Expenditure.Platform, _ = totalIncomeAmount.
			Mul(v).
			Round(precision).
			Float64() // 平台佣金
	}
	expenditureProduct := decimal.NewFromFloat(0) // 商品成本
	expenditurePacking := decimal.NewFromFloat(0) // 包装成本
	for _, good := range order.GoodsInfo.TongToolGoodsInfoList {
		var costPrice float64
		if good.GoodsAverageCost > 0 {
			costPrice = good.GoodsAverageCost // 货品平均成本
		} else if good.GoodsCurrentCost > 0 {
			costPrice = good.GoodsCurrentCost // 货品成本（最新成本）
		} else if good.ProductAverageCost > 0 {
			costPrice = good.ProductAverageCost // 商品平均成本
		} else {
			costPrice = good.ProductCurrentCost // 商品成本
		}
		if costPrice > 0 && good.Quantity > 0 {
			expenditureProduct = expenditureProduct.Add(decimal.NewFromFloat(costPrice).Mul(decimal.NewFromInt(int64(good.Quantity))))
		}
		// 包装成本
		if good.GoodsPackagingCost > 0 {
			expenditurePacking = expenditurePacking.Add(decimal.NewFromFloat(good.GoodsPackagingCost))
		}
	}
	if !expenditureProduct.IsZero() {
		oa.IncomeExpenditure.Expenditure.Product, _ = expenditureProduct.Round(precision).Float64()
	}
	if !expenditurePacking.IsZero() {
		oa.IncomeExpenditure.Expenditure.Package, _ = expenditurePacking.Round(precision).Float64()
	}
	summaryExpenditure := decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Product).
		Add(decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.VAT)).
		Add(decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Shipping)).
		Add(decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Platform)).
		Add(decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Package)).
		Add(decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Other))
	oa.Summary.Expenditure, _ = summaryExpenditure.Round(precision).Float64()
	summaryProfit := totalIncomeAmount.Sub(summaryExpenditure)
	oa.Summary.Profit, _ = summaryProfit.Round(precision).Float64()
	oa.Proportion.Product, _ = expenditureProduct.
		Div(totalIncomeAmount).
		Mul(decimal100).
		Round(precision).
		Float64()
	oa.Proportion.Shipping, _ = decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Shipping).
		Div(totalIncomeAmount).
		Mul(decimal100).
		Round(precision).
		Float64()
	oa.Proportion.Other, _ = decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Other).
		Div(totalIncomeAmount).
		Mul(decimal100).
		Round(precision).
		Float64()
	oa.Proportion.Profit, _ = summaryProfit.
		Div(totalIncomeAmount).
		Mul(decimal100).
		Round(precision).
		Float64()

	totalQuantity := decimal.NewFromInt(int64(oa.TotalQuantity))
	for i, item := range items {
		quantity := decimal.NewFromInt(int64(item.Quantity))
		item.Expenditure.Platform, _ = decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Platform).
			Div(totalQuantity).
			Mul(quantity).
			Round(precision).
			Float64()
		item.Expenditure.VAT, _ = decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.VAT).
			Div(totalQuantity).
			Mul(quantity).
			Round(precision).
			Float64()
		item.Expenditure.Package, _ = decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Package).
			Div(totalQuantity).
			Mul(quantity).
			Round(precision).
			Float64()
		item.Expenditure.Shipping, _ = decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Shipping).
			Div(totalQuantity).
			Mul(quantity).
			Round(precision).
			Float64()
		item.Expenditure.Other, _ = decimal.NewFromFloat(oa.IncomeExpenditure.Expenditure.Other).
			Div(totalQuantity).
			Mul(quantity).
			Round(precision).
			Float64()
		items[i] = item
	}
	oa.Items = items
	return oa
}

// ExchangeTo 兑换
func (oa OrderAmount) ExchangeTo(currency string) (newOA OrderAmount, err error) {
	if v, ok := oa.config.rates[currency]; ok {
		precision := oa.config.precision
		rate := decimal.NewFromFloat(v)
		newOA = oa
		newOA.Currency = currency
		for i, item := range newOA.Items {
			newOA.Items[i].Price, _ = decimal.NewFromFloat(item.Price).Div(rate).Round(precision).Float64()
			newOA.Items[i].Amount, _ = decimal.NewFromFloat(item.Amount).Div(rate).Round(precision).Float64()
			newOA.Items[i].Expenditure.Platform, _ = decimal.NewFromFloat(item.Expenditure.Platform).Div(rate).Round(precision).Float64()
			newOA.Items[i].Expenditure.VAT, _ = decimal.NewFromFloat(item.Expenditure.VAT).Div(rate).Round(precision).Float64()
			newOA.Items[i].Expenditure.Package, _ = decimal.NewFromFloat(item.Expenditure.Package).Div(rate).Round(precision).Float64()
			newOA.Items[i].Expenditure.Shipping, _ = decimal.NewFromFloat(item.Expenditure.Shipping).Div(rate).Round(precision).Float64()
			newOA.Items[i].Expenditure.Other, _ = decimal.NewFromFloat(item.Expenditure.Other).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeExpenditure.Income.Product > 0 {
			newOA.IncomeExpenditure.Income.Product, _ = decimal.NewFromFloat(newOA.IncomeExpenditure.Income.Product).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeExpenditure.Income.Shipping > 0 {
			newOA.IncomeExpenditure.Income.Shipping, _ = decimal.NewFromFloat(newOA.IncomeExpenditure.Income.Shipping).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeExpenditure.Expenditure.Product > 0 {
			newOA.IncomeExpenditure.Expenditure.Product, _ = decimal.NewFromFloat(newOA.IncomeExpenditure.Expenditure.Product).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeExpenditure.Expenditure.Platform > 0 {
			newOA.IncomeExpenditure.Expenditure.Platform, _ = decimal.NewFromFloat(newOA.IncomeExpenditure.Expenditure.Platform).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeExpenditure.Expenditure.VAT > 0 {
			newOA.IncomeExpenditure.Expenditure.VAT, _ = decimal.NewFromFloat(newOA.IncomeExpenditure.Expenditure.VAT).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeExpenditure.Expenditure.Shipping > 0 {
			newOA.IncomeExpenditure.Expenditure.Shipping, _ = decimal.NewFromFloat(newOA.IncomeExpenditure.Expenditure.Shipping).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeExpenditure.Expenditure.Other > 0 {
			newOA.IncomeExpenditure.Expenditure.Other, _ = decimal.NewFromFloat(newOA.IncomeExpenditure.Expenditure.Other).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeExpenditure.Expenditure.Package > 0 {
			newOA.IncomeExpenditure.Expenditure.Package, _ = decimal.NewFromFloat(newOA.IncomeExpenditure.Expenditure.Package).Div(rate).Round(precision).Float64()
		}
		if newOA.Summary.Income > 0 {
			newOA.Summary.Income, _ = decimal.NewFromFloat(newOA.Summary.Income).Div(rate).Round(precision).Float64()
		}
		if newOA.Summary.Expenditure > 0 {
			newOA.Summary.Expenditure, _ = decimal.NewFromFloat(newOA.Summary.Expenditure).Div(rate).Round(precision).Float64()
		}
		if newOA.Summary.Profit > 0 {
			newOA.Summary.Profit, _ = decimal.NewFromFloat(newOA.Summary.Profit).Div(rate).Round(precision).Float64()
		}
	} else {
		err = fmt.Errorf("无效的币种：%s", currency)
	}
	return
}

// ExchangeMoney 兑换指定的值，且以指定的货币形式返回
func (oa OrderAmount) ExchangeMoney(currency string, value float64) (money float64, err error) {
	if v, ok := oa.config.rates[currency]; ok {
		if value == 0 {
			return
		}
		money, _ = decimal.NewFromFloat(value).
			Div(decimal.NewFromFloat(v)).
			Round(oa.config.precision).
			Float64()
	} else {
		err = fmt.Errorf("无效的币种：%s", currency)
	}
	return
}
