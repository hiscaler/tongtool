package erp3

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	jsoniter "github.com/json-iterator/go"
)

// 入库单列表
// https://open.tongtool.com/apiDoc.html#/?docId=214c62308f3c46b0b673a923e853dec3

// 收货记录
type receipt struct {
	BatchNumber           int    `json:"batchNumber"`           // 收货数量
	CancelQuantity        int    `json:"cancelQuantity"`        // 取消数量
	CreatedBy             string `json:"createdBy"`             // (收货)创建人
	CreatedTime           string `json:"createdTime"`           // (收货)创建时间
	GoodsSKU              string `json:"goodsSku"`              // 	收货产品SKU
	MerchantId            string `json:"merchantId"`            // 	商户编号
	ProductGoodsId        string `json:"productGoodsId"`        // 	货品流水号（货品Id）
	ReceiptBatchId        string `json:"receiptBatchId"`        // 	入库单货品明细收货批次ID
	ReceiptBatchNo        string `json:"receiptBatchNo"`        // 	批次号（业务使用LC+年月日+8位seq）
	ReceiptNo             string `json:"receiptNo"`             // 	入库单编号
	WarehouseBlockCode    string `json:"warehouseBlockCode"`    // 	仓库库区代码
	WarehouseLocationCode string `json:"warehouseLocationCode"` // 	仓库库位代码
	WarehouseLocationId   string `json:"warehouseLocationId"`   // 	库位ID
}

// 质检
type qualityInspection struct {
	CreatedBy                 string `json:"createdBy"`                 // (质检)创建人
	CreatedTime               string `json:"createdTime"`               // (质检)创建时间
	GoodsSKU                  string `json:"goodsSku"`                  // 货品SKU
	MerchantId                string `json:"merchantId"`                // 商户编号
	PassCheckNumber           int    `json:"passCheckNumber"`           // 已检通过数量
	ProblemCheckNumber        int    `json:"problemCheckNumber"`        // 已检问题数量
	ReceiptBatchCheckDetailId string `json:"receiptBatchCheckDetailId"` // 质检明细ID
	ReceiptBatchCheckDetailNo string `json:"receiptBatchCheckDetailNo"` // 质检明细编号（业务使用QC+年月日+8位seq）
	ReceiptBatchCheckNo       string `json:"receiptBatchCheckNo"`       // 质检编号
	ReceiptBatchNo            string `json:"receiptBatchNo"`            // 批次号
	ReceiptNo                 string `json:"receiptNo"`                 // 入库单号
}

// 入库单货品信息
type detail struct {
	DoneNumber          int    `json:"doneNumber"`          // 已收货数量
	ExpectedNumber      int    `json:"expectedNumber"`      // 	预期数量
	GoodsCnDesc         string `json:"goodsCnDesc"`         // 	产品中文描述
	GoodsEnDesc         string `json:"goodsEnDesc"`         // 	产品英文描述
	GoodsSKU            string `json:"goodsSku"`            // 	货品 SKU
	MerchantId          string `json:"merchantId"`          // 	商户编号
	ProductGoodsId      string `json:"productGoodsId"`      // 	货品流水号（货品Id）
	ReceiptDetailId     string `json:"receiptDetailId"`     // 	入库单货品明细ID
	ReceiptDetailStatus string `json:"receiptDetailStatus"` // 	状态（0：创建状态、1：部分收货、2：完全收货、3：超量收货）
	ReceiptNo           string `json:"receiptNo"`           // 	入库单编号
	Source              string `json:"source"`              // 	来源（P：有源入库、U：无源入库）
}

// 上架
type shelve struct {
	CreatedBy                  string `json:"createdBy"`                  // (上架)创建人
	CreatedTime                string `json:"createdTime"`                // (上架)创建时间
	GoodsSKU                   string `json:"goodsSku"`                   // 货品SKU
	MerchantId                 string `json:"merchantId"`                 // 商户编号
	ReceiptBatchNo             string `json:"receiptBatchNo"`             // 批次号
	ReceiptCheckDetailShelveNo string `json:"receiptCheckDetailShelveNo"` // 上架单编号
	ReceiptCheckShelveDetailId string `json:"receiptCheckShelveDetailId"` // 上架明细ID
	ReceiptCheckShelveDetailNo string `json:"receiptCheckShelveDetailNo"` // 上架明细编号(业务使用SJ+年月日+8位seq)
	ReceiptNo                  string `json:"receiptNo"`                  // 入库单号
	ShelveNumber               int    `json:"shelveNumber"`               // 已上架数量
	WarehouseBlockCode         string `json:"warehouseBlockCode"`         // 仓库库区代码
	WarehouseLocationCode      string `json:"warehouseLocationCode"`      // 仓库库位代码
	WarehouseLocationId        string `json:"warehouseLocationId"`        // 上架库位ID
}

type StockInSheet struct {
	AbnormalStatus    string              `json:"abnormalStatus"`    // 收货异常状态,格式为:无源入库+部分收货+超出收货（0：表示没有、1：表示有异常）
	CreatedTime       string              `json:"createdTime"`       // 创建时间
	MerchantId        string              `json:"merchantId"`        // 商户编号
	ReceiptBatchList  []receipt           `json:"receiptBatchList"`  // 收货记录
	ReceiptCheckList  []qualityInspection `json:"receiptCheckList"`  // 质检记录
	ReceiptDetailList []detail            `json:"receiptDetailList"` // 入库单货品信息
	ReceiptId         string              `json:"receiptId"`         // 入库单ID
	ReceiptNo         string              `json:"receiptNo"`         // 入库单编号(业务使用RK+年月日+8位seq)
	ReceiptShelveList []shelve            `json:"receiptShelveList"` // 上架记录
	ReceiptStatus     string              `json:"receiptStatus"`     //	入库单状态(0：入库单创建、1：入库单取消、2：入库单关闭、3：收货中)
	ReceiptType       string              `json:"receiptType"`       //	入库类型(0：采购入库、1：生产入库、2：调拨入库、3：退货入库、4：其他入库)
	ReferenceNo       string              `json:"referenceNo"`       //	参考编号
	ReferenceNo2      string              `json:"referenceNo2"`      //	参考编号2
	UpdatedTime       string              `json:"updatedTime"`       //	更新时间
	WarehouseId       string              `json:"warehouseId"`       // (收货)仓库ID
	WarehouseName     string              `json:"warehouseName"`     // (收货)仓库名称
}

type StockInSheetsQueryParams struct {
	CreatedEndTime   string `json:"createdEndTime,omitempty"`   // 创建结束时间
	CreatedStartTime string `json:"createdStartTime,omitempty"` // 创建开始时间
	MerchantId       string `json:"merchantId"`                 // 商户号
	ReceiptNo        string `json:"receiptNo,omitempty"`        // 入库单号
	UpdatedEndTime   string `json:"updatedEndTime,omitempty"`   // 更新结束时间
	UpdatedStartTime string `json:"updatedStartTime,omitempty"` // 更新开始时间
	Paging
}

func (m StockInSheetsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreatedStartTime, validation.When(m.CreatedStartTime != "", validation.Date(constant.DateFormat).Error("创建开始时间格式错误"))),
		validation.Field(&m.CreatedEndTime, validation.When(m.CreatedEndTime != "", validation.Date(constant.DateFormat).Error("创建结束时间格式错误"))),
		validation.Field(&m.UpdatedStartTime, validation.When(m.UpdatedStartTime != "", validation.Date(constant.DateFormat).Error("更新开始时间格式错误"))),
		validation.Field(&m.UpdatedEndTime, validation.When(m.UpdatedEndTime != "", validation.Date(constant.DateFormat).Error("更新结束时间格式错误"))),
	)
}

func (s service) StockInSheets(params StockInSheetsQueryParams) (items []StockInSheet, nextToken string, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.NextToken, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
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
	res := struct {
		tongtool.Response
		Datas struct {
			NextToken string         `json:"nextToken"`
			Result    []StockInSheet `json:"result"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/wmsReceipt/query")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			items = res.Datas.Result
			nextToken = res.Datas.NextToken
			isLastPage = nextToken == ""
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
