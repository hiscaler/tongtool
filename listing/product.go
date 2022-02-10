package listing

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
)

// 修改售卖资料
// https://open.tongtool.com/apiDoc.html#/?docId=4fdc65752aab409c93895327c65c1a81

// ProductBaseInfo 基础信息
type ProductBaseInfo struct {
	FullText            string  `json:"fullText"`            // 全文检索(分词)
	IsSimpleMode        string  `json:"isSimpleMode"`        // 图库是否简易模式(N,Y)
	IsStockProduct      string  `json:"isStockProduct"`      // 是否库存商品(N,Y)
	MerchantId          string  `json:"merchantId"`          // 商户编号
	PackageHeight       float64 `json:"packageHeight"`       // 包装高度
	PackageLength       float64 `json:"packageLength"`       // 包装长度
	PackageWeight       int     `json:"packageWeight"`       // 包装重量
	PackageWidth        float64 `json:"packageWidth"`        // 包装宽度
	PrimaryAttribute    string  `json:"primaryAttribute"`    // 主属性(橱窗图属性)
	ProductCategoryId   string  `json:"productCategoryId"`   // 产品类目顺序号
	ProductCategoryText string  `json:"productCategoryText"` // 产品类名称
	ProductHeight       float64 `json:"productHeight"`       // 商品高度
	ProductId           string  `json:"productId"`           // 产品顺序号
	ProductLength       float64 `json:"productLength"`       // 商品长度
	ProductName         string  `json:"productName"`         // 商品名称
	ProductRegisterType string  `json:"productRegisterType"` // 是否带电 1:带电
	ProductStatus       string  `json:"productStatus"`       // 商品状态0-停售,1-在售,2-试卖,3-部分停售,4-清仓库,5-部分清仓
	ProductType         string  `json:"productType"`         // 产品类型 1-单属性,2-多属性,3-捆绑,4-多属性单卖
	ProductWeight       int     `json:"productWeight"`       // 商品重量(克)
	ProductWidth        float64 `json:"productWidth"`        // 商品宽度
	PurchaseCost        string  `json:"purchaseCost"`        // 采购成本
	Responsible         string  `json:"responsible"`         // 责任人
	SKU                 string  `json:"sku"`                 // 商品编号
}

// ProductCustomAttribute 自定义属性
type ProductCustomAttribute struct {
	CreatedBy         string `json:"createdBy"`         // 创建人
	CreatedDate       string `json:"createdDate"`       // 创建时间
	CustomAttributeId string `json:"customAttributeId"` // 自定义属性顺序号
	MerchantId        string `json:"merchantId"`        // 商户编号
	ProductId         string `json:"productId"`         // 产品顺序号
	SortNo            int    `json:"sortNo"`            // 排序号
	UpdatedBy         string `json:"updatedBy"`         // 创建人
	UpdatedDate       string `json:"updatedDate"`       // 更新时间
	VariationName     string `json:"variationName"`     // 属性名称英文US
	VariationNameCN   string `json:"variationNameCN"`   // 属性名中文
	VariationNameFRA  string `json:"variationNameFRA"`  // 属性名称法语
	VariationNameGER  string `json:"variationNameGER"`  // 属性名称德语
	VariationNameIT   string `json:"variationNameIT"`   // 属性名称意大利语
	VariationNamePOL  string `json:"variationNamePOL"`  // 属性名波兰语
	VariationNamePOR  string `json:"variationNamePOR"`  // 属性名称葡萄牙语
	VariationNameSPN  string `json:"variationNameSPN"`  // 属性名称西班牙语
	VariationValue    string `json:"variationValue"`    // 属性值英文US
	VariationValueCN  string `json:"variationValueCN"`  // 属性值中文
	VariationValueFRA string `json:"variationValueFRA"` // 属性值法语
	VariationValueGER string `json:"variationValueGER"` // 属性值德语
	VariationValueIT  string `json:"variationValueIT"`  // 属性值意大利语
	VariationValuePOL string `json:"variationValuePOL"` // 属性值波兰语
	VariationValuePOR string `json:"variationValuePOR"` // 属性值葡萄牙语
	VariationValueSPN string `json:"variationValueSPN"` // 属性值西班牙语
}

type ProductTitle struct {
	GroupId        string `json:"groupId"`        // 关键词组Id
	MerchantId     string `json:"merchantId"`     // 商户编号
	ProductTitleId string `json:"productTitleId"` // 商品标题关键字顺序号
	SortNo         int    `json:"sortNo"`         // 关键顺序
	Title          string `json:"title"`          // 标题关键字
	Typ            string `json:"type"`           // 关键字类型
}

// ProductTitleGroup 关键词组
type ProductTitleGroup struct {
	GroupId           string         `json:"groupId"`           // 商品关键词组Id
	Language          string         `json:"language"`          // 语言
	MerchantId        string         `json:"merchantId"`        // 商户编号
	ProductDescribeId string         `json:"productDescribeId"` // 商品描述顺序号
	ProductTitleList  []ProductTitle `json:"productTitleList"`  // 关键词信息
	SaleAccountIds    string         `json:"saleAccountIds"`    // 适用平台与账号
	SortNo            int            `json:"sortNo"`            // 排序
}

// ProductDescription 描述和标题
type ProductDescription struct {
	BaseRichText          string              `json:"baseRichText"`          // (简单)富文本描述
	Content               string              `json:"content"`               // 富文本描述
	Highlights            string              `json:"highlights"`            // 亮点描述
	Language              string              `json:"language"`              // 语言
	MerchantId            string              `json:"merchantId"`            // 商户编号
	MobileContent         string              `json:"mobileContent"`         // 移动端描述
	PackageContent        string              `json:"packageContent"`        // 包裹信息(描述)
	ProductDescribeId     string              `json:"productDescribeId"`     // 商品描述顺序号
	ProductId             string              `json:"productId"`             // 产品顺序号
	ProductTitleGroupList []ProductTitleGroup `json:"productTitleGroupList"` // 关键词组
	TextDescribe          string              `json:"textDescribe"`          // 纯文本描述
}

// ProductVideo 视频
type ProductVideo struct {
	CreatedBy           string `json:"createdBy"`           // 创建人
	CreatedDate         string `json:"createdDate"`         // 创建时间
	ImagePath           string `json:"imagePath"`           // 图片路径
	IsUploadedAbroad    string `json:"isUploadedAbroad"`    // 是否已上传国外服务器
	MerchantId          string `json:"merchantId"`          // 商户编号
	ProductVideoGroupId string `json:"productVideoGroupId"` //	视频组顺序号
	ProductVideoId      string `json:"productVideoId"`      // 视频顺序号
	SortNo              int    `json:"sortNo"`              // 视频顺序
	UpdatedBy           string `json:"updatedBy"`           // 创建人
	UpdatedDate         string `json:"updatedDate"`         // 更新时间
	VideoPath           string `json:"videoPath"`           // 视频路径
	VideoType           string `json:"videoType"`           // 视频类型(预留字段)
}

// ProductVideoGallery 视频组
type ProductVideoGallery struct {
	CreatedBy           string         `json:"createdBy"`           // 创建人
	CreatedDate         string         `json:"createdDate"`         // 创建时间
	MerchantId          string         `json:"merchantId"`          // 商户编号
	ProductId           string         `json:"productId"`           // 产品id
	ProductVideoGroupId string         `json:"productVideoGroupId"` // 视频组顺序号
	SaleAccountIds      string         `json:"saleAccountIds"`      // 适用平台与账号
	SortNo              int            `json:"sortNo"`              // 视频组顺序
	UpdatedBy           string         `json:"updatedBy"`           // 创建人
	UpdatedDate         string         `json:"updatedDate"`         // 更新时间
	VideoGroupType      string         `json:"videoGroupType"`      // 视频组类型,A-视频库,L-视频组
	VideoList           []ProductVideo `json:"videoList"`           // 视频列表
}

// ProductGoodsVariation 售卖商品多属性信息
type ProductGoodsVariation struct {
	GoodsDetailId     string `json:"goodsDetailId"`     // 货品顺序号
	GoodsVariationId  string `json:"goodsVariationId"`  // 商品属性顺序号
	MerchantId        string `json:"merchantId"`        // 商户编号
	SortNo            int    `json:"sortNo"`            // 排序号
	VariationName     string `json:"variationName"`     // 属性名称英文US
	VariationNameAU   string `json:"variationNameAU"`   // 属性名英语AU
	VariationNameCHT  string `json:"variationNameCHT"`  // 属性名中文
	VariationNameFra  string `json:"variationNameFra"`  // 属性名称法语
	VariationNameGer  string `json:"variationNameGer"`  // 属性名称德语
	VariationNameIt   string `json:"variationNameIt"`   // 属性名称意大利语
	VariationNamePol  string `json:"variationNamePol"`  // 属性名称波兰语
	VariationNamePor  string `json:"variationNamePor"`  // 属性名称葡萄牙语
	VariationNameSpn  string `json:"variationNameSpn"`  // 属性名称西班牙语
	VariationNameUK   string `json:"variationNameUK"`   // 属性名英语uk
	VariationValue    string `json:"variationValue"`    // 属性值英文US
	VariationValueAU  string `json:"variationValueAU"`  // 属性值英语AU
	VariationValueCHT string `json:"variationValueCHT"` // 属性值中文
	VariationValueFra string `json:"variationValueFra"` // 属性值法语
	VariationValueGer string `json:"variationValueGer"` // 属性值德语
	VariationValueIt  string `json:"variationValueIt"`  // 属性值意大利语
	VariationValuePol string `json:"variationValuePol"` // 属性值波兰语
	VariationValuePor string `json:"variationValuePor"` // 属性值葡萄牙语
	VariationValueSpn string `json:"variationValueSpn"` // 属性值西班牙语
	VariationValueUK  string `json:"variationValueUK"`  // 属性值英语uk
}

// ProductStockGoodsInfo 关联库存
type ProductStockGoodsInfo struct {
	GoodsDetailId      string `json:"goodsDetailId"`      // 货品顺序号
	MappingId          string `json:"mappingId"`          // 关联序列号
	MerchantId         string `json:"merchantId"`         // 商户编号
	Quantity           int    `json:"quantity"`           // 数量
	SortNo             int    `json:"sortNo"`             // 排序号
	StockGoodsDetailId string `json:"stockGoodsDetailId"` // 库存货品顺序号
}

// ProductGoodsInfo 售卖商品多属性信息
type ProductGoodsInfo struct {
	DistributeGoodsDetailId string                  `json:"distributeGoodsDetailId"` // 分销货品顺序号
	GoodHeight              float64                 `json:"goodHeight"`              // 货品高度
	GoodLength              float64                 `json:"goodLength"`              // 货品长度
	GoodPurchaseCost        string                  `json:"goodPurchaseCost"`        // 货品采购成本
	GoodWeight              int                     `json:"goodWeight"`              // 货品重量 (克)
	GoodWidth               float64                 `json:"goodWidth"`               // 货品宽度
	GoodsDetailId           string                  `json:"goodsDetailId"`           // 货品顺序号
	GoodsVariationList      []ProductGoodsVariation `json:"goodsVariationList"`      // 售卖商品多属性信息
	MerchantId              string                  `json:"merchantId"`              // 商户编号
	ProductId               string                  `json:"productId"`               // 图片组顺序号
	SKU                     string                  `json:"sku"`                     // 属性SKU
	SortNo                  int                     `json:"sortNo"`                  // 排序号
	StockGoodsDetailId      string                  `json:"stockGoodsDetailId"`      // 售卖货品顺序号
	StockGoodsInfos         []ProductStockGoodsInfo `json:"stockGoodsInfos"`         // 关联库存
	SupplierSKU             string                  `json:"supplierSku"`             // supplierSku
}

// ProductImage 商品图片
type ProductImage struct {
	ImageAddress         string `json:"imageAddress"`         // 图片地址
	ImageType            string `json:"imageType"`            // 图片类型(W-橱窗图,D-描述图)
	IsUploadedAbroad     string `json:"isUploadedAbroad"`     // 是否已上传国外服务器
	MerchantId           string `json:"merchantId"`           // 商户编号
	OriginalImageAddress string `json:"originalImageAddress"` // 原始图片地址
	ProductImageGroupId  string `json:"productImageGroupId"`  // 图片组顺序号
	ProductImageId       string `json:"productImageId"`       // 图片顺序号
	SortNo               int    `json:"sortNo"`               // 图片顺序
}

// ProductImageGroup 图片分组
type ProductImageGroup struct {
	ImageGroupType      string         `json:"imageGroupType"`      // 图片组类型,A-橱窗图,D-描述图(详情图)，V-多属性图,L-图片组列表,
	MerchantId          string         `json:"merchantId"`          // 商户编号
	ProductId           string         `json:"productId"`           // 图片组顺序号
	ProductImageGroupId string         `json:"productImageGroupId"` // 图片组顺序号
	ProductImageList    []ProductImage `json:"productImageList"`    // 图片列表
	SaleAccountIds      string         `json:"saleAccountIds"`      // 适用平台与账号
	SortNo              int            `json:"sortNo"`              // 图片组顺序
}

// ProductLabelName 标签
type ProductLabelName struct {
	LabelName      string `json:"labelName"`      // 规格名称
	MerchantId     string `json:"merchantId"`     // 商户编号
	ProductId      string `json:"productId"`      // 产品顺序号
	ProductLabelId string `json:"productLabelId"` // 商品标签顺序号
}

// ProductMonitor 来源
type ProductMonitor struct {
	MerchantId       string `json:"merchantId"`       // 商户编号
	MonitorLink      string `json:"monitorLink"`      // 来源URL
	ProductId        string `json:"productId"`        // 库存产品编号
	ProductMonitorId string `json:"productMonitorId"` // 库存产品来源URL编号
}

// ProductNote 备注
type ProductNote struct {
	Content       string `json:"content"`       // 备注内容
	MerchantId    string `json:"merchantId"`    // 商户编号
	ProductId     string `json:"productId"`     // 产品顺序号
	ProductNoteId string `json:"productNoteId"` // 产品备注顺序号
}

// ProductVariationImage 主属性图片列表
type ProductVariationImage struct {
	ImageAddress   string `json:"imageAddress"`   // 图片地址
	ImageId        string `json:"imageId"`        // 顺序号
	MerchantId     string `json:"merchantId"`     // 商户编号
	ProductId      string `json:"productId"`      // 产品顺序号
	SortNo         int    `json:"sortNo"`         // 图片顺序
	VariationValue string `json:"variationValue"` // 主属性值
}

// Product 售卖资料
type Product struct {
	BaseInfo            ProductBaseInfo          `json:"baseInfo"`            // BaseInfo
	CustomAttributeList []ProductCustomAttribute `json:"customAttributeList"` // 自定义属性
	DescribeList        []ProductDescription     `json:"describeList"`        // 描述和标题
	GalleryVideoList    []ProductVideoGallery    `json:"galleryVideoList"`    // 视频组
	GoodsInfoList       []ProductGoodsInfo       `json:"goodsInfoList"`       // 售卖商品多属性信息
	ImageGroupList      []ProductImageGroup      `json:"imageGroupList"`      // 图片分组列表
	ImageList           []ProductImageGroup      `json:"imageList"`           // 图片列表
	LabelNames          []ProductLabelName       `json:"labelNames"`          // 标签
	MonitorList         []ProductMonitor         `json:"monitorList"`         // 来源
	NoteList            []ProductNote            `json:"noteList"`            // 备注
	VariationImagesList []ProductVariationImage  `json:"variationImagesList"` // 主属性图片列表
}

// UpdateProductRequest 商品更新请求
type UpdateProductRequest struct {
	MerchantId           string    `json:"merchantId"`           // 商户号
	ProductInfoParamList []Product `json:"productInfoParamList"` // 售卖资料
}

func (m UpdateProductRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductInfoParamList, validation.Required.Error("售卖资料不能为空")),
	)
}

// UpdateProduct 修改售卖资料
func (s service) UpdateProduct(req UpdateProductRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/product/updateProductInfo")
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

// 删除售卖资料
// https://open.tongtool.com/apiDoc.html#/?docId=591e82951e8542018fabfb16f2a3764d

type DeleteProductRequest struct {
	MerchantId string   `json:"merchantId"` // 商户编号
	ProductIds []string `json:"productIds"` // 售卖产品id
}

func (m DeleteProductRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductIds, validation.Required.Error("售卖产品id不能为空")),
	)
}

// DeleteProduct 删除售卖资料
func (s service) DeleteProduct(req DeleteProductRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/product/deleteProductInfo")
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

// 批量获取售卖详情
// https://open.tongtool.com/apiDoc.html#/?docId=9af3a1d913b7431b8605a0cdae5c6aeb

type ProductQueryParams struct {
	MerchantId    string   `json:"merchantId"`    // 商户号
	ProductIdList []string `json:"productIdList"` // 产品Id
	SKUList       []string `json:"skuList"`       // 产品sku
}

func (m ProductQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductIdList, validation.When(len(m.SKUList) == 0, validation.Required.Error("产品Id不能为空"))),
		validation.Field(&m.SKUList, validation.When(len(m.ProductIdList) == 0, validation.Required.Error("产品sku不能为空"))),
	)
}

// Products 批量获取售卖详情
func (s service) Products(req ProductQueryParams) (items []Product, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		tongtool.Response
		Datas []Product `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/product/getProductInfoByParamList")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas
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
