package erp2

import (
	"testing"
)

func TestService_TrackingNumbers(t *testing.T) {
	_, ttService := newTestTongTool()
	params := TrackingNumberQueryParams{}
	params.OrderIds = []string{"bad.order.id"}
	trackingNumbers := make([]TrackingNumber, 0)
	for {
		pageTrackingNumbers, isLastPage, err := ttService.TrackingNumbers(params)
		if err != nil {
			t.Errorf("ttService.TrackingNumbers error: %s", err.Error())
		} else {
			trackingNumbers = append(trackingNumbers, pageTrackingNumbers...)
		}
		if isLastPage || err != nil {
			break
		}
		params.PageNo++
	}
	if len(trackingNumbers) != len(params.OrderIds) {
		t.Errorf("order")
	}

	// 未提供订单集合参数
	_, _, err := ttService.TrackingNumbers(TrackingNumberQueryParams{})
	if err != nil {
		t.Errorf("ttService.TrackingNumbers error: %s", err.Error())
	}
}
