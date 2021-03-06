package erp2

import (
	"errors"
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestOrderAmount(t *testing.T) {
	orderNumber := "O1234567"
	order, _, err := ttService.Order(orderNumber)
	if err != nil {
		if errors.Is(err, tongtool.ErrNotFound) {
			t.Errorf("%s not exists in tongtool", orderNumber)
		} else {
			t.Errorf("ttService.Order error: %s", err.Error())
		}
	} else if !strings.EqualFold(order.OrderIdCode, orderNumber) {
		t.Errorf("order.OrderIdKey %s not match %s", order.OrderIdCode, orderNumber)
	} else {
		orderAmount := NewOrderAmount(order, map[string]float64{
			constant.USD: 6.3927,
			constant.CNY: 1,
		}, 2, 22.45, 10)
		assert.Equal(t, 8.35, orderAmount.Summary.Expenditure, "order 1")
		assert.Equal(t, 98.72, orderAmount.Summary.Income, "order 2")
		newOrder, err := orderAmount.ExchangeTo(constant.USD)
		assert.Equal(t, nil, err, "newOrder 1")
		fmt.Println(jsonx.ToPrettyJson(orderAmount))
		fmt.Println(jsonx.ToPrettyJson(newOrder))
	}
}
