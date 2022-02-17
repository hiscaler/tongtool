package erp2

import (
	"testing"
)

func TestService_PaypalTransaction(t *testing.T) {
	params := PaypalTransactionsQueryParams{}
	transactions := make([]PaypalTransaction, 0)
	for {
		pageOrders, isLastPage, err := ttService.PaypalTransaction(params)
		if err != nil {
			t.Errorf("ttService.PaypalTransaction error: %s", err.Error())
		} else {
			transactions = append(transactions, pageOrders...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
}
