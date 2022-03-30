package logistics

type Paging struct {
	NextToken string `json:"nextToken,omitempty"` // 下一页标志,如果查询接口返回了则需要将该参数传入再次查询
	PageSize  int    `json:"limit,omitempty"`     // 分页返回的一页最大限制数量，可以传入 1 到 300 的值，默认为 50
}

func (p *Paging) SetPagingVars(nextToken string, pageSize, maxPageSize int) *Paging {
	p.NextToken = nextToken
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	p.PageSize = pageSize
	return p
}
