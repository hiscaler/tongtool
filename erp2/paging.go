package erp2

type Paging struct {
	PageNo   int `json:"pageNo,omitempty"`   // 查询页数
	PageSize int `json:"pageSize,omitempty"` // 每页返回数量，默认值：100, 最大值：100，超过最大值以最大值数量返回
}

func (p *Paging) SetPagingVars(page, pageSize, maxPageSize int) *Paging {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	p.PageNo = page
	p.PageSize = pageSize
	return p
}
