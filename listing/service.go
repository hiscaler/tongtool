package listing

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	ProductCategories(params ProductCategoryQueryParams) (items []ProductCategory, err error) // 商品类目列表
	CreateProductCategory(req CreateProductCategoryRequest) error                             // 添加商品类目
	UpdateProductCategory(req UpdateProductCategoryRequest) error                             // 更新商品类目
	DeleteProductCategory(req DeleteProductCategoryRequest) error                             // 删除商品类目
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
