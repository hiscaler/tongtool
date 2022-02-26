package erp2

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/keyx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/constant"
)

// SaleAccount 商户账号信息
type SaleAccount struct {
	SaleAccountId    string   `json:"saleAccountId"`    // 通途账户id
	Account          string   `json:"account"`          // 账户
	AccountCode      string   `json:"accountCode"`      // 账户简码
	PlatformId       string   `json:"platformId"`       // 平台id
	SiteIds          []string `json:"siteIds"`          // 站点id列表
	SiteCountryCodes []string `json:"siteCountryCodes"` // 站点国家代码列表
	Status           string   `json:"status"`           // 账号状态（0：停用、1：启用）
	StatusBoolean    bool     `json:"statusBoolean"`    // 账号状态布尔值
}

type SaleAccountsQueryParams struct {
	Paging
	MerchantId string `json:"merchantId"`            // 商户ID
	PlatformId string `json:"platform_id,omitempty"` // 平台 id
}

func (m SaleAccountsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PlatformId, validation.When(m.PlatformId != "", validation.In(PlatformAmazon, PlatformEBay, PlatformWish, PlatformAliExpress, Platform1688, PlatformShopify, PlatformTopHatter, PlatformFunPinPin, PlatformCouPang, PlatformWalmart, PlatformVOVA, PlatformLazada, PlatformJDInternational, PlatformCdiscount, PlatformNewegg, PlatformRakutenFR, PlatformDHgate, PlatformShopee, PlatformMercadoLibre, PlatformJoom, PlatformMyCom, PlatformFactoryMarket, PlatformYandex, PlatformFXXT, PlatformJDID, PlatformJDGlobal, PlatformTeezily, PlatformAlibabaInternational, PlatformMeesho, PlatformShopline, PlatformJDTH, PlatformAllegro, PlatformBackMarket, PlatformThisShop, PlatformKauflandDE, PlatformMercadoLibreGlobal, PlatformRakutenDE, PlatformXShoppy, PlatformPassfeed, PlatformShopLazza, PlatformDaraz, PlatformTaoBao, PlatformLinio, PlatformB2W, PlatformFunPinPin2, PlatformShopBase, PlatformFordeal, PlatformShoplus, PlatformOnBuy, PlatformManoMano, PlatformShoptima, PlatformShoprises).Error("无效的平台代码"))),
	)
}

// SaleAccounts 商户账号列表
// https://open.tongtool.com/apiDoc.html#/?docId=1e81e4bbae0b4d60b5f7777fc629ba2a
func (s service) SaleAccounts(params SaleAccountsQueryParams) (items []SaleAccount, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.MerchantId = s.tongTool.MerchantId
	params.SetPagingVars(params.PageNo, params.PageSize, s.tongTool.QueryDefaultValues.PageSize)
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
	items = make([]SaleAccount, 0)
	res := struct {
		tongtool.Response
		Datas struct {
			Array    []SaleAccount `json:"array"`
			PageNo   int           `json:"pageNo"`
			PageSize int           `json:"pageSize"`
		} `json:"datas,omitempty"`
	}{}
	resp, err := s.tongTool.Client.R().
		SetBody(params).
		SetResult(&res).
		Post("/openapi/tongtool/merchantSaleAccountQuery")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = tongtool.ErrorWrap(res.Code, res.Message); err == nil {
			for _, item := range res.Datas.Array {
				if params.PlatformId != "" && item.PlatformId != params.PlatformId {
					continue
				}

				item.StatusBoolean = item.Status == "1"
				if item.SiteIds == nil {
					item.SiteIds = []string{}
					item.SiteCountryCodes = []string{}
				} else {
					siteCountryCodes := make([]string, len(item.SiteIds))
					for i, siteId := range item.SiteIds {
						siteCountryCodes[i] = getSiteCountryCodeById(siteId)
					}
					item.SiteCountryCodes = siteCountryCodes
				}
				items = append(items, item)
			}
			isLastPage = len(res.Datas.Array) < params.PageSize
		}
	} else {
		if e := json.Unmarshal(resp.Body(), &res); e == nil {
			err = tongtool.ErrorWrap(res.Code, res.Message)
		} else {
			err = errors.New(resp.Status())
		}
	}
	if err != nil {
		return
	}

	if s.tongTool.EnableCache && len(items) > 0 {
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

// 根据站点 id 获取站点所在国家代码
func getSiteCountryCodeById(siteId string) string {
	if siteId == "" {
		return ""
	}

	code := ""
	switch siteId {
	case "100002":
		code = constant.CountryCodeAmerica
	case "100003":
		code = constant.CountryCodeUnitedKingdom
	case "100004":
		code = constant.CountryCodeCanada
	case "100005":
		code = constant.CountryCodeGermany
	case "100006":
		code = constant.CountryCodeSpain
	case "100007":
		code = constant.CountryCodeFrance
	case "100008":
		code = constant.CountryCodeItaly
	case "100009":
		code = constant.CountryCodeJapan
	case "100010":
		code = constant.CountryCodeMexico
	case "100011":
		code = constant.CountryCodeAustralian
	case "100012":
		code = constant.CountryCodeIndia
	case "100013":
		code = constant.CountryCodeUnitedArabEmirates
	case "100014":
		code = constant.CountryCodeTurkey
	case "100015":
		code = constant.CountryCodeSingapore
	case "100016":
		code = constant.CountryCodeNetherlands
	case "100017":
		code = constant.CountryCodeBrazil
	case "100018":
		code = constant.CountryCodeSaudiArabia
	case "100019":
		code = constant.CountryCodeSweden
	case "100020":
		code = constant.CountryCodePoland
	}
	return code
}
