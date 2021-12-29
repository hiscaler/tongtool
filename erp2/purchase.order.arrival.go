package erp2

import (
	"encoding/json"
	"errors"
	"github.com/hiscaler/tongtool"
)

type PurchaseOrderArrivalItem struct {
	ArrivalGoodsList  []PurchaseOrderArrivalGoodsItem `json:"arrivalGoodsList"`  // 采购到货明细
	Freight           float64                         `json:"freight"`           // 运费
	PurchaseOrderCode string                          `json:"purchaseOrderCode"` // 采购单号
	Remark            string                          `json:"remark"`            // 到货备注
}

// PurchaseOrderArrivalGoodsItem 采购到货项
type PurchaseOrderArrivalGoodsItem struct {
	GoodsDetailId        string `json:"goodsDetailId"`        // 通途货品ID
	InQuantity           int    `json:"inQuantity"`           // 到货数量
	IsReplace            string `json:"isReplace"`            // 是否是变参替换到货：[Y：是]
	ReplaceGoodsDetailId string `json:"replaceGoodsDetailId"` // 变参替换的通途货品ID
	ReplaceQuantity      int    `json:"replaceQuantity"`      // 变参替换的到货数量
}

type PurchaseOrderArrivalRequest struct {
	MerchantId          string                     `json:"merchantId"`          // 商户ID
	PurchaseArrivalList []PurchaseOrderArrivalItem `json:"purchaseArrivalList"` // 采购到货列表
}

// PurchaseOrderArrival 采购单到货
// https://open.tongtool.com/apiDoc.html#/?docId=ee942453af114a7686d0c8d5187988f2
func (s service) PurchaseOrderArrival(req PurchaseOrderArrivalRequest) (err error) {
	req.MerchantId = s.tongTool.MerchantId
	res := struct {
		result
		Datas interface{} `json:"datas"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(req).
		SetResult(&res).
		Post("/openapi/tongtool/purchaseArrival")
	if err == nil {
		code := 0
		message := ""
		if resp.IsSuccess() {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			if e := json.Unmarshal(resp.Body(), &res); e == nil {
				err = tongtool.ErrorWrap(code, message)
			} else {
				err = errors.New(resp.Status())
			}
		}
	}

	return
}
