package pui

// PackageUploadItemDHL DHL
type PackageUploadItemDHL struct {
	ContentIndicator string  `json:"contentIndicator"` // 带电类型(00:不含电池,01:含锂电池,04:含蓄电池)
	CustomerPrefix   string  `json:"customerPrefix"`   // 物流商提供customerPrefix
	Incoterm         string  `json:"incoterm"`         // 交货方式(DDU或者DDP)
	InsuranceValue   float64 `json:"insuranceValue"`   // 保险金
	PickupAccountId  string  `json:"pickupAccountId"`  // 物流商提供pickupAccountId
	SoldToAccountId  string  `json:"soldToAccountId"`  // 物流商提供soldToAccountId
}

// PackageUploadItemDHLExpress DHL快递
type PackageUploadItemDHLExpress struct {
	// 贸易条款:CFR-(Cost and freight);CIF-(Cost, insurance, freight);CIP-(Carriage and insurance paid to);CPT-(Carriage paid to);DAF-(Delivered at frontier);DAP-(Delivered At Place);DAT-(Delivered At Terminal);DDP-(Delivered Duty Paid);DDU-(Delivered Duty unpaid);DEQ-(Delivered ex quay);DES-(Delivered ex ship);DVU-(Delivered Duty Paid VAT Unpaid);EXW-(Ex works);FAS-(Free alongside Ship);FCA-(Free Carrier);FOB-(Free on Board);DPU-(Delivered at Place Unloaded)
	DHLTermsOfTrade string `json:"dhlTermsOfTrade"`
	// 出口类型:P-Permanent;T-Temporary;R-Return For Repair;M-Used Exhibition Goods To Origin;I-Intercompany Use;C-Commercial Purpose Or Sale;E-Personal Belongings or Personal Use;S-Sample;G-Gift;U-Return To Origin;W-Warranty Replacement;D-Diplmatic Goods;F-Defenece Material 默认值：P
	ExportReasonType     string `json:"exportReasonType"`
	InvoiceType          string `json:"invoiceType"`          // 发票类型:CMI-商业发票;PFI-形式发票 默认值：CMI
	IsDutiable           string `json:"isDutiable"`           // 征税规则:Y:包裹;N：文件
	Password string `json:"password"`                         // 账号密码
	PlaceOfIncoterm      string `json:"placeOfIncoterm"`      // 贸易条款所适用的港口名称
	PltService           string `json:"pltService"`           // 是否开启PLT服务：PLTservice-PLT服务只适用于包裹类运单，不适用于文件类运单
	ShipperAccountNumber string `json:"shipperAccountNumber"` // 发件人账号
	SignatureImage       string `json:"signatureImage"`       // 签名图像：图片base64加密串
	SiteId               string `json:"siteId"`               // 账号ID
}
