package erp3

type Paging struct {
	NextToken string `json:"nextToken,omitempty"` // 下一页标志，如果查询接口返回了则需要将该参数传入再次查询
	PageSize  int    `json:"pageSize,omitempty"`  // 每页显示数量，范围 [1~500]
}

func (p *Paging) SetPagingVars(nextToken string, pageSize, maxPageSize int) *Paging {
	p.NextToken = nextToken
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	p.PageSize = pageSize
	return p
}
