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

// ProductInfo 售卖资料
type ProductInfo struct {
	BaseInfo            ProductBaseInfo          `json:"baseInfo"`            // BaseInfo
	CustomAttributeList []ProductCustomAttribute `json:"customAttributeList"` // 自定义属性
	DescribeList        []ProductDescription     `json:"describeList"`        // 描述和标题
	GalleryVideoList    []ProductVideoGallery    `json:"galleryVideoList"`    // 视频组
}

// UpdateProductRequest Todo 结构不完整，需要完善
type UpdateProductRequest struct {
	MerchantId           string        `json:"merchantId"`           // 商户号
	ProductInfoParamList []ProductInfo `json:"productInfoParamList"` // 售卖资料
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
