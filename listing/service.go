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
	Tags(req TagQueryParams) (items []Tag, err error)                                         // 标签列表
	CreateTag(req CreateTagRequest) error                                                     // 添加标签
	UpdateTag(req UpdateTagRequest) error                                                     // 更新标签
	DeleteTag(req DeleteTagRequest) error                                                     // 删除标签
	Warehouses(req WarehouseQueryParams) (items []Warehouse, err error)                       // 仓库列表
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
