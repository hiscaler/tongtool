package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderAmount(t *testing.T) {
	order := Order{
		OrderIdCode:               "O123",
		ShippingFee:               10,
		ShippingFeeIncome:         5,
		ShippingFeeIncomeCurrency: USD,
	}
	orderAmount := NewOrderAmount(order, map[string]float64{
		USD: 6.3,
		CNY: 1,
	}, 2)
	assert.Equal(t, 14.73, orderAmount.TotalExpenditureAmount, "order 1")
	assert.Equal(t, 31.5, orderAmount.IncomeAmount.Shipping, "order 2")
	newOrder, err := orderAmount.ExchangeTo(USD)
	assert.Equal(t, nil, err, "newOrder 1")
	fmt.Println(jsonx.ToJson(orderAmount, "{}"))
	fmt.Println(jsonx.ToJson(newOrder, "{}"))
}
