package erp2

import "testing"

func TestService_QuotePrices(t *testing.T) {
	params := QuotedPriceQueryParams{
		QuotedPriceDateBegin: "2018-01-01 00:00:00",
		QuotedPriceDateEnd:   "2018-01-02 00:00:00",
		SKU:                  "Lillian201309130003",
	}
	quotedPrices := make([]QuotedPrice, 0)
	for {
		pageQuotedPrices, isLastPage, err := ttService.QuotePrices(params)
		if err != nil {
			t.Errorf("ttService.QuotePrices error: %s", err.Error())
		} else {
			quotedPrices = append(quotedPrices, pageQuotedPrices...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
}
