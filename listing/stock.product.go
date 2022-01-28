package listing

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
)

// 保存库存产品资料
// https://open.tongtool.com/apiDoc.html#/?docId=c39bd0a8801a48ea9ece608ef236e314

// StockProductBaseInfo 库存商品商品信息
type StockProductBaseInfo struct {
	CreatedBy           string  `json:"createdBy"`           // 创建人
	CreatedDate         string  `json:"createdDate"`         // 创建时间
	FromType            string  `json:"fromType"`            // 来源：L系统，N新品开发系统
	FullText            string  `json:"fullText"`            // 全文检索(分词)
	IsSimpleMode        string  `json:"isSimpleMode"`        // 图库是否简易模式:Y是，N否
	MerchantId          string  `json:"merchantId"`          // 商户编号
	PrimaryAttribute    string  `json:"primaryAttribute"`    // 主属性(橱窗图属性名称)
	ProductCategoryId   string  `json:"productCategoryId"`   // 产品类目Id
	ProductCategoryText string  `json:"productCategoryText"` // 产品类名称
	ProductHeight       float64 `json:"productHeight"`       // 商品高度
	ProductId           string  `json:"productId"`           // 商品ID
	ProductLength       int     `json:"productLength"`       // 商品长度
	ProductName         string  `json:"productName"`         // 商品名称
	ProductRegisterType string  `json:"productRegisterType"` // 是否带电 1:带电
	ProductType         string  `json:"productType"`         // 产品类型 1-单属性,2-多属性,3-捆绑,4-多属性单卖
	ProductWeight       int     `json:"productWeight"`       // 商品重量(克)
	ProductWidth        float64 `json:"productWidth"`        // 商品宽度
	PurchaseCost        float64 `json:"purchaseCost"`        // 采购成本(元)
	Responsible         string  `json:"responsible"`         // 责任人
	SKU                 string  `json:"sku"`                 // 商品编号
	UpdatedBy           string  `json:"updatedBy"`           // 修改人
	UpdatedDate         string  `json:"updatedDate"`         // 修改时间
}

// StockProductURL 库存产品来源 URL
type StockProductURL struct {
	MerchantId  string `json:"merchantId"`  // 商户编号
	MonitorLink string `json:"monitorLink"` // 来源 URL 内容
	ProductId   string `json:"productId"`   // 库存产品编号
}

// StockProductDescription 库存商品描述
type StockProductDescription struct {
	BaseRichText      string `json:"baseRichText"`      //	(简单)富文本描述
	Content           string `json:"content"`           //	富文本描述
	CreatedBy         string `json:"createdBy"`         // 创建人
	CreatedDate       string `json:"createdDate"`       // 创建时间
	Highlights        string `json:"highlights"`        // 亮点描述
	Language          string `json:"language"`          //	单选:EN英语,GER德语,FRA法语,SPN西班牙语,IT意大利语,POR葡萄牙语,CN中文,RUS俄语,TH泰语,AR阿拉伯语
	MerchantId        string `json:"merchantId"`        // 商户编号
	PackageContent    string `json:"packageContent"`    // 包裹信息(描述)
	ProductDescribeId string `json:"productDescribeId"` // 商品描述Id
	ProductId         string `json:"productId"`         // 商品ID
}

// StockProductGoodsInfo 库存货品信息
type StockProductGoodsInfo struct {
	CreatedBy        string  `json:"createdBy"`        // 创建人
	CreatedDate      string  `json:"createdDate"`      //	创建时间
	GoodHeight       float64 `json:"goodHeight"`       //	货品高度
	GoodLength       float64 `json:"goodLength"`       //	货品长度
	GoodPurchaseCost string  `json:"goodPurchaseCost"` //	货品采购成本
	GoodWeight       int     `json:"goodWeight"`       //	货品重量(克)
	GoodWidth        int     `json:"goodWidth"`        // 货品宽度
	GoodsDetailId    string  `json:"goodsDetailId"`    //	货品ID
	MerchantId       string  `json:"merchantId"`       //	商户编号
	ProductId        string  `json:"productId"`        // 图片组顺序号
	SKU              string  `json:"sku"`              // 属性SKU
	SortNo           int     `json:"sortNo"`           // 排序号(同一商品下，从1开始递增)
}

// StockProductNote 库存商品备注
type StockProductNote struct {
	Content       string `json:"content"`       // 备注内容
	CreatedBy     string `json:"createdBy"`     //	创建人
	CreatedDate   string `json:"createdDate"`   //	创建时间
	MerchantId    string `json:"merchantId"`    //	商户编号
	ProductId     string `json:"productId"`     // 商品Ids
	ProductNoteId string `json:"productNoteId"` //	产品备注Id
	UpdatedBy     string `json:"updatedBy"`     //	修改人
	UpdatedDate   string `json:"updatedDate"`   //	修改时间
}

// StockProductLabel 库存商品标签
type StockProductLabel struct {
	CreatedBy      string `json:"createdBy"`      // 创建人
	CreatedDate    string `json:"createdDate"`    // 创建时间
	LabelName      string `json:"labelName"`      // 规格名称
	MerchantId     string `json:"merchantId"`     // 商户编号
	ProductId      string `json:"productId"`      // 商品ID
	ProductLabelId string `json:"productLabelId"` // 商品标签Id
	UpdatedBy      string `json:"updatedBy"`      // 创建人
	UpdatedDate    string `json:"updatedDate"`    // 更新时间
}

// StockProductImage 库存商品图片
type StockProductImage struct {
	CreatedBy      string `json:"createdBy"`      // 创建人
	CreatedDate    string `json:"createdDate"`    // 创建时间
	ImageAddress   string `json:"imageAddress"`   // 图片地址
	ImageType      string `json:"imageType"`      // 图片类型(A-图库,D-描述图,V-属性图)
	MerchantId     string `json:"merchantId"`     // 商户Id
	ProductId      string `json:"productId"`      // 图片组顺序号
	ProductImageId string `json:"productImageId"` // 图片顺序号
	SortNo         int    `json:"sortNo"`         // 图片顺序(同一商品下从1开始递增)
	UpdatedBy      string `json:"updatedBy"`      // 创建人
	UpdatedDate    string `json:"updatedDate"`    // 更新时间
}

// StockProductVariationImage 库存产品主属性图片信息
type StockProductVariationImage struct {
	CreatedBy      string `json:"createdBy"`      // 创建人
	CreatedDate    string `json:"createdDate"`    // 创建时间
	ImageAddress   string `json:"imageAddress"`   // 图片地址
	ImageId        string `json:"imageId"`        // 主属性图片Id
	MerchantId     string `json:"merchantId"`     // 商户编号
	ProductId      string `json:"productId"`      // 商品Id
	SortNo         int    `json:"sortNo"`         // 图片顺序（同一商品下从1开始递增）
	UpdatedBy      string `json:"updatedBy"`      // 创建人
	UpdatedDate    string `json:"updatedDate"`    // 更新时间
	VariationValue string `json:"variationValue"` // 主属性值
}

// UpsertStockProductRequest 库存商品更新请求
type UpsertStockProductRequest struct {
	BaseInfo            StockProductBaseInfo         `json:"baseInfo"`            // 库存商品商品信息
	MonitorList         []StockProductURL            `json:"monitorList"`         // 库存产品来源 URL
	DescribeParamList   []StockProductDescription    `json:"describeParamList"`   // 库存商品描述
	GoodsInfoParamList  []StockProductGoodsInfo      `json:"goodsInfoParamList"`  // 库存货品信息
	NoteList            []StockProductNote           `json:"noteList"`            // 库存商品备注
	LabelList           []StockProductLabel          `json:"labelList"`           // 库存商品标签
	ImageList           []StockProductImage          `json:"imageList"`           // 库存商品图片
	VariationImagesList []StockProductVariationImage `json:"variationImagesList"` // 库存产品主属性图片信息
	DataType            string                       `json:"dataType"`            //	数据内容（"baseInfo,picture,description"）包含其中的一个或多个,逗号分隔
	MerchantId          string                       `json:"merchantId"`          // 商户编号
	RequestType         int                          `json:"requestType"`         // 请求类型-0创建，1更新
	UploadPicToTongTool bool                         `json:"uploadPicToTongtool"` // 是否上传图片至通途空间
}

func (m UpsertStockProductRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.RequestType, validation.In(0, 1).Error("错误的请求类型")),
	)
}

// UpsertStockProduct 添加/更新库存产品资料
func (s service) UpsertStockProduct(req UpsertStockProductRequest) error {
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
		Post("/openapi/tongtool/listing/stock/saveStockProductInfo")
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
