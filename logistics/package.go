package logistics

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
	"strconv"
	"strings"
)

// 获取包裹信息
// https://open.tongtool.com/apiDoc.html#/?docId=68251fa7414b43b5a1b0eddc898f6319

// 物流商授权信息
type apiParam struct {
	Key   string `json:"apiKey"`   // 物流商授权代码
	Name  string `json:"apiName"`  // 物流商授权名称
	Value string `json:"apiValue"` // 物流商授权值
}

// 申报信息
type declaration struct {
	DeclareCnName    string  `json:"declareCnName"`    // 申报中文品名
	DeclareCurrency  string  `json:"declareCurrency"`  // 申报币种，三字货币代码
	DeclareEnName    string  `json:"declareEnName"`    // 申报英文品名
	DeclareNumber    int     `json:"declareNumber"`    // 申报产品数量
	DeclareProductId string  `json:"declareProductId"` // 申报的产品ID,通常为平台上的产品id
	DeclareURL       string  `json:"declareUrl"`       // 申报的产品链接
	DeclareValue     float64 `json:"declareValue"`     // 单个货品申报价值
	DeclareWeight    float64 `json:"declareWeight"`    // 单个货品申报重量，单位克
	GoodsSKU         string  `json:"goodsSku"`         // 申报产品SKU
	HsCode           string  `json:"hsCode"`           // 海关编码
	Material         string  `json:"material"`         // 申报的产品材质
	Purpose          string  `json:"purpose"`          // 申报的产品用途
}

// 邮寄方式扩展设置信息
type extendParameter struct {
	Code  string `json:"parameterCode"`  // 扩展参数代码
	Name  string `json:"parameterName"`  // 扩展参数名称
	Type  string `json:"parameterType"`  // input 输入字符串,checkbox 多选按钮,select单选按钮
	Value string `json:"parameterValue"` // 如果是多选或者单选按钮，需要提供取值范围，多个取值之间用逗号或者分号隔开。
}

// 配货信息
type picking struct {
	CargoSpace        string `json:"cargoSpace"`        // 货位号
	EbayItemId        string `json:"ebayItemId"`        //	Ebay订单货品ID
	EbayTransactionId string `json:"ebayTransactionId"` //	Ebay订单交易ID
	ProductName       string `json:"productName"`       //	货品配货名称
	ProductWeight     int    `json:"productWeight"`     //	货品单品重量（单位克）
	Quantity          int    `json:"quantity"`          //	数量
	Remark            string `json:"remark"`            //	备注
	SKU               string `json:"sku"`               //	货品SKU（默认先传仓库专用货号）
	WarehouseName     string `json:"warehouseName"`     //	仓库代码
}

type Package struct {
	APIParamArray          []apiParam        `json:"apiParamArray"`          // 物流商授权信息
	BuyerPassportCode      string            `json:"buyerPassportCode"`      // 收件人识别编号
	CarrierOrderId         string            `json:"carrierOrderId"`         // 物流商系统单号未在物流商系统下单状态的订单为空
	DeclarationArray       []declaration     `json:"declarationArray"`       // 申报信息列表
	EbayBuyerId            string            `json:"ebayBuyerId"`            // Ebay买家ID
	EbaySellerId           string            `json:"ebaySellerId"`           // Ebay卖家ID
	ExtendParameterArray   []extendParameter `json:"extendParameterArray"`   // 邮寄方式扩展设置信息
	Height                 float64           `json:"height"`                 //	包裹高，单位cm
	IOSSMethod             string            `json:"iossMethod"`             //	预缴增值税方式（IOSS、no-IOSS、other）
	IOSSNo                 string            `json:"iossNo"`                 //	卖家 IOSS 税号
	LastSyncTime           string            `json:"lastsyncTime"`           //	订单状态的最后更新时间
	Length                 float64           `json:"length"`                 //	包裹长（单位cm）
	MerchantId             string            `json:"merchantId"`             //	通途商户号
	PickingArray           []picking         `json:"pickingArray"`           //	配货信息
	PlatformId             string            `json:"platformId"`             // 平台类型
	RecipientAddress1      string            `json:"recipientAddress1"`      //	收件人地址1
	RecipientAddress2      string            `json:"recipientAddress2"`      //	收件人地址2
	RecipientCity          string            `json:"recipientCity"`          //	收件人城市
	RecipientCompany       string            `json:"recipientCompany"`       //	收件人公司
	RecipientCountry       string            `json:"recipientCountry"`       //	收件人国家二字代码
	RecipientCountryCnName string            `json:"recipientCountryCnName"` //	收件人国家中文名称
	RecipientCountryEnName string            `json:"recipientCountryEnName"` //	收件人国家英文名称
	RecipientEmail         string            `json:"recipientEmail"`         //	收件人电子邮箱
	RecipientMobile        string            `json:"recipientMobile"`        //	收件人手机
	RecipientName          string            `json:"recipientName"`          //	收件人姓名
	RecipientPostalCode    string            `json:"recipientPostalCode"`    //	收件人邮编
	RecipientState         string            `json:"recipientState"`         //	收件人省州
	RecipientTelephone     string            `json:"recipientTelephone"`     //	收件人电话
	SalesRecordNumber      string            `json:"salesRecordNumber"`      // 包裹订单号（如果有多订单情况会以|做区分）
	SenderAddress1         string            `json:"senderAddress1"`         //	寄件人地址1
	SenderAddress2         string            `json:"senderAddress2"`         //	寄件人地址2
	SenderCity             string            `json:"senderCity"`             //	寄件人城市
	SenderCompany          string            `json:"senderCompany"`          //	寄件人公司
	SenderCountry          string            `json:"senderCountry"`          //	寄件人国家
	SenderEmail            string            `json:"senderEmail"`            //	寄件人电子邮箱
	SenderMobile           string            `json:"senderMobile"`           //	寄件人手机
	SenderName             string            `json:"senderName"`             //	寄件人姓名
	SenderPostalCode       string            `json:"senderPostalCode"`       //	寄件人邮编
	SenderState            string            `json:"senderState"`            //	寄件人省州
	SenderTelephone        string            `json:"senderTelephone"`        //	寄件人电话
	ShippingMethodCode     string            `json:"shippingMethodCode"`     //	物流渠道代码
	TrackingNumber         string            `json:"trackingNumber"`         //	未在物流商系统下单状态的订单或无跟踪号的渠道的订单为空
	TtPacketId             string            `json:"ttPacketId"`             //	通途包裹号
	TtPacketStatus         string            `json:"ttPacketStatus"`         //	通途包裹状态WAIT_UPLOAD 等待在物流商系统下单,WAIT_CONFIRM 等待在物流商系统交运,CONFIRM客户已经交运但是没有发货,WAIT_CANCEL等待在物流商系统取消,FAILURE物流商系统处理失败
	VatNo                  string            `json:"vatNo"`                  //	卖家VAT税号
	VirtualTrackingNumber  string            `json:"virtualTrackingNumber"`  //	虚拟跟踪号
	Width                  float64           `json:"width"`                  //	包裹宽，单位cm
}

type PackagesQueryParams struct {
	Paging
	MerchantId         string `json:"merchantId"`                   // 商戶号
	OrderStatus        string `json:"orderStatus,omitempty"`        // 通途订单状态
	ShippingMethodCode string `json:"shippingMethodCode,omitempty"` // 渠道代码
	Since              string `json:"since"`                        // 查询的起始时间
}

func (m PackagesQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Since,
			validation.Required.Error("起始时间不能为空"),
			validation.Date(constant.DatetimeFormat).Error("起始时间格式错误"),
		),
	)
}

func (s service) Packages(params PackagesQueryParams) (items []Package, nextToken string, isLastPage bool, err error) {
	params.MerchantId = s.tongTool.MerchantId
	if params.PageSize <= 0 || params.PageSize > s.tongTool.QueryDefaultValues.PageSize {
		params.PageSize = s.tongTool.QueryDefaultValues.PageSize
	}
	if err != nil {
		return
	}
	res := struct {
		tongtool.Response
		NextToken string `json:"nextToken"`
		Datas     struct {
			ACK          string    `json:"ack"`          // 响应结果（Success：成功、Failure：失败）
			ErrorCode    string    `json:"errorCode"`    // 错误代码
			ErrorMessage string    `json:"errorMessage"` // 错误信息
			NextToken    string    `json:"nextToken"`    // 是否有下一页，有下一页返回下一页的token
			OrderArray   []Package `json:"orderArray"`   // 订单列表（调用失败或查询无结果是为空）
		} `json:"datas"`
		PageNo   int `json:"pageNo"`
		PageSize int `json:"pageSize"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/logi/getOrder")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			if strings.EqualFold(res.Datas.ACK, "Success") {
				items = res.Datas.OrderArray
				nextToken = res.Datas.NextToken
				isLastPage = len(items) <= params.PageSize
			} else {
				errorCode, _ := strconv.Atoi(res.Datas.ErrorCode)
				err = tongtool.ErrorWrap(errorCode, res.Datas.ErrorMessage)
			}
		}
	} else {
		if e := json.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	return
}
