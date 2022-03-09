package erp2

import (
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
	Channel  float64 `json:"channel"`  // 渠道收费
	VAT      float64 `json:"vat"`      // VAT（订单总金额 * VAT 汇率）
	Product  float64 `json:"product"`  // 商品成本
	Shipping float64 `json:"shipping"` // 运费
}

type OrderAmount struct {
	Number                 string                 `json:"number"`                   // 订单号
	ExchangeRates          map[string]float64     `json:"exchange_rates"`           // 汇率转换
	TotalQuantity          int                    `json:"total_quantity"`           // 商品总数量
	IncomeAmount           OrderIncomeAmount      `json:"income_amount"`            // 收入
	ExpenditureAmount      OrderExpenditureAmount `json:"expenditure_amount"`       // 支出
	TotalIncomeAmount      float64                `json:"total_income_amount"`      // 总收入金额
	TotalExpenditureAmount float64                `json:"total_expenditure_amount"` // 总支出成本
}

// 兑换
func currencyExchange(value float64, exchangeRate map[string]float64, currency string) decimal.Decimal {
	decimalValue := decimal.NewFromFloat(value)
	if rate, ok := exchangeRate[currency]; ok {
		decimalValue = decimalValue.Mul(decimal.NewFromFloat(rate))
	}
	return decimalValue
}

func NewOrderAmount(order Order, exchangeRates map[string]float64, precision int32) *OrderAmount {
	oa := &OrderAmount{
		Number:                 order.OrderIdCode,
		ExchangeRates:          exchangeRates,
		TotalQuantity:          0,
		IncomeAmount:           OrderIncomeAmount{},
		ExpenditureAmount:      OrderExpenditureAmount{},
		TotalIncomeAmount:      0,
		TotalExpenditureAmount: 0,
	}
	tmp := decimal.NewFromFloat(0)
	totalIncomeAmount := decimal.NewFromFloat(0)
	incomeProduct := decimal.NewFromFloat(0)
	for _, detail := range order.OrderDetails {
		value := currencyExchange(detail.TransactionPrice, exchangeRates, order.OrderAmountCurrency).Mul(decimal.NewFromInt(int64(detail.Quantity)))
		incomeProduct.Add(value)
		totalIncomeAmount.Add(value)
		oa.TotalQuantity += detail.Quantity
	}
	oa.IncomeAmount.Product, _ = incomeProduct.Round(precision).Float64() // 商品收入
	tmp = currencyExchange(order.ShippingFeeIncome, exchangeRates, order.ShippingFeeIncomeCurrency)
	oa.IncomeAmount.Shipping, _ = tmp.Round(precision).Float64()
	totalIncomeAmount = totalIncomeAmount.Add(tmp)
	oa.TotalIncomeAmount, _ = totalIncomeAmount.Round(precision).Float64()
	oa.ExpenditureAmount.Shipping = order.ShippingFee
	oa.ExpenditureAmount.Channel, _ = incomeProduct.Add(decimal.NewFromFloat(oa.ExpenditureAmount.Shipping)).
		Mul(decimal.NewFromFloat(0.15)).
		Round(precision).
		Float64()
	if order.StoreCountryCode() == constant.CountryCodeUnitedKingdom {
		// ((商品金额 + 客户支付的运费) / 1.2 * 0.2) 简化后为 ((商品金额 + 客户支付的运费) / 6)
		oa.ExpenditureAmount.VAT, _ = incomeProduct.Add(decimal.NewFromFloat(oa.ExpenditureAmount.Shipping)).
			Div(decimal.NewFromInt(6)).
			Round(precision).
			Float64()
	}

	// 商品成本
	productAmount := decimal.NewFromFloat(0)
	for _, good := range order.GoodsInfo.TongToolGoodsInfoList {
		var costPrice float64
		if good.ProductCurrentCost > 0 {
			costPrice = good.ProductCurrentCost
		} else if good.GoodsAverageCost > 0 {
			costPrice = good.GoodsAverageCost
		} else {
			costPrice = good.GoodsCurrentCost
		}
		if costPrice > 0 {
			productAmount.Add(decimal.NewFromFloat(costPrice).Mul(decimal.NewFromInt(int64(good.Quantity))))
		}
	}
	oa.ExpenditureAmount.Product, _ = productAmount.Round(precision).Float64()
	oa.TotalExpenditureAmount, _ = decimal.NewFromFloat(oa.ExpenditureAmount.Product).
		Add(decimal.NewFromFloat(oa.ExpenditureAmount.VAT)).
		Add(decimal.NewFromFloat(oa.ExpenditureAmount.Shipping)).
		Add(decimal.NewFromFloat(oa.ExpenditureAmount.Channel)).
		Round(precision).
		Float64()

	return oa
}
