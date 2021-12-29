package erp2

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/pkg/in"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	StatusBoolean        bool            `json:"statusBoolean"`        // 商品删除状态布尔值
	SupplierName         string          `json:"supplier_name"`        // 供应商名称
	UpdatedDate          int             `json:"updatedDate"`          // 产品信息修改时间
}

// GoodsDetailIndex 商品详情下标值
func (p Product) GoodsDetailIndex() int {
	index := -1
	for i, d := range p.GoodsDetail {
		if strings.EqualFold(d.GoodsSKU, p.SKU) {
			index = i
			break
		}
	}
	return index
}

// Image 商品图片
func (p Product) Image() (path string) {
	index := p.GoodsDetailIndex()
	if index != -1 && len(p.ProductImgList) > 0 && (index+1) <= len(p.ProductImgList) {
		path = p.ProductImgList[index].ImageGroupId
	}
	return
}

// SaveImage 下载并保存图片
func (p Product) SaveImage(saveDir string) (imagePath string, err error) {
	img := p.Image()
	if img == "" {
		err = errors.New("image path is empty")
		return
	}
	response, err := http.Get(img)
	if err == nil {
		defer response.Body.Close()
		var b []byte
		b, err = ioutil.ReadAll(response.Body)
		if err == nil {
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
			name := slug.Make(p.SKU)
			dirs := []string{saveDir}
			for i := 0; i < len(name); i += 2 {
				dirs = append(dirs, name[i:i+2])
				if len(dirs) >= 3 {
					break
				}
			}
			filename := path.Join(dirs...)

			dirExists := false
			fi, e := os.Stat(filename)
			if !os.IsNotExist(e) {
				dirExists = !fi.IsDir()
			}

			if !dirExists {
				if err = os.MkdirAll(filename, os.ModePerm); err != nil {
					return
				}
			}

			randomNumberFunc := func(len int) string {
				str := "0123456789"
				number := ""
				bigInt := big.NewInt(int64(bytes.NewBufferString(str).Len()))
				for i := 0; i < len; i++ {
					randomInt, _ := rand.Int(rand.Reader, bigInt)
					number += string(str[randomInt.Int64()])
				}
				return number
			}

			imagePath = path.Join(filename, fmt.Sprintf("%s-%s%s", name, randomNumberFunc(8), imageExt))
			err = os.WriteFile(imagePath, b, 0666)
		}
	}
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
	ProductWidth         float64                    `json:"ProductWidth"`         // 商品尺寸(宽cm)
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
	CategoryName     string   `json:"category_name,omitempty"`
	MerchantId       string   `json:"merchantId"`
	PageNo           int      `json:"pageNo"`
	PageSize         int      `json:"pageSize"`
	ProductStatus    string   `json:"productStatus,omitempty"`
	ProductType      string   `json:"productType"`
	SKUAliases       []string `json:"skuAliases,omitempty"`
	SKUs             []string `json:"skus,omitempty"`
	SupplierName     string   `json:"supplierName,omitempty"`
	UpdatedDateBegin string   `json:"updatedDateBegin,omitempty"`
	UpdatedDateEnd   string   `json:"updatedDateEnd,omitempty"`
}

// CreateProduct 创建商品
// https://open.tongtool.com/apiDoc.html#/?docId=43a41f3680e04756a122d8671f2fc0ca
func (s service) CreateProduct(req CreateProductRequest) error {
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

// UpdateProduct 更新商品
// https://open.tongtool.com/apiDoc.html#/?docId=a928207c94184649be852b120a9f4044
func (s service) UpdateProduct(req UpdateProductRequest) error {
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

// Products 根据指定参数查询商品列表
// https://open.tongtool.com/apiDoc.html#/?docId=919e8fff6c8047deb77661f4d8c92a3a
func (s service) Products(params ProductQueryParams) (items []Product, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	if len(params.SKUs) > 10 {
		err = errors.New("skus 参数长度不能大于 10 个")
	} else if len(params.SKUAliases) > 10 {
		err = errors.New("skuAliases 参数长度不能大于 10 个")
	}
	if err != nil {
		return
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
					items[i].StatusBoolean = in.StringIn(item.Status, "0", "null", "")
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
	return
}

// Product 根据 SKU 或 SKU 别名查询单个商品
func (s service) Product(typ string, sku string, isAlias bool) (item Product, err error) {
	if len(sku) == 0 {
		return item, errors.New("invalid param values")
	}

	if !in.StringIn(typ, ProductTypeNormal, ProductTypeVariable, ProductTypeBinding) {
		typ = ProductTypeNormal
	}

	params := ProductQueryParams{
		MerchantId:  s.tongTool.MerchantId,
		ProductType: typ,
		PageNo:      1,
		PageSize:    s.tongTool.QueryDefaultValues.PageSize,
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
						if strings.EqualFold(sku, p.SKU) {
							exists = true
							item = p
						}
					default:
						if isAlias {
							for _, label := range p.LabelList {
								if strings.EqualFold(sku, label.SKULabel) && p.StatusBoolean {
									exists = true
									item = p
								}
							}
						} else {
							if strings.EqualFold(sku, p.SKU) && p.StatusBoolean {
								exists = true
								item = p
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
