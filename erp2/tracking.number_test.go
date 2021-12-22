package erp2

import (
	"fmt"
	"testing"
)

func TestService_TrackingNumbers(t *testing.T) {
	_, ttService := newTestTongTool()
	params := TrackingNumberQueryParams{}
	params.OrderIds = []string{"L-M20211221152430918"}
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
	fmt.Println(fmt.Sprintf("%#v", trackingNumbers))
}
