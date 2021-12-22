package erp3

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	Products(params ProductQueryParams) (items []Product, nextToken string, isLastPage bool, err error) // 商品列表
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
