package erp3

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hiscaler/tongtool"
	"strings"
	"time"
)

// 商品状态
const (
	ProductStatusHaltSales     = "0" // 0 停售
	ProductStatusOnSale        = "1" // 1 在售
	ProductStatusTrySale       = "2" // 2 试卖
	ProductStatusClearanceSale = "4" // 4 清仓
)

// 销售类型
const (
	ProductTypeNormal   = iota // 0 普通销售
	ProductTypeVariable        // 1 变参销售
	ProductTypeBinding         // 2 捆绑销售
)

const (
	ProductSaleTypeNormal   = "0" // 普通销售
	ProductSaleTypeVariable = "1" // 变参销售
)

// 详细描述语言
const (
	ProductDetailDescriptionLanguageDe   = "de-de" // 德语
	ProductDetailDescriptionLanguageEnGb = "en-gb" // 英语(英国)
	ProductDetailDescriptionLanguageEnUs = "en-us" // 英语(美国)
	ProductDetailDescriptionLanguageEs   = "es-es" // 西班牙语
	ProductDetailDescriptionLanguageFr   = "fr-fr" // 法语
	ProductDetailDescriptionLanguageIt   = "it-it" // 意大利语
	ProductDetailDescriptionLanguagePl   = "pl-pl" // 波兰语
	ProductDetailDescriptionLanguagePt   = "pt-pt" // 葡萄牙语
	ProductDetailDescriptionLanguageRu   = "ru-ru" // 俄语
	ProductDetailDescriptionLanguageZhCn = "zh-cn" // 简体中文
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
	GoodsSKU      string  `json:"goodsSku"`
	GoodsWeight   float64 `json:"goodsWeight"`
	GoodsAveCost  float64 `json:"goodsAveCost"`
	GoodsCurCost  float64 `json:"goodsCurCost"`
	GoodsDetailId string  `json:"goodsDetailId"`
}

type Label struct {
	SKULabel string `json:"skuLabel"`
}

// ProductAccessory 商品配件
type ProductAccessory struct {
	AccessoriesName     string `json:"accessoriesName"`     // 配件名称
	AccessoriesQuantity int    `json:"accessoriesQuantity"` // 配件个数
}

// ProductQualityMeasure 质检标准
type ProductQualityMeasure struct {
	ItemName    string `json:"itemName"`    // 质检项
	ItemValue   string `json:"itemValue"`   // 质检标准
	MeasureName string `json:"measureName"` // 质检类目
}

// ProductSupplier 供应商
type ProductSupplier struct {
	SupplierName        string `json:"supplierName"`
	MinPurchaseQuantity int    `json:"minPurchaseQuantity"`
	PurchaseRemark      string `json:"purchaseRemark"`
}

// ProductGoodsVariation 货品属性
type ProductGoodsVariation struct {
	VariationName  string `json:"variationName"` // 规格名称
	VariationValue string `json:"variationName"` // 规格值
}

// ProductAttribute 商品属性
type ProductAttribute struct {
	AttributeKey   string `json:"attributeKey"`   // 属性key
	AttributeValue string `json:"attributeValue"` // 配件value
}

// ProductGoods 变参货品列表，创建变参销售商品时必填
type ProductGoods struct {
	GoodsAverageCost float64                 `json:"goodsAverageCost"` // 货品平均成本
	GoodsCurrentCost float64                 `json:"goodsCurrentCost"` // 货品成本(最新成本)
	GoodsSKU         string                  `json:"goodsSku"`         // 货号(SKU)
	GoodsWeight      int                     `json:"goodsWeight"`      // 货品重量(克)
	GoodsVariation   []ProductGoodsVariation `json:"goodsVariation"`   // 货品属性列表
}

// ProductDetailDescription 详细描述
type ProductDetailDescription struct {
	Content      string `json:"content"`      // 详细描述内容
	DescLanguage string `json:"descLanguage"` // 详细描述语言，德语:de-de,英语(英国):en-gb,英语(美国):en-us,西班牙语:es-es,法语:fr-fr,意大利语:it-it,波兰语:pl-pl,葡萄牙语:pt-pt,俄语:ru-ru,简体中文:zh-cn
	Title        string `json:"title"`        // 详细描述标题
}

type CreateProductRequest struct {
	Accessories          []ProductAccessory         `json:"accessories"`          // 商品配件列表
	QualityMeasures      []ProductQualityMeasure    `json:"qualityMeasures"`      // 质检标准列表
	Suppliers            []ProductSupplier          `json:"suppliers"`            // 供应商列表，默认第一个为首选供应商
	Attributes           []ProductAttribute         `json:"attributes"`           // 商品属性列表
	DetailDescriptions   []ProductDetailDescription `json:"detailDescriptions"`   // 详细描述列表
	Goods                []ProductGoods             `json:"goods"`                // 变参货品列表，创建变参销售商品时必填
	BrandCode            string                     `json:"brandCode"`            // 品牌,请输入通途erp中存在的品牌名称，不存在的将不保存
	CategoryCode         string                     `json:"categoryCode"`         // 分类,请输入通途erp中存在的分类名称，不存在的将不保存
	DeclareCnName        string                     `json:"declareCnName"`        // 商品中文报关名称
	DeclareEnName        string                     `json:"declareEnName"`        // 商品英文报关名称
	DetailImageUrls      []string                   `json:"detailImageUrls"`      // 详细描述图片url列表
	DeveloperName        string                     `json:"developerName"`        // 业务开发员,请输入通途erp中存在的用户名称，不存在的将不保存
	EnablePackageNum     int                        `json:"enablePackageNum"`     // 可包装个数
	HsCode               string                     `json:"hsCode"`               // 海关编码
	ImgUrls              []string                   `json:"imgUrls"`              // 商品图片url列表，第一个图片默认为主图
	InquirerName         string                     `json:"inquirerName"`         // 采购询价员,请输入通途erp中存在的用户名称，不存在的将不保存
	MerchantId           string                     `json:"merchantId"`           // 商户ID
	PackageHeight        float64                    `json:"packageHeight"`        // 包裹尺寸(高cm)
	PackageLength        float64                    `json:"packageLength"`        // 包裹尺寸(长cm)
	PackageMaterial      string                     `json:"packageMaterial"`      // 包装材料
	PackageWidth         float64                    `json:"packageWidth"`         // 包裹尺寸(宽cm)
	PackagingCost        float64                    `json:"packagingCost"`        // 包装成本
	PackagingWeight      float64                    `json:"packagingWeight"`      // 商品包装重量(g)
	ProductAverageCost   float64                    `json:"productAverageCost"`   // 平均成本（CNY）
	ProductCode          string                     `json:"productCode"`          // 商品编号PCL
	ProductCurrentCost   float64                    `json:"productCurrentCost"`   // 当前成本（CNY）
	ProductFeature       string                     `json:"productFeature"`       // 产品特点
	ProductGuideCost     float64                    `json:"productGuideCost"`     // 指导成本（CNY）
	ProductHeight        float64                    `json:"productHeight"`        // 商品尺寸(高cm)
	ProductLabelIds      []string                   `json:"productLabelIds"`      // 特性标签列表,请输入通途erp中存在的特性标签，不存在的将不保存
	ProductLength        float64                    `json:"productLength"`        // 商品尺寸(长cm)
	ProductName          string                     `json:"productName"`          // 商品名
	ProductPackingEnName string                     `json:"productPackingEnName"` // 英文配货名称
	ProductPackingName   string                     `json:"productPackingName"`   // 中文配货名称
	ProductRemark        string                     `json:"productRemark"`        // 产品备注
	ProductStatus        string                     `json:"productStatus"`        // 商品状态；停售：0，在售：1，试卖：2，清仓：4
	ProductWeight        int                        `json:"productWeight"`        // 商品重量
	ProductWidth         float64                    `json:"productWidth"`         // 商品尺寸(宽cm)
	PurchaserName        string                     `json:"purchaserName"`        // 采购员,请输入通途erp中存在的用户名称，不存在的将不保存
	SalesType            string                     `json:"salesType"`            // 销售类型；普通销售：0，变参销售：1；暂不支持其他类型
}

type UpdateProductRequest struct {
	DeclareCnName        string  `json:"declareCnName"`        // 商品中文报关名称
	DeclareEnName        string  `json:"declareEnName"`        // 商品英文报关名称
	EnablePackageNum     int     `json:"enablePackageNum"`     // 可包装个数
	HsCode               string  `json:"hsCode"`               // 海关编码
	MerchantId           string  `json:"merchantId"`           // 商户ID
	PackageHeight        float64 `json:"packageHeight"`        // 包裹尺寸(高cm)
	PackageLength        float64 `json:"packageLength"`        // 包裹尺寸(长cm)
	PackageWidth         float64 `json:"packageWidth"`         // 包裹尺寸(宽cm)
	PackagingCost        float64 `json:"packagingCost"`        // 包装成本
	PackagingWeight      float64 `json:"packagingWeight"`      // 商品包装重量(g)
	ProductAverageCost   float64 `json:"productAverageCost"`   // 平均成本（CNY），变参销售不支持修改
	ProductCurrentCost   float64 `json:"productCurrentCost"`   // 当前成本（CNY），变参销售不支持修改
	ProductFeature       string  `json:"productFeature"`       // 产品特点
	ProductGuideCost     float64 `json:"productGuideCost"`     // 指导成本（CNY），变参销售不支持修改
	ProductHeight        float64 `json:"productHeight"`        // 商品尺寸(高cm)
	ProductId            string  `json:"productId"`            // 商品ID
	ProductLength        float64 `json:"productLength"`        // 商品尺寸(长cm)
	ProductName          string  `json:"productName"`          // 商品名
	ProductPackingEnName string  `json:"productPackingEnName"` // 英文配货名称
	ProductPackingName   string  `json:"productPackingName"`   // 中文配货名称
	ProductRemark        string  `json:"productRemark"`        // 产品备注
	ProductStatus        string  `json:"productStatus"`        // 商品状态；停售：0，在售：1，试卖：2，清仓：4
	ProductWeight        int     `json:"productWeight"`        // 商品重量
	ProductWidth         float64 `json:"productWidth"`         // 商品尺寸(宽cm)
	SalesType            string  `json:"salesType"`            // 销售类型；普通销售：0，变参销售：1；暂不支持其他类型
}

type ProductQueryParams struct {
	MerchantId         string   `json:"merchantId"`
	NextToken          string   `json:"nextToken,omitempty"`
	ProductCategoryId  string   `json:"product_category_id,omitempty"`
	ProductGoodsIdList []string `json:"productGoodsIdList,omitempty"`
	QueryType          string   `json:"queryType,omitempty"`
	SKUList            []string `json:"skuList,omitempty"`
	UpdatedStartTime   string   `json:"updatedStartTime,omitempty"`
	UpdatedEndTime     string   `json:"updatedEndTime,omitempty"`
	PageSize           int      `json:"pageSize"`
}

type productResult struct {
	result
	Datas struct {
		Array    []Product `json:"array"`
		PageNo   int       `json:"pageNo"`
		PageSize int       `json:"pageSize"`
	} `json:"datas,omitempty"`
}

// CreateProduct 创建商品
// https://open.tongtool.com/apiDoc.html#/?docId=43a41f3680e04756a122d8671f2fc0ca
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
			if cpr.Code != tongtool.OK {
				err = errors.New(cpr.Message)
			}
		} else {
			err = errors.New(r.Status())
		}
	}

	return err
}

// UpdateProduct 更新商品
// https://open.tongtool.com/apiDoc.html#/?docId=a928207c94184649be852b120a9f4044
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
			if cpr.Code != tongtool.OK {
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

// Products 根据指定参数查询商品列表
// https://open.tongtool.com/apiDoc.html#/?docId=919e8fff6c8047deb77661f4d8c92a3a
func (s service) Products(params ProductQueryParams) (items []Product, nextToken string, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	if err != nil {
		return
	}
	items = make([]Product, 0)
	res := productResult{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/product/query")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				isLastPage = len(items) <= params.PageSize
			}
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}
	return
}
