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
	UpsertStockProduct(req UpsertStockProductRequest) error              // 保存库存产品资料
	UpsertSaleAccount(req UpsertSaleAccountRequest) error                // 保存店铺信息
	UpsertUser(req UpsertUserRequest) error                              // 保存用户信息
	SaveUserAccount(req UpsertUserAccountRequest) error                  // 保存用户店铺信息
	UpdateProduct(req UpdateProductRequest) error                        // 修改售卖资料
	DeleteProduct(req DeleteProductRequest) error                        // 删除售卖资料
	Products(req ProductQueryParams) (items []Product, err error)        // 批量获取售卖详情
}

func NewService(tt *tongtool.TongTool) Service {
	return service{tt}
}
