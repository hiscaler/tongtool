package listing

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	jsoniter "github.com/json-iterator/go"
)

const (
	ProductTag                   = "product"                   // 产品资料标签
	EbayDraftTag                 = "ebayDraft"                 // eBay草稿标签
	EbayListingTag               = "ebayListing"               // eBay在线标签
	AliExpressDraftTag           = "aliexpressDraft"           // 速卖通草稿标签
	AliExpressListingTag         = "aliexpressListing"         // 速卖通在线标签
	AmazonDraftTag               = "amazonDraft"               // 亚马逊草稿标签
	AmazonListingTag             = "amazonListing"             // 亚马逊在线标签
	WishDraftTag                 = "wishDraft"                 // Wish草稿标签
	WishListingTag               = "wishListing"               // Wish在线标签
	LazadaDraftTag               = "lazadaDraft"               // Lazada草稿标签
	LazadaListingTag             = "lazadaListing"             // Lazada在线标签
	PmDraftTag                   = "pmDraft"                   // PM草稿标签
	PmListingTag                 = "pmListing"                 // PM在线标签
	NeweggDraftTag               = "neweggDraft"               // newegg草稿标签
	NeweggListing                = "neweggListing"             // newegg在线标签
	MercadolibreDraftTag         = "mercadolibreDraft"         // Mercadolibre草稿标签
	MercadolibreListingTag       = "mercadolibreListing"       // Mercadolibre在线标签
	ShopeeDraftTag               = "shopeeDraft"               // shopee草稿标签
	ShopeeListingTag             = "shopeeListing"             // shopee在线标签
	JoomDraftTag                 = "joomDraft"                 // joom草稿标签
	JoomListingTag               = "joomListing"               // joom在线标签
	YandexDraftTag               = "yandexDraft"               // yandex草稿标签
	TeezilyDraftTag              = "teezilyDraft"              // teezily草稿标签
	TeezilyListingTag            = "teezilyListing"            // Teezily在线标签
	ListingTemplateTag           = "listingTemplate"           // 刊登模板标签
	ShopifyDraftTag              = "shopifyDraft"              // Shopify草稿标签
	ShopifyListingTag            = "shopifyListing"            // Shopify在线标签
	JdeptDraftTag                = "jdeptDraft"                // 京东全球售草稿标签
	JdeptListingTag              = "jdeptListing"              // 京东全球售在线标签
	jdidDraftTag                 = "jdidDraft"                 // 京东全球售在线标签
	JdidListingTag               = "jdidListing"               // 京东印尼在线标签
	ShoplineDraftTag             = "shoplineDraft"             // SHOPLINE草稿标签
	ShoplineListingTag           = "shoplineListing"           // Shopline在线标签
	ShoplusDraftTag              = "shoplusDraft"              // SHOPLUS草稿标签
	ShoplusListingTag            = "shoplusListing"            // SHOPLUS在线标签
	VovaDraftTag                 = "vovaDraft"                 // Vova草稿标签
	VovaListingTag               = "vovaListing"               // Vova在线标签
	EtsyDraftTag                 = "etsyDraft"                 // Etsy草稿标签
	EtsyListingTag               = "etsyListing"               // Etsy在线标签
	DataCollectionTag            = "dataCollection"            // 数据采集标签
	ThisshopDraftTag             = "thisshopDraft"             // Thisshop草稿标签
	ThisshopListingTag           = "thisshopListing"           // Thisshop在线标签
	MycomDraftTag                = "mycomDraft"                // Mycom草稿标签
	MycomListingTag              = "mycomListing"              // Mycom在线标签
	AllegroDraftTag              = "allegroDraft"              // Allegro草稿标签
	AllegroListingTag            = "allegroListing"            // Allegro在线标签
	GlobalMercadolibreDraftTag   = "globalMercadolibreDraft"   // GlobalMercadolibre草稿标签
	GlobalMercadolibreListingTag = "globalMercadolibreListing" // GlobalMercadolibre在线标签
	PassfeedDraftTag             = "passfeedDraft"             // passfeed 在线标签
	passfeedListingTag           = "passfeedListing"           // passfeed 在线标签
	AlibabagjDraftTag            = "alibabagjDraft"            // 阿里巴巴国际草稿标签
	AlibabagjListingTag          = "alibabagjListing"          // 阿里巴巴国际在线标签
	ShoplazzaDraftTag            = "shoplazzaDraft"            // 店匠草稿标签
	ShoplazzaListingTag          = "shoplazzaListing"          // 店匠在线标签
	TbgspDraftTag                = "tbgspDraft"                // 阿里分销草稿标签
	TbgspListingTag              = "tbgspListing"              // 阿里分销在线标签
	B2wDraftTag                  = "b2wDraft"                  // B2W草稿标签
	B2wListingTag                = "b2wListing"                // B2W草稿标签
)

type Tag struct {
	LabelId   string `json:"labelId"`   // 标签ID
	LabelName string `json:"labelName"` // 标签名称
	LabelType string `json:"labelType"` // 标签类别
}

// 标签列表
// https://open.tongtool.com/apiDoc.html#/?docId=f22b1937adf04312974634495a9bbb6e

type TagsQueryParams struct {
	LabelId    string `json:"labelId,omitempty"`   // 标签ID
	LabelName  string `json:"labelName,omitempty"` // 标签名称
	LabelType  string `json:"labelType,omitempty"` // 标签类别
	MerchantId string `json:"merchantId"`          // 商户号
}

func (m TagsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.LabelType, validation.When(m.LabelType != "", validation.In(
			ProductTag,
			EbayDraftTag, EbayListingTag,
			AliExpressDraftTag, AliExpressListingTag,
			AmazonDraftTag, AmazonListingTag,
			WishDraftTag, WishListingTag,
			LazadaDraftTag, LazadaListingTag,
			PmDraftTag, PmListingTag,
			NeweggDraftTag, NeweggListing,
			MercadolibreDraftTag, MercadolibreListingTag,
			ShopeeDraftTag, ShopeeListingTag,
			JoomDraftTag, JoomListingTag,
			YandexDraftTag,
			TeezilyDraftTag, TeezilyListingTag,
			ListingTemplateTag,
			ShopifyDraftTag, ShopifyListingTag,
			JdeptDraftTag, JdeptListingTag,
			jdidDraftTag, JdidListingTag,
			ShoplineDraftTag, ShoplineListingTag,
			ShoplusDraftTag, ShoplusListingTag,
			VovaDraftTag, VovaListingTag,
			EtsyDraftTag, EtsyListingTag,
			DataCollectionTag,
			ThisshopDraftTag, ThisshopListingTag,
			MycomDraftTag, MycomListingTag,
			AllegroDraftTag, AllegroListingTag,
			GlobalMercadolibreDraftTag, GlobalMercadolibreListingTag,
			PassfeedDraftTag, passfeedListingTag,
			AlibabagjDraftTag, AlibabagjListingTag,
			ShoplazzaDraftTag, ShoplazzaListingTag,
			TbgspDraftTag, TbgspListingTag,
			B2wDraftTag, B2wListingTag,
		).Error("无效的标签类别"))),
	)
}

func (s service) Tags(params TagsQueryParams) (items []Tag, err error) {
	params.MerchantId = s.tongTool.MerchantId
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = keyx.Generate(params)
		if b, e := s.tongTool.Cache.Get(cacheKey); e == nil {
			if e = jsoniter.Unmarshal(b, &items); e == nil {
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
	items = make([]Tag, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []Tag `json:"array"`
			PageNo   int   `json:"pageNo"`
			PageSize int   `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/listing/productTag/getProductTag")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Array
		}
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	if err != nil {
		return
	}

	if s.tongTool.EnableCache && len(items) > 0 {
		if b, e := jsoniter.Marshal(&items); e == nil {
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

// 添加标签
// https://open.tongtool.com/apiDoc.html#/?docId=2daeddee878b457ea736a4de7b67717d

type CreateTagRequest struct {
	LabelId    string `json:"labelId"`    // 标签ID
	LabelName  string `json:"labelName"`  // 标签名称
	LabelType  string `json:"labelType"`  // 标签类别
	MerchantId string `json:"merchantId"` // 商户号
}

func (m CreateTagRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.LabelName, validation.Required.Error("标签名称不能为空")),
		validation.Field(&m.LabelType, validation.Required.Error("标签类别不能为空")),
	)
}

func (s service) CreateTag(req CreateTagRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := tongtool.Response{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/productTag/createProductTag")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		err = tongtool.ErrorWrap(res.Code, res.Message)
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}

// 替换标签库标签
// https://open.tongtool.com/apiDoc.html#/?docId=453ecba3472f41da85eb76bbef08da4a

type UpdateTagRequest struct {
	DestinationLabelName string `json:"destinationLabelName"` // 目的标签名称
	LabelType            string `json:"labelType"`            // 标签类型
	MerchantId           string `json:"merchantId"`           // 商户号
	OriginalLabelName    string `json:"originalLabelName"`    // 原标签名称
}

func (m UpdateTagRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DestinationLabelName, validation.Required.Error("目的标签名称不能为空")),
		validation.Field(&m.LabelType, validation.Required.Error("标签类别不能为空")),
		validation.Field(&m.OriginalLabelName, validation.Required.Error("原标签名称不能为空")),
	)
}

func (s service) UpdateTag(req UpdateTagRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := tongtool.Response{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/productTag/replaceLabelLibrary")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		err = tongtool.ErrorWrap(res.Code, res.Message)
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}

// 删除标签
// https://open.tongtool.com/apiDoc.html#/?docId=9ebbd46fa32642c58e56cb9e7383050a

type DeleteTagRequest struct {
	LabelId    string `json:"labelId"`    // 标签ID
	LabelType  string `json:"labelType"`  // 标签类别
	MerchantId string `json:"merchantId"` // 商户号
}

func (m DeleteTagRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.LabelId, validation.Required.Error("标签 ID 不能为空")),
	)
}

func (s service) DeleteTag(req DeleteTagRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.MerchantId = s.tongTool.MerchantId
	res := tongtool.Response{}
	resp, err := s.tongTool.Client.R().
		SetResult(&res).
		SetBody(req).
		Post("/openapi/tongtool/listing/productTag/removeProductTag")
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		err = tongtool.ErrorWrap(res.Code, res.Message)
	} else {
		if e := jsoniter.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return err
}
