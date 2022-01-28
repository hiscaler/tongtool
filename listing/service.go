package listing

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	ProductCategories(params ProductCategoryQueryParams) (items []ProductCategory, err error)
	CreateProductCategory(req CreateProductCategoryRequest) error
	UpdateProductCategory(req UpdateProductCategoryRequest) error
	DeleteProductCategory(req DeleteProductCategoryRequest) error
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
