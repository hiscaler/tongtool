package erp2

import (
	"testing"
)

func TestService_PaypalTransaction(t *testing.T) {
	params := PaypalTransactionsQueryParams{}
	params.PageNo = 1
	transactions := make([]PaypalTransaction, 0)
	for {
		pageOrders, isLastPage, err := ttService.PaypalTransactions(params)
		if err != nil {
			t.Errorf("ttService.PaypalTransactions error: %s", err.Error())
		} else {
			transactions = append(transactions, pageOrders...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
}
