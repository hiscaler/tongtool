package erp2

// PurchaseOrderStatusNtoS 采购单状态数字转字符
func PurchaseOrderStatusNtoS(n string) (s string) {
	switch n {
	case PurchaseOrderNumberStatusDelivering:
		s = PurchaseOrderStatusDelivering
	case PurchaseOrderNumberStatusPReceivedAndWaitM:
		s = PurchaseOrderStatusPReceivedAndWaitM
	case PurchaseOrderNumberStatusPartialReceived:
		s = PurchaseOrderStatusPartialReceived
	case PurchaseOrderNumberStatusReceived:
		s = PurchaseOrderStatusReceived
	case PurchaseOrderNumberStatusCancel:
		s = PurchaseOrderStatusCancel
	}
	return
}
