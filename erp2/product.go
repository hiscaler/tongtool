package erp2

import (
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gosimple/slug"
	"github.com/hiscaler/gox/filex"
	"github.com/hiscaler/gox/inx"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/gox/randx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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
	ProductTypeNormal   = "0" // 0 普通销售
	ProductTypeVariable = "1" // 1 变参销售
	ProductTypeBinding  = "2" // 2 捆绑销售
)

const (
	ProductSaleTypeNormal   = "0" // 普通销售
	ProductSaleTypeVariable = "1" // 变参销售
)

// 详细描述语言
const (
	ProductDetailDescriptionGerman            = "de-de" // 德语
	ProductDetailDescriptionBritishEnglish    = "en-gb" // 英语(英国)
	ProductDetailDescriptionAmericanEnglish   = "en-us" // 英语(美国)
	ProductDetailDescriptionSpanish           = "es-es" // 西班牙语
	ProductDetailDescriptionFrench            = "fr-fr" // 法语
	ProductDetailDescriptionItalian           = "it-it" // 意大利语
	ProductDetailDescriptionPolish            = "pl-pl" // 波兰语
	ProductDetailDescriptionPortuguese        = "pt-pt" // 葡萄牙语
	ProductDetailDescriptionRussian           = "ru-ru" // 俄语
	ProductDetailDescriptionSimplifiedChinese = "zh-cn" // 简体中文
)

// Product 通途商品
type Product struct {
	BrandName            string          `json:"brandName"`            // 品牌名称
	CategoryName         string          `json:"categoryName"`         // 分类名称
	CreatedDate          int             `json:"createdDate"`          // 产品创建时间
	DeclareCnName        string          `json:"declareCnName"`        // 商品中文报关名称
	DeclareEnName        string          `json:"declareEnName"`        // 商品英文报关名称
	DeveloperName        string          `json:"developerName"`        // 业务开发员名称
	EnablePackageNum     int             `json:"enablePackageNum"`     // 可包装个数
	GoodsDetail          []ProductDetail `json:"goodsDetail"`          // 商品明细
	HsCode               string          `json:"hsCode"`               // 海关编码
	InquirerName         string          `json:"inquirerName"`         // 采购询价员名称
	LabelList            []Label         `json:"labelList"`            // 产品标签
	LabelName            string          `json:"labelName"`            // 标签名称
	PackageCost          float64         `json:"packageCost"`          // 商品包装成本
	PackageHeight        float64         `json:"packageHeight"`        // 包裹尺寸(高cm)
	PackageLength        float64         `json:"packageLength"`        // 包裹尺寸(长cm)
	PackageMaterialName  string          `json:"packageMaterialName"`  // 包装名称
	PackageWeight        float64         `json:"packageWeight"`        // 商品包装重量(克)
	PackageWidth         float64         `json:"packageWidth"`         // 包裹尺寸(宽cm)
	ProductCode          string          `json:"productCode"`          // 商品编号
	ProductFeature       string          `json:"productFeature"`       // 产品特点
	ProductHeight        float64         `json:"productHeight"`        // 商品尺寸(高cm)
	ProductImgList       []ProductImage  `json:"productImgList"`       // 产品图片
	ProductLength        float64         `json:"productLength"`        // 商品尺寸(长cm)
	ProductName          string          `json:"productName"`          // 产品名称
	ProductPackingEnName string          `json:"productPackingEnName"` // 商品英文报关名称
	ProductPackingName   string          `json:"productPackingName"`   // 中文配货名称
	ProductWidth         float64         `json:"productWidth"`         // 商品尺寸(宽cm)
	ProductId            string          `json:"product_id"`           // 产品id
	PurchaseName         string          `json:"purchaseName"`         // 采购员名称
	PurchaserId          string          `json:"purchaserId"`          // 采购员id
	SKU                  string          `json:"sku"`                  // 商品sku
	Status               string          `json:"status"`               // 商品删除状态,1:删除,null或0：未删除
	SupplierName         string          `json:"supplierName"`         // 供应商名称
	UpdatedDate          int             `json:"updatedDate"`          // 产品信息修改时间
	// 自定义字段
	IsDeleted bool `json:"isDeleted"` // 商品是否删除
}

// GoodsDetailIndex 商品详情下标值
func (p Product) GoodsDetailIndex() int {
	for k, v := range p.GoodsDetail {
		if strings.EqualFold(v.GoodsSKU, p.SKU) {
			return k
		}
	}
	return -1
}

// Image 商品图片
func (p Product) Image() (path string) {
	n := len(p.ProductImgList)
	if n == 0 {
		return
	}
	index := p.GoodsDetailIndex()
	if index >= 0 && index < n {
		// 0 ~ n-1
		path = p.ProductImgList[index].ImageGroupId
	}
	return
}

// ImageIsNormalized 图片地址是否规范
// 判断依据：必须为网络地址（不会判断地址的有效性），且后续的地址中只能包含大小写字母、数字、点、横杠、左斜线
// 作用：通途的图片地址中可能包含中文和其他一些特殊字符，导致在数据在某些场景下不能正确地使用，所以如果地址不是规范的，可以使用 SaveImage 方法将图片下载保存到本地再使用
func (p Product) ImageIsNormalized() bool {
	s := p.Image()
	if s == "" {
		return false
	}

	re, _ := regexp.Compile(`^(?i)http[s]?://[a-zA-Z0-9.-/]+$`)
	return re.MatchString(s)
}

// SaveImage 下载并保存图片
func (p Product) SaveImage(saveDir string) (imagePath string, err error) {
	img := p.Image()
	if img == "" {
		err = errors.New("image path is empty")
		return
	}
	response, err := http.Get(img)
	if err != nil {
		return
	}

	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var imageExt string
	switch http.DetectContentType(b) {
	case "image/jpeg":
		imageExt = ".jpg"
	case "image/png":
		imageExt = ".png"
	case "image/gif":
		imageExt = ".gif"
	case "image/bmp":
		imageExt = ".bmp"
	case "image/webp":
		imageExt = ".webp"
	default:
		imageExt = filepath.Ext(img)
	}
	replacer := strings.NewReplacer("-", "", "_", "")
	name := replacer.Replace(slug.Make(p.SKU))
	maxDirLevels := 2
	dirs := make([]string, 0)
	if saveDir != "" {
		maxDirLevels++
		dirs = append(dirs, saveDir)
	}

	n := len(name)
	for i := 0; i < n; i += 2 {
		j := 2
		if i >= n {
			j = 1
		}
		dirs = append(dirs, name[i:j])
		if len(dirs) >= maxDirLevels {
			break
		}
	}
	filename := path.Join(dirs...)
	if !filex.Exists(filename) {
		if err = os.MkdirAll(filename, os.ModePerm); err != nil {
			return
		}
	}

	imagePath = path.Join(filename, fmt.Sprintf("%s-%s%s", name, randx.Number(8), imageExt))
	err = os.WriteFile(imagePath, b, 0666)

	return
}

// ProductImage 商品图片
type ProductImage struct {
	ImageGroupId string `json:"imageGroupId"` // 图片url
}

// ProductDetail 通途商品详情
type ProductDetail struct {
	GoodsAveCost  float64 `json:"goodsAveCost"`  // 商品平均成本
	GoodsCurCost  float64 `json:"goodsCurCost"`  // 商品当前成本
	GoodsDetailId string  `json:"goodsDetailId"` // 货品ID
	GoodsSKU      string  `json:"goodsSku"`      // 商品sku
	GoodsWeight   float64 `json:"goodsWeight"`   // 货品重量(克)
}

type Label struct {
	SKULabel string `json:"skuLabel"` // 商品标签
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
	SupplierName        string `json:"supplierName"`        // 供应商名称
	MinPurchaseQuantity int    `json:"minPurchaseQuantity"` // 最小采购量
	PurchaseRemark      string `json:"purchaseRemark"`      // 采购备注
}

// ProductGoodsVariation 货品属性
type ProductGoodsVariation struct {
	VariationName  string `json:"variationName"`  // 规格名称
	VariationValue string `json:"variationValue"` // 规格值
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

// 创建商品

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

func (m CreateProductRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductName, validation.Required.Error("商品名称不能为空")),
		validation.Field(&m.ProductCode, validation.Required.Error("商品 SKU 不能为空")),
		validation.Field(&m.EnablePackageNum, validation.Min(1).Error("可包装数量不能小于 1")),
		validation.Field(&m.ProductStatus, validation.Required.Error("商品状态不可为空"), validation.In(ProductStatusHaltSales, ProductStatusOnSale, ProductStatusTrySale, ProductStatusClearanceSale).Error("无效的商品状态")),
		validation.Field(&m.SalesType, validation.Required.Error("销售类型不可为空"), validation.In(ProductSaleTypeNormal, ProductSaleTypeVariable).Error("无效的销售类型")),
		validation.Field(&m.Goods, validation.When(m.SalesType == ProductSaleTypeVariable, validation.Required.Error("变参货品不能为空"))),
		validation.Field(&m.Accessories, validation.When(len(m.Accessories) > 0, validation.By(func(value interface{}) error {
			items, ok := value.([]ProductAccessory)
			if !ok {
				return errors.New("无效的商品配件")
			}
			for i, item := range items {
				if item.AccessoriesName == "" {
					return fmt.Errorf("数据 %d 中配件名称不能为空", i+1)
				}
				if item.AccessoriesQuantity <= 0 {
					return fmt.Errorf("数据 %d 中配件数量不能小于 1", i+1)
				}
			}
			return nil
		}))),
		validation.Field(&m.DetailImageUrls, validation.When(len(m.DetailImageUrls) > 0), validation.Each(is.URL.Error("无效的地址"))),
		validation.Field(&m.ImgUrls, validation.When(len(m.ImgUrls) > 0), validation.Each(is.URL.Error("无效的地址"))),
		validation.Field(&m.DetailDescriptions, validation.When(len(m.DetailDescriptions) > 0, validation.By(func(value interface{}) error {
			items, ok := value.([]ProductDetailDescription)
			if !ok {
				return errors.New("无效的商品详细描述")
			}
			for _, item := range items {
				err := validation.ValidateStruct(&item,
					validation.Field(&item.Title, validation.Required.Error("详细描述标题不能为空")),
					validation.Field(&item.DescLanguage, validation.In(
						ProductDetailDescriptionGerman,
						ProductDetailDescriptionBritishEnglish,
						ProductDetailDescriptionAmericanEnglish,
						ProductDetailDescriptionSpanish,
						ProductDetailDescriptionFrench,
						ProductDetailDescriptionItalian,
						ProductDetailDescriptionPolish,
						ProductDetailDescriptionPortuguese,
						ProductDetailDescriptionRussian,
						ProductDetailDescriptionSimplifiedChinese,
					).Error("无效的描叙语言")),
					validation.Field(&item.Content, validation.Required.Error("详细描述内容不能为空")),
				)
				if err != nil {
					return err
				}
			}
			return nil
		}))),
	)
}

// CreateProduct 创建商品
// https://open.tongtool.com/apiDoc.html#/?docId=43a41f3680e04756a122d8671f2fc0ca
func (s service) CreateProduct(req CreateProductRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		result
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/createProduct")
	if err == nil {
		if resp.IsSuccess() {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	return err
}

// 更新商品

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

func (m UpdateProductRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductName, validation.Required.Error("商品名称不能为空")),
		validation.Field(&m.EnablePackageNum, validation.Min(1).Error("可包装数量不能小于 1")),
		validation.Field(&m.ProductStatus, validation.Required.Error("商品状态不可为空"), validation.In(ProductStatusHaltSales, ProductStatusOnSale, ProductStatusTrySale, ProductStatusClearanceSale).Error("无效的商品状态")),
		validation.Field(&m.SalesType, validation.Required.Error("销售类型不可为空"), validation.In(ProductSaleTypeNormal, ProductSaleTypeVariable).Error("无效的销售类型")),
	)
}

// UpdateProduct 更新商品
// https://open.tongtool.com/apiDoc.html#/?docId=a928207c94184649be852b120a9f4044
func (s service) UpdateProduct(req UpdateProductRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		result
		Datas  string      `json:"datas"`
		Others interface{} `json:"others"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/updateProduct")
	if err == nil {
		if resp.IsSuccess() {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	return err
}

// 查询商品

type ProductQueryParams struct {
	Paging
	CategoryName     string   `json:"category_name,omitempty"`    // 分类名称
	MerchantId       string   `json:"merchantId"`                 // 商户ID
	ProductStatus    string   `json:"productStatus,omitempty"`    // 商品状态：1试卖、2正常
	ProductType      string   `json:"productType"`                // 销售类型：0, 普通销售/1,变参销售/2,捆绑销售
	SKUAliases       []string `json:"skuAliases,omitempty"`       // SKU别名数组，长度不超过10
	SKUs             []string `json:"skus,omitempty"`             // SKU数组，长度不超过10
	SupplierName     string   `json:"supplierName,omitempty"`     // 供应商
	UpdatedDateBegin string   `json:"updatedDateBegin,omitempty"` // 更新时间查询的起始时间
	UpdatedDateEnd   string   `json:"updatedDateEnd,omitempty"`   // 更新时间查询的结束时间
}

func (m ProductQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductStatus, validation.In("1", "2").Error("无效的商品状态")),
		validation.Field(&m.ProductType, validation.In(ProductTypeNormal, ProductTypeVariable, ProductTypeBinding)),
		validation.Field(&m.SKUs, validation.When(len(m.SKUs) > 0, validation.By(func(value interface{}) error {
			items, ok := value.([]string)
			if !ok {
				return errors.New("无效的 SKU 数据")
			}
			if len(items) > 10 {
				return errors.New("SKU 数据不能多于 10 个")
			}
			return nil
		}))),
		validation.Field(&m.SKUAliases, validation.When(len(m.SKUAliases) > 0, validation.By(func(value interface{}) error {
			items, ok := value.([]string)
			if !ok {
				return errors.New("无效的 SKU 别名数据")
			}
			if len(items) > 10 {
				return errors.New("SKU 别名数据不能多于 10 个")
			}
			return nil
		}))),
		validation.Field(&m.UpdatedDateBegin, validation.When(m.UpdatedDateBegin != "", validation.Date(constant.DatetimeFormat).Error(" 更新时间查询的起始时间格式无效"))),
		validation.Field(&m.UpdatedDateEnd, validation.When(m.UpdatedDateEnd != "", validation.Date(constant.DatetimeFormat).Error(" 更新时间查询的结束时间格式无效"))),
	)
}

// Products 根据指定参数查询商品列表
// https://open.tongtool.com/apiDoc.html#/?docId=919e8fff6c8047deb77661f4d8c92a3a
func (s service) Products(params ProductQueryParams) (items []Product, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
	if len(params.SKUs) > 10 {
		err = errors.New("skus 参数长度不能大于 10 个")
	} else if len(params.SKUAliases) > 10 {
		err = errors.New("skuAliases 参数长度不能大于 10 个")
	}
	if err != nil {
		return
	}

	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = keyx.Generate(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = json.Unmarshal(b, &items); e == nil {
				return
			} else {
				s.tongTool.Logger.Printf(`cache data unmarshal error
 DATA: %s
ERROR: %s
`, string(b), e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("get cache %s error: %s", cacheKey, e.Error())
		}
	}
	items = make([]Product, 0)
	res := struct {
		result
		Datas struct {
			Array    []Product `json:"array"`
			PageNo   int       `json:"pageNo"`
			PageSize int       `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/goodsQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				for i, item := range items {
					items[i].IsDeleted = item.Status == "1"
				}
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

	if err == nil && s.tongTool.EnableCache && len(items) > 0 {
		if b, e := json.Marshal(&items); e == nil {
			e = s.tongTool.Cache.Set(cacheKey, b)
			if e != nil {
				s.tongTool.Logger.Printf("set cache %s error: %s", cacheKey, e.Error())
			}
		} else {
			s.tongTool.Logger.Printf("items marshal error: %s", err.Error())
		}
	}
	return
}

// Product 根据 SKU 或 SKU 别名查询单个商品
func (s service) Product(typ string, sku string, isAlias bool) (item Product, err error) {
	if !inx.StringIn(typ, ProductTypeNormal, ProductTypeVariable, ProductTypeBinding) {
		typ = ProductTypeNormal
	}

	params := ProductQueryParams{
		MerchantId:  s.tongTool.MerchantId,
		ProductType: typ,
	}
	if isAlias {
		params.SKUAliases = []string{sku}
	} else {
		params.SKUs = []string{sku}
	}

	exists := false
	for {
		items := make([]Product, 0)
		isLastPage := false
		items, isLastPage, err = s.Products(params)
		if err == nil {
			if len(items) == 0 {
				err = tongtool.ErrNotFound
			} else {
				for _, p := range items {
					switch typ {
					case ProductTypeVariable:
						if strings.EqualFold(sku, p.ProductCode) || strings.EqualFold(sku, p.SKU) {
							exists = true
							item = p
						} else {
							for _, detail := range p.GoodsDetail {
								if strings.EqualFold(sku, detail.GoodsSKU) {
									exists = true
									item = p
									break
								}
							}
						}
					default:
						if !p.IsDeleted {
							if isAlias {
								for _, label := range p.LabelList {
									if strings.EqualFold(sku, label.SKULabel) {
										exists = true
										item = p
										break
									}
								}
							} else {
								if strings.EqualFold(sku, p.SKU) {
									exists = true
									item = p
								}
							}
						}
					}

					if exists {
						break
					}
				}
			}
		}
		if isLastPage || exists || err != nil {
			break
		}
		params.PageNo++
	}

	if err == nil && !exists {
		err = tongtool.ErrNotFound
	}

	return
}

// ProductExists 根据 SKU 或 SKU 别名查询单个商品是否存在
func (s service) ProductExists(typ string, sku string, isAlias bool) bool {
	if _, err := s.Product(typ, sku, isAlias); err == nil {
		return true
	} else {
		return false
	}
}
