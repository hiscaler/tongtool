package erp2

import "github.com/hiscaler/tongtool/constant"

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

func exchangeAfter(value float64, exchangeRate map[string]float64, currency string) float64 {
	if rate, ok := exchangeRate[currency]; ok {
		return value * rate
	}
	return value
}

func NewOrderAmount(order Order, exchangeRates map[string]float64) *OrderAmount {
	oa := &OrderAmount{
		Number:                 order.OrderIdCode,
		ExchangeRates:          exchangeRates,
		TotalQuantity:          0,
		IncomeAmount:           OrderIncomeAmount{},
		ExpenditureAmount:      OrderExpenditureAmount{},
		TotalIncomeAmount:      0,
		TotalExpenditureAmount: 0,
	}
	var totalIncomeAmount float64
	for _, detail := range order.OrderDetails {
		value := exchangeAfter(detail.TransactionPrice, exchangeRates, order.OrderAmountCurrency) * float64(detail.Quantity)
		oa.IncomeAmount.Product += value
		totalIncomeAmount += value
		oa.TotalQuantity += detail.Quantity
	}
	oa.TotalIncomeAmount = totalIncomeAmount
	oa.IncomeAmount.Shipping = exchangeAfter(order.ShippingFeeIncome, exchangeRates, order.ShippingFeeIncomeCurrency)
	oa.ExpenditureAmount.Shipping = order.ShippingFee
	oa.ExpenditureAmount.Channel = (oa.IncomeAmount.Product + oa.IncomeAmount.Shipping) * 0.15
	if order.StoreCountryCode() == constant.CountryCodeUnitedKingdom {
		// ((商品金额 + 客户支付的运费) / 1.2 * 0.2) 简化后为 ((商品金额 + 客户支付的运费) / 6)
		oa.ExpenditureAmount.VAT = (oa.IncomeAmount.Product + oa.IncomeAmount.Shipping) / 6
	}
	// 商品成本
	for _, good := range order.GoodsInfo.TongToolGoodsInfoList {
		var costPrice float64
		if good.ProductCurrentCost > 0 {
			costPrice = good.ProductCurrentCost
		} else if good.GoodsAverageCost > 0 {
			costPrice = good.GoodsAverageCost
		} else {
			costPrice = good.GoodsCurrentCost
		}
		oa.ExpenditureAmount.Product = costPrice
	}
	oa.TotalExpenditureAmount = oa.ExpenditureAmount.Product +
		oa.ExpenditureAmount.VAT +
		oa.ExpenditureAmount.Shipping +
		oa.ExpenditureAmount.Channel
	return oa
}
