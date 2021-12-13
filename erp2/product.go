package erp2

import (
	"errors"
	"fmt"
	"github.com/hiscaler/tongtool"
	"strconv"
	"strings"
	"time"
)

// 销售类型
const (
	ProductTypeNormal   = iota // 0 普通销售
	ProductTypeVariable        // 1 变参销售
	ProductTypeBinding         // 2 捆绑销售
)

// Product 通途商品
type Product struct {
	ProductId          string  `json:"product_id"`
	ProductCode        string  `json:"productCode"`
	BrandName          string  `json:"brandName"`
	CategoryName       string  `json:"categoryName"`
	DeclareCnName      string  `json:"declareCnName"`
	DeclareEnName      string  `json:"declareEnName"`
	DeveloperName      string  `json:"developerName"`
	HsCode             string  `json:"hsCode"`
	InquirerName       string  `json:"inquirerName"`
	PackageCost        float64 `json:"packageCost"`
	PackageHeight      float64 `json:"packageHeight"`
	PackageLength      float64 `json:"packageLength"`
	SKU                string  `json:"sku"`
	ProductName        string  `json:"productName"`        // 产品名称
	ProductPackingName string  `json:"productPackingName"` // 中文配货名称
	ProductImgList     []struct {
		ImageGroupId string `json:"imageGroupId"`
	} `json:"productImgList"` // 产品图片
	LabelList []struct {
		SKULabel string `json:"skuLabel"`
	} `json:"labelList"`
	Status       string          `json:"status"`
	SupplierName string          `json:"supplier_name"`
	GoodsDetail  []ProductDetail `json:"goodsDetail"`
	CreatedDate  int             `json:"createdDate"`
	UpdatedDate  time.Time       `json:"updated_date"`
}

// ProductDetail 通途商品详情
type ProductDetail struct {
	GoodsSku      string  `json:"goodsSku"`
	GoodsWeight   float64 `json:"goodsWeight"`
	GoodsAveCost  float64 `json:"goodsAveCost"`
	GoodsCurCost  float64 `json:"goodsCurCost"`
	GoodsDetailId string  `json:"goodsDetailId"`
}

type Label struct {
	SKULabel string `json:"skuLabel"`
}

type ProductSupplier struct {
	SupplierName        string `json:"supplierName"`
	MinPurchaseQuantity int    `json:"minPurchaseQuantity"`
	PurchaseRemark      string `json:"purchaseRemark"`
}

type CreateProductRequest struct {
	ProductCode          string            `json:"productCode"` // 商品编号PCL
	ProductName          string            `json:"productName"` // 商品名
	ProductLabelIds      []string          `json:"productLabelIds"`
	ProductPackingEnName string            `json:"productPackingEnName"` // 英文配货名称
	ProductPackingName   string            `json:"productPackingName"`   // 中文配货名称
	DeclareCnName        string            `json:"declareCnName"`        // 商品中文报关名称
	DeclareEnName        string            `json:"declareEnName"`        // 商品英文报关名称
	HsCode               string            `json:"hsCode"`               // 海关编码
	ImgUrls              []string          `json:"imgUrls"`              // 商品图片url列表，第一个图片默认为主图
	DeveloperName        string            `json:"developerName"`        // 业务开发员,请输入通途erp中存在的用户名称，不存在的将不保存
	PurchaserName        string            `json:"purchaserName"`        // 采购员
	MerchantId           string            `json:"merchantId"`           // 商户ID
	ProductStatus        int               `json:"productStatus"`        // 商品状态；停售：0，在售：1，试卖：2，清仓：4
	ProductRemark        string            `json:"productRemark"`        // 产品备注
	SalesType            int               `json:"salesType"`            // 销售类型；普通销售：0，变参销售：1；暂不支持其他类型
	ProductWeight        float64           `json:"productWeight"`        // 商品重量
	ProductCurrentCost   float64           `json:"productCurrentCost"`   // 当前成本
	CategoryCode         string            `json:"categoryCode"`         // 分类名称
	Suppliers            []ProductSupplier `json:"suppliers"`            // 供应商列表，默认第一个为首选供应商
	PackageLength        float64           `json:"packageLength"`        // 包裹尺寸(长cm)
	PackageWidth         float64           `json:"packageWidth"`         // 包裹尺寸(宽cm)
	PackageHeight        float64           `json:"packageHeight"`        // 包裹尺寸(高cm)
}

type UpdateProductRequest struct {
	ProductId            string  `json:"productId"`            // 商品ID
	ProductName          string  `json:"productName"`          // 商品名
	ProductPackingEnName string  `json:"productPackingEnName"` // 英文配货名称
	ProductPackingName   string  `json:"productPackingName"`   // 中文配货名称
	DeclareCnName        string  `json:"declareCnName"`        // 商品中文报关名称
	DeclareEnName        string  `json:"declareEnName"`        // 商品英文报关名称
	HsCode               string  `json:"hsCode"`               // 海关编码
	MerchantId           string  `json:"merchantId"`           // 商户ID
	ProductStatus        int     `json:"productStatus"`        // 商品状态；停售：0，在售：1，试卖：2，清仓：4
	ProductRemark        string  `json:"productRemark"`        // 产品备注
	SalesType            int     `json:"salesType"`            // 销售类型；普通销售：0，变参销售：1；暂不支持其他类型
	ProductWeight        float64 `json:"productWeight"`        // 商品重量
	ProductCurrentCost   float64 `json:"productCurrentCost"`   // 当前成本
	PackageLength        float64 `json:"packageLength"`        // 包裹尺寸(长cm)
	PackageWidth         float64 `json:"packageWidth"`         // 包裹尺寸(宽cm)
	PackageHeight        float64 `json:"packageHeight"`        // 包裹尺寸(高cm)
}

type ProductQueryParams struct {
	CategoryName     string   `json:"category_name,omitempty"`
	MerchantId       string   `json:"merchantId"`
	PageNo           int      `json:"pageNo"`
	PageSize         int      `json:"pageSize"`
	ProductStatus    string   `json:"productStatus,omitempty"`
	ProductType      string   `json:"productType"`
	SkuAliases       []string `json:"skuAliases,omitempty"`
	Skus             []string `json:"skus,omitempty"`
	SupplierName     string   `json:"supplierName,omitempty"`
	UpdatedDateBegin string   `json:"updatedDateBegin,omitempty"`
	UpdatedDateEnd   string   `json:"updatedDateEnd,omitempty"`
}

type productResult struct {
	result
	Datas struct {
		Array    []Product `json:"array"`
		PageNo   int       `json:"pageNo"`
		PageSize int       `json:"pageSize"`
	}
}

func (s service) CreateProduct(req CreateProductRequest) error {
	type createProductResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	cpr := createProductResponse{}
	req.MerchantId = s.tongTool.MerchantId
	r, err := s.tongTool.Client.R().SetResult(&cpr).SetBody(req).Post("/openapi/tongtool/createProduct")
	if err == nil {
		if r.IsSuccess() {
			if cpr.Code != 200 {
				err = errors.New(cpr.Message)
			}
		} else {
			err = errors.New(r.Status())
		}
	}

	return err
}

func (s service) UpdateProduct(req UpdateProductRequest) error {
	type updateProductResponse struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Datas   string      `json:"datas"`
		Others  interface{} `json:"others"`
	}
	cpr := updateProductResponse{}
	req.MerchantId = s.tongTool.MerchantId
	r, err := s.tongTool.Client.R().SetResult(&cpr).SetBody(req).Post("/openapi/tongtool/updateProduct")
	if err == nil {
		if r.IsSuccess() {
			if cpr.Code != 200 {
				msg := strings.TrimSpace(cpr.Message)
				if msg == "" {
					msg = fmt.Sprintf("code: %d", cpr.Code)
				}
				err = errors.New(msg)
			}
		} else {
			err = errors.New(r.Status())
		}
	}

	return err
}

// Products 根据制定参数查询商品列表
func (s service) Products(params ProductQueryParams) (items []Product, isLastPage bool, err error) {
	if params.PageNo <= 0 {
		params.PageNo = s.tongTool.QueryDefaultValues.PageNo
	}
	if params.PageSize <= 0 {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	items = make([]Product, 0)
	res := productResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/goodsQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.HasError(res.Code); err == nil {
				items = res.Datas.Array
				isLastPage = len(items) <= params.PageSize
			}
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}

// Product 根据 SKU 或 SKU 别名查询单个商品
func (s service) Product(typ int, skus []string, isAlias bool) (item Product, err error) {
	if len(skus) == 0 {
		return item, errors.New("invalid param values")
	}

	if typ != ProductTypeVariable && typ != ProductTypeBinding {
		typ = ProductTypeNormal
	}

	params := ProductQueryParams{
		MerchantId:  s.tongTool.MerchantId,
		ProductType: strconv.Itoa(typ),
		PageNo:      1,
		PageSize:    s.tongTool.QueryDefaultValues.PageSize,
	}
	if isAlias {
		params.SkuAliases = skus
	} else {
		params.Skus = skus
	}

	inFunc := func(v string) bool {
		for _, sku := range skus {
			if strings.EqualFold(v, sku) {
				return true
			}
		}
		return false
	}
	exists := false
	for {
		products := make([]Product, 0)
		isLastPage := false
		products, isLastPage, err = s.Products(params)
		if err == nil {
			for _, p := range products {
				if isAlias {
					for _, label := range p.LabelList {
						if inFunc(label.SKULabel) {
							item = p
							exists = true
						}
					}
				} else {
					if inFunc(p.SKU) {
						item = p
						exists = true
					}
				}
				if exists {
					break
				}
			}
		}
		if isLastPage || exists || err != nil {
			break
		}
		params.PageNo++
	}

	if err == nil && !exists {
		err = errors.New("not found")
	}

	return
}

// ProductExists 根据 SKU 或 SKU 别名查询单个商品是否存在
func (s service) ProductExists(typ int, skus []string, isAlias bool) bool {
	if _, err := s.Product(typ, skus, isAlias); err == nil {
		return true
	} else {
		return false
	}
}
