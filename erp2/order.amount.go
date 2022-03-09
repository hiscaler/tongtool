package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/constant"
	"github.com/shopspring/decimal"
)

// 订单金额计算
// 计算后的值统一为人民币，如果需要使用其他币种，需要自行转换

// OrderIncomeAmount 订单收入
type OrderIncomeAmount struct {
	Product     float64 `json:"product"`      // 商品金额
	ProductTax  float64 `json:"product_tax"`  // 商品税金
	Shipping    float64 `json:"shipping"`     // 实付运费
	ShippingTax float64 `json:"shipping_tax"` // 实付运费税金
}

// OrderExpenditureAmount 订单支出
type OrderExpenditureAmount struct {
	Product  float64 `json:"product"`  // 商品成本
	Channel  float64 `json:"channel"`  // 渠道收费
	VAT      float64 `json:"vat"`      // 增值税（订单总金额 * 汇率）
	Shipping float64 `json:"shipping"` // 运费
}

type orderAmountConfig struct {
	rates     map[string]float64 // 汇率转换
	precision int32              // 保留精度
}

type OrderAmount struct {
	config                 orderAmountConfig      // 设置
	Number                 string                 `json:"number"`                   // 订单号
	Currency               string                 `json:"currency"`                 // 货币
	TotalQuantity          int                    `json:"total_quantity"`           // 商品总数量
	IncomeAmount           OrderIncomeAmount      `json:"income_amount"`            // 收入
	ExpenditureAmount      OrderExpenditureAmount `json:"expenditure_amount"`       // 支出
	TotalIncomeAmount      float64                `json:"total_income_amount"`      // 总收入金额
	TotalExpenditureAmount float64                `json:"total_expenditure_amount"` // 总支出成本
}

// 货币金额转换
func currencyExchange(value float64, exchangeRate map[string]float64, currency string) decimal.Decimal {
	decimalValue := decimal.NewFromFloat(value)
	if rate, ok := exchangeRate[currency]; ok {
		decimalValue = decimalValue.Div(decimal.NewFromFloat(rate))
	}
	return decimalValue
}

// NewOrderAmount
//
// exchangeRates 以人民币为基准，比如美元兑人民币 1:6.3，对应的设置为：
// map[string]float64{"USD": 0.1587}
// 默认情况下，转换后的币种为人民币，如果需要转换为其他币种，请使用 ExchangeTo 函数获取
func NewOrderAmount(order Order, exchangeRates map[string]float64, precision int32) *OrderAmount {
	oa := &OrderAmount{
		config:        orderAmountConfig{rates: exchangeRates, precision: precision},
		Number:        order.OrderIdCode,
		Currency:      CNY,
		TotalQuantity: 0,
		IncomeAmount: OrderIncomeAmount{
			Product:     0,
			ProductTax:  0,
			Shipping:    0,
			ShippingTax: 0,
		},
		ExpenditureAmount: OrderExpenditureAmount{
			Channel:  0,
			VAT:      0,
			Product:  0,
			Shipping: order.ShippingFee,
		},
		TotalIncomeAmount:      0,
		TotalExpenditureAmount: 0,
	}
	totalIncomeAmount := decimal.NewFromFloat(0)
	incomeProduct := decimal.NewFromFloat(0)
	for _, detail := range order.OrderDetails {
		value := currencyExchange(detail.TransactionPrice, exchangeRates, order.OrderAmountCurrency).Mul(decimal.NewFromInt(int64(detail.Quantity)))
		incomeProduct.Add(value)
		totalIncomeAmount.Add(value)
		oa.TotalQuantity += detail.Quantity
	}
	oa.IncomeAmount.Product, _ = incomeProduct.Round(precision).Float64() // 商品收入
	if order.TaxIncome > 0 {
		incomeProductTax := currencyExchange(order.TaxIncome, exchangeRates, order.TaxCurrency)
		oa.IncomeAmount.ProductTax, _ = incomeProductTax.Round(precision).Float64()
		totalIncomeAmount.Sub(incomeProductTax) // 商品金额 - 税金
	}
	incomeShipping := currencyExchange(order.ShippingFeeIncome, exchangeRates, order.ShippingFeeIncomeCurrency)
	oa.IncomeAmount.Shipping, _ = incomeShipping.Round(precision).Float64()
	totalIncomeAmount = totalIncomeAmount.Add(incomeShipping)
	oa.TotalIncomeAmount, _ = totalIncomeAmount.Round(precision).Float64()
	oa.ExpenditureAmount.Channel, _ = totalIncomeAmount.
		Mul(decimal.NewFromFloat(0.15)).
		Round(precision).
		Float64()
	if order.StoreCountryCode() == constant.CountryCodeUnitedKingdom {
		// ((商品金额 + 客户支付的运费) / 1.2 * 0.2) 简化后为 ((商品金额 + 客户支付的运费) / 6)
		oa.ExpenditureAmount.VAT, _ = incomeProduct.Add(incomeShipping).
			Div(decimal.NewFromInt(6)).
			Round(precision).
			Float64()
	}

	// 商品成本
	expenditureProduct := decimal.NewFromFloat(0)
	for _, good := range order.GoodsInfo.TongToolGoodsInfoList {
		var costPrice float64
		if good.ProductCurrentCost > 0 {
			costPrice = good.ProductCurrentCost
		} else if good.GoodsAverageCost > 0 {
			costPrice = good.GoodsAverageCost
		} else {
			costPrice = good.GoodsCurrentCost
		}
		if costPrice > 0 && good.Quantity > 0 {
			expenditureProduct.Add(decimal.NewFromFloat(costPrice).Mul(decimal.NewFromInt(int64(good.Quantity))))
		}
	}
	oa.ExpenditureAmount.Product, _ = expenditureProduct.Round(precision).Float64()
	oa.TotalExpenditureAmount, _ = decimal.NewFromFloat(oa.ExpenditureAmount.Product).
		Add(decimal.NewFromFloat(oa.ExpenditureAmount.VAT)).
		Add(decimal.NewFromFloat(oa.ExpenditureAmount.Shipping)).
		Add(decimal.NewFromFloat(oa.ExpenditureAmount.Channel)).
		Round(precision).
		Float64()
	return oa
}

// ExchangeTo 兑换
func (oa OrderAmount) ExchangeTo(currency string) (newOA OrderAmount, err error) {
	if v, ok := oa.config.rates[currency]; ok {
		precision := oa.config.precision
		rate := decimal.NewFromFloat(1).Div(decimal.NewFromFloat(v))
		newOA = oa
		newOA.Currency = currency

		if newOA.IncomeAmount.Product > 0 {
			newOA.IncomeAmount.Product, _ = decimal.NewFromFloat(newOA.IncomeAmount.Product).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeAmount.ProductTax > 0 {
			newOA.IncomeAmount.ProductTax, _ = decimal.NewFromFloat(newOA.IncomeAmount.ProductTax).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeAmount.Shipping > 0 {
			newOA.IncomeAmount.Shipping, _ = decimal.NewFromFloat(newOA.IncomeAmount.Shipping).Div(rate).Round(precision).Float64()
		}
		if newOA.IncomeAmount.ShippingTax > 0 {
			newOA.IncomeAmount.ShippingTax, _ = decimal.NewFromFloat(newOA.IncomeAmount.ShippingTax).Div(rate).Round(precision).Float64()
		}
		if newOA.ExpenditureAmount.Product > 0 {
			newOA.ExpenditureAmount.Product, _ = decimal.NewFromFloat(newOA.ExpenditureAmount.Product).Div(rate).Round(precision).Float64()
		}
		if newOA.ExpenditureAmount.Channel > 0 {
			newOA.ExpenditureAmount.Channel, _ = decimal.NewFromFloat(newOA.ExpenditureAmount.Channel).Div(rate).Round(precision).Float64()
		}
		if newOA.ExpenditureAmount.VAT > 0 {
			newOA.ExpenditureAmount.VAT, _ = decimal.NewFromFloat(newOA.ExpenditureAmount.VAT).Div(rate).Round(precision).Float64()
		}
		if newOA.ExpenditureAmount.Shipping > 0 {
			newOA.ExpenditureAmount.Shipping, _ = decimal.NewFromFloat(newOA.ExpenditureAmount.Shipping).Div(rate).Round(precision).Float64()
		}
		if newOA.TotalIncomeAmount > 0 {
			newOA.TotalIncomeAmount, _ = decimal.NewFromFloat(newOA.TotalIncomeAmount).Div(rate).Round(precision).Float64()
		}
		if newOA.TotalExpenditureAmount > 0 {
			newOA.TotalExpenditureAmount, _ = decimal.NewFromFloat(newOA.TotalExpenditureAmount).Div(rate).Round(precision).Float64()
		}
	} else {
		err = fmt.Errorf("无效的币种：%s", currency)
	}
	return
}
