package erp2

import (
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_TrackingNumbers(t *testing.T) {
	_, ttService := newTestTongTool()
	// 验证返回结果数量
	params := TrackingNumberQueryParams{}
	params.OrderIds = []string{"bad.order.id1", "bad.order.id2", "L-M20211221152430918"}
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
		t.Errorf("except return %d records, actual return %d records", len(params.OrderIds), len(trackingNumbers))
	} else {
		t.Log(cast.ToJson(trackingNumbers))
	}

	// 验证未提供订单集合参数
	_, _, err := ttService.TrackingNumbers(TrackingNumberQueryParams{})
	if err != nil {
		t.Errorf("ttService.TrackingNumbers error: %s", err.Error())
	}
}
