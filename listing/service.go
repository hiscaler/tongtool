package listing

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	Categories(params CategoryQueryParams) (items []Category, err error) // 类目列表
	CreateCategory(req CreateCategoryRequest) error                      // 添加类目
	UpdateCategory(req UpdateCategoryRequest) error                      // 更新类目
	DeleteCategory(req DeleteCategoryRequest) error                      // 删除类目
	Tags(req TagQueryParams) (items []Tag, err error)                    // 标签列表
	CreateTag(req CreateTagRequest) error                                // 添加标签
	UpdateTag(req UpdateTagRequest) error                                // 更新标签
	DeleteTag(req DeleteTagRequest) error                                // 删除标签
	Warehouses(req WarehouseQueryParams) (items []Warehouse, err error)  // 仓库列表
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
