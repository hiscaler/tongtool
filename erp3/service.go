package erp3

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	Products(params ProductsQueryParams) (items []Product, nextToken string, isLastPage bool, err error) // 商品列表
}

func NewService(tt *tongtool.TongTool) Service {
	tt.QueryDefaultValues.PageSize = 500
	return service{tt}
}
