package logistics

// 获取包裹信息
// https://open.tongtool.com/apiDoc.html#/?docId=68251fa7414b43b5a1b0eddc898f6319

type Package struct {
}

type PackagesQueryParams struct {
	OrderStatus        string `json:"orderStatus"`        // 通途订单状态
	ShippingMethodCode string `json:"shippingMethodCode"` // 渠道代码
	Since              string `json:"since"`              // 查询的起始时间
}
