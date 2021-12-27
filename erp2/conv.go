package erp2

// PurchaseOrderStatusNtoS 采购单状态数字转字符
func PurchaseOrderStatusNtoS(n string) (s string) {
	switch n {
	case PurchaseOrderStatusDelivering:
		s = "delivering"
	case PurchaseOrderStatuspReceivedAndWaitM:
		s = "pReceivedAndWaitM"
	case PurchaseOrderStatusPartialReceived:
		s = "partialReceived"
	case PurchaseOrderStatusReceived:
		s = "Received"
	case PurchaseOrderStatusCancel:
		s = "cancel"
	}
	return
}
