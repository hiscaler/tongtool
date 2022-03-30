package logistics

import (
	"github.com/hiscaler/tongtool"
)

type service struct {
	tongTool *tongtool.TongTool
}

type Service interface {
	Packages(params PackagesQueryParams) (items []Package, nextToken string, isLastPage bool, err error) // 获取包裹信息
	WriteBackPackageProcessingResult(req PackageWriteBackRequest) error                                  // 回写包裹处理结果
	WriteBackPackageDeliveryInformation(req PackageWriteBackRequest) error                               // 回写包裹发货信息
}

func NewService(tt *tongtool.TongTool) Service {
	tt.QueryDefaultValues.PageSize = 300
	return service{tt}
}
