package erp2

import "errors"

type Supplier struct {
	AccountName         string  `json:"accountName"`
	Bank                string  `json:"bank"`
	BillingCycle        float64 `json:"billingCycle"`
	BillingCycleUnit    string  `json:"billingCycleUnit"`
	CityCnName          string  `json:"cityCnName"`
	ClearingForm        string  `json:"clearingForm"`
	ClearingRemark      string  `json:"clearingRemark"`
	CorporationFullName string  `json:"corporationFullname"`
	CountryCnName       string  `json:"countryCnName"`
	Description         string  `json:"description"`
	DetailAddress       string  `json:"detailAddress"`
	Email               string  `json:"email"`
	FaxNumber           string  `json:"faxNumber"`
	FullAddress         string  `json:"fullAddress"`
	IsDefault           string  `json:"isDefult"`
	Linkman             string  `json:"linkman"`
	PayeeAccount        string  `json:"payeeAccount"`
	PaymentMode         string  `json:"paymentMode"`
	PostalCode          string  `json:"postalCode"`
	QqNumber            string  `json:"qqNumber"`
	StateCnName         string  `json:"stateCnName"`
	SupplierCode        string  `json:"supplierCode"`
	SupplierGrade       string  `json:"supplierGrade"`
	SupplierId          string  `json:"supplierId"`
	Telephone           string  `json:"telephone"`
	WwNumber            string  `json:"wwNumber"`
	ZipCode             string  `json:"zipCode"`
}

type SuppliersQueryParams struct {
	MerchantId string `json:"merchantId"`
	PageNo     int    `json:"pageNo,omitempty"`
	PageSize   int    `json:"pageSize,omitempty"`
}

type supplierResult struct {
	result
	Datas struct {
		Array    []Supplier `json:"array"`
		PageNo   int        `json:"pageNo"`
		PageSize int        `json:"pageSize"`
	}
}

func (s service) Suppliers(params SuppliersQueryParams) (items []Supplier, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = s.tongTool.QueryDefaultValues.PageNo
	}
	if params.PageSize <= 0 {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	params.MerchantId = s.tongTool.MerchantId
	items = make([]Supplier, 0)
	res := supplierResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/supplierQuery")
	if err == nil {
		if resp.IsSuccess() {
			items = res.Datas.Array
			isLastPage = len(items) < params.PageSize
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}
