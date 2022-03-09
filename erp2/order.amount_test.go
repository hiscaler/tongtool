package erp2

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
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
	if orderAmount.TotalExpenditureAmount != 14.73 {
		t.Errorf("TotalExpenditureAmount excepted %f, actual %f", 10.0, orderAmount.TotalExpenditureAmount)
	}
	if orderAmount.IncomeAmount.Shipping != 31.5 {
		t.Errorf("IncomeAmount.Shipping excepted %f, actual %f", 10.0, orderAmount.IncomeAmount.Shipping)
	}
	fmt.Println(jsonx.ToJson(orderAmount, "{}"))
}
