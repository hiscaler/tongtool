package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	"github.com/hiscaler/tongtool/pkg/cache"
)

// 执行状态
const (
	AfterSaleExecutionStatusRefunded           = "refunded"           // 已退款
	AfterSaleExecutionStatusRefundFail         = "refundFail"         // 退款失败
	AfterSaleExecutionStatusDepatchOneMore     = "depatchOneMore"     // 已补发货
	AfterSaleExecutionStatusRGoodStockInCancel = "rGoodStockInCancel" // 取消退款入库

)

// AfterSaleItem 售后项
type AfterSaleItem struct {
	GoodsSKU    string `json:"goods_sku"`     // 货品sku
	ProductName string `json:"producrt_name"` // 货品名称
	Quantity    int    `json:"quantity"`      // 数量
}

type AfterSaleService struct {
	Refunded      bool `json:"refunded"`       // 退款
	ReturnedGoods bool `json:"returned_goods"` // 退货
	ReissueGoods  bool `json:"reissue_goods"`  // 补发
}

// AfterSale 售后单
type AfterSale struct {
	AfterSaleServiceType   string          `json:"after_sale_service_type"`   // 退款退货补发类型:以三位数表示 首位表示是否退款，第二位表示是否退货，第三位表示是否补发 0--表示不执行对应操作，1--表示执行以上操作 如：111表示退款退货补发，001表示补发
	ApproveStatus          string          `json:"approve_status"`            // 审批状态 0或null-未提交 1-已提交 2-审批退回 3-审批通过
	BuyerAccountId         string          `json:"buyer_account_id"`          // 买家ID
	BuyerCountryCode       string          `json:"buyer_country_code"`        // 买家国家编码
	BuyerReturnTrackingNum string          `json:"buyer_return_tracking_num"` // 买家退货的跟踪单号
	CancelReturnType       string          `json:"cancel_return_type"`        // 是否取消退货(1：取消退货)
	CarrierId              string          `json:"carrier_id"`                // 承运人顺序号
	CreateDate             string          `json:"create_date"`               // 售后单创建时间
	CreatedBy              string          `json:"created_by"`                // 创建人
	GoodsList              []AfterSaleItem `json:"goodsList"`                 // 货品信息
	OldBuyerCountryCode    string          `json:"old_buyer_country_code"`    // old买家国家编码
	OrderId                string          `json:"order_id"`                  // 通途销售单内部单号
	OrderOwner             string          `json:"order_owner"`               // 订单来源
	RefundAmount           float64         `json:"refund_amount"`             // 退款金额
	RefundCurrency         string          `json:"refund_currency"`           // 审批状态 0或null-未提交 1-已提交 2-审批退回 3-审批通过
	RefundStatus           string          `json:"refund_status"`             // 退款状态 0或null--未退款 1--退款失败 2--退款完成
	SalePlatform           string          `json:"sale_platform"`             // 平台编码
	SalesRecordNumber      string          `json:"sales_record_number"`       // 订单号
	SellerAccountId        string          `json:"seller_account_id"`         // 卖家ID
	ShippingMethod         string          `json:"shipping_method"`           // 邮寄方式顺序号
	SubmittedBy            string          `json:"submitted_by"`              // 提交人
	SubmittedDate          string          `json:"submitted_date"`            // 售后单提交时间
	WarehouseStoreType     string          `json:"warehouse_store_type"`      // 仓库退货补发状态
	// 扩展属性
	AfterSaleService AfterSaleService `json:"after_sale_service"` // 退款退货补发
}

type AfterSaleQueryParams struct {
	ApproveStatus   string `json:"approveStatus"`   // 审核状态：approveStatus格式错误 、applying/等待提交、approving/等待审批、approved/审批通过
	CreatedDateFrom string `json:"createdDateFrom"` // 售后创建开始时间
	CreatedDateTo   string `json:"createdDateTo"`   // 售后创建结束时间
	MerchantId      string `json:"merchantId"`
	OrderId         string `json:"orderId"`            // 订单ID
	PageNo          int    `json:"pageNo,omitempty"`   // 查询页数
	PageSize        int    `json:"pageSize,omitempty"` // 每页数量,默认值：100,最大值100，超过最大值以最大值数量返回
	ZhiXingStatus   string `json:"zhixingStatus"`      // 执行状态: refunded/已退款、refundFail/退款失败、depatchOneMore/已补发货、rGoodStockInCancel/取消退款入库
}

func (m AfterSaleQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreatedDateFrom, validation.Date(constant.DatetimeFormat).Error("售后创建开始时间格式无效")),
		validation.Field(&m.CreatedDateTo, validation.Date(constant.DatetimeFormat).Error("售后创建结束时间格式无效")),
		validation.Field(&m.ZhiXingStatus, validation.In(AfterSaleExecutionStatusRefunded, AfterSaleExecutionStatusRefundFail, AfterSaleExecutionStatusDepatchOneMore, AfterSaleExecutionStatusRGoodStockInCancel)),
	)
}

// AfterSales 查询售后单信息
// https://open.tongtool.com/apiDoc.html#/?docId=4406fd10e6d34dae994043e6d52c4e33
func (s service) AfterSales(params AfterSaleQueryParams) (items []AfterSale, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}
	params.MerchantId = s.tongTool.MerchantId
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	var cacheKey string
	if s.tongTool.EnableCache {
		cacheKey = cache.GenerateKey(params)
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
	items = make([]AfterSale, 0)
	res := struct {
		result
		Datas struct {
			Array    []AfterSale `json:"array"`
			PageNo   int         `json:"pageNo"`
			PageSize int         `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/afterSalesQuery")
	if err == nil {
		if resp.IsSuccess() {
			if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
				items = res.Datas.Array
				for i, item := range items {
					ass := AfterSaleService{}
					if item.AfterSaleServiceType[0:1] == "1" {
						ass.Refunded = true
					}
					if item.AfterSaleServiceType[1:2] == "1" {
						ass.ReturnedGoods = true
					}
					if item.AfterSaleServiceType[2:3] == "1" {
						ass.ReissueGoods = true
					}
					items[i].AfterSaleService = ass
				}
				isLastPage = len(items) < params.PageSize
			}
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(res.Code, res.Message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	if err == nil && s.tongTool.EnableCache {
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
