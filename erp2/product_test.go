package erp2

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/tongtool"
	"github.com/hiscaler/tongtool/config"
	"os"
	"testing"
)

var ttInstance *tongtool.TongTool
var ttService Service

func TestMain(m *testing.M) {
	b, err := os.ReadFile("../config/config_test.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var c config.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	ttInstance = tongtool.NewTongTool(c)
	ttService = NewService(ttInstance)
	m.Run()
}

func TestService_Products(t *testing.T) {
	params := ProductsQueryParams{
		ProductType: ProductTypeNormal,
		SKUs:        []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
	}
	products, _, err := ttService.Products(params)
	if err != nil {
		t.Errorf("ttService.Products error: %s", err.Error())
	} else {
		t.Log(jsonx.ToJson(products, "[]"))
	}
}

func TestService_ProductByNormalType(t *testing.T) {
	typ := ProductTypeNormal
	sku := "tt-sku-a"
	isAlias := false
	product, err := ttService.Product(typ, sku, isAlias)
	if err == nil {
		fmt.Println("sku is ", product.SKU)
		fmt.Println(jsonx.ToJson(product, "[]"))
	} else {
		t.Error(err.Error())
	}
	exists, err := ttService.ProductExists(typ, sku, isAlias)
	if !exists {
		t.Errorf("sku %s is not exists, error: %s", sku, err.Error())
	}
}

// 变体商品查询
func TestService_ProductByVariableType(t *testing.T) {
	typ := ProductTypeVariable
	sku := "00145_2"
	isAlias := false
	product, err := ttService.Product(typ, sku, isAlias)
	if err == nil {
		fmt.Println("sku is ", product.SKU)
		fmt.Println(jsonx.ToJson(product, "[]"))
	} else {
		t.Error(err.Error())
	}
	exists, err := ttService.ProductExists(typ, sku, isAlias)
	if !exists {
		t.Errorf("sku %s is not exists, error: %s", sku, err.Error())
	}
}

func TestService_CreateProduct(t *testing.T) {
	req := CreateProductRequest{
		ProductCode:          "tt-sku-c",
		ProductName:          "NETGEAR 路由器",
		ProductPackingEnName: "NETGEAR 4-Stream WiFi 6 Router (R6700AXS) – with 1-Year Armor Cybersecurity Subscription - AX1800 Wireless Speed (Up to 1.8 Gbps) | Coverage up to 1,500 sq. ft., 20+ devices, AX WiFi 6 w/ 1yr Security",
		ProductPackingName:   "NETGEAR 4-Stream WiFi 6 Router (R6700AXS) – with 1-Year Armor Cybersecurity Subscription - AX1800 Wireless Speed (Up to 1.8 Gbps) | Coverage up to 1,500 sq. ft., 20+ devices, AX WiFi 6 w/ 1yr Security",
		DeclareCnName:        "NETGEAR 路由器",
		DeclareEnName:        "NETGEAR 4-Stream WiFi 6 Router (R6700AXS) – with 1-Year Armor Cybersecurity Subscription - AX1800 Wireless Speed (Up to 1.8 Gbps) | Coverage up to 1,500 sq. ft., 20+ devices, AX WiFi 6 w/ 1yr Security",
		HsCode:               "123456",
		ImgUrls: []string{
			"https://m.media-amazon.com/images/I/518c11AD-0L._AC_UY218_.jpg",
		},
		DeveloperName:      "张三",
		PurchaserName:      "李四",
		ProductStatus:      ProductStatusOnSale,
		ProductRemark:      "test",
		SalesType:          ProductSaleTypeNormal,
		ProductCurrentCost: 12,
		ProductWeight:      100,
		CategoryCode:       "未分类",
		ProductLabelIds:    []string{"a", "b"},
		PackageLength:      20,
		PackageWidth:       120,
		PackageHeight:      30,
		EnablePackageNum:   1,
		Accessories: []ProductAccessory{
			{"", 0},
		},
	}
	err := ttService.CreateProduct(req)
	if err == nil {
		fmt.Println("Create product successful.")
	} else {
		t.Errorf("Create product failed, error: %s", err.Error())
	}
}

func TestService_CreateVariableProduct(t *testing.T) {
	req := CreateProductRequest{
		ProductCode:          "tt-sku-a-variable-1",
		ProductName:          "NETGEAR 路由器",
		ProductPackingEnName: "NETGEAR 4-Stream WiFi 6 Router (R6700AXS) – with 1-Year Armor Cybersecurity Subscription - AX1800 Wireless Speed (Up to 1.8 Gbps) | Coverage up to 1,500 sq. ft., 20+ devices, AX WiFi 6 w/ 1yr Security",
		ProductPackingName:   "NETGEAR 4-Stream WiFi 6 Router (R6700AXS) – with 1-Year Armor Cybersecurity Subscription - AX1800 Wireless Speed (Up to 1.8 Gbps) | Coverage up to 1,500 sq. ft., 20+ devices, AX WiFi 6 w/ 1yr Security",
		DeclareCnName:        "NETGEAR 路由器",
		DeclareEnName:        "NETGEAR 4-Stream WiFi 6 Router (R6700AXS) – with 1-Year Armor Cybersecurity Subscription - AX1800 Wireless Speed (Up to 1.8 Gbps) | Coverage up to 1,500 sq. ft., 20+ devices, AX WiFi 6 w/ 1yr Security",
		HsCode:               "123456",
		ImgUrls: []string{
			"https://m.media-amazon.com/images/I/518c11AD-0L._AC_UY218_.jpg",
		},
		DeveloperName:      "张三",
		PurchaserName:      "李四",
		ProductStatus:      ProductStatusOnSale,
		ProductRemark:      "test",
		SalesType:          ProductSaleTypeVariable,
		ProductCurrentCost: 12,
		ProductWeight:      100,
		CategoryCode:       "未分类",
		ProductLabelIds:    []string{"a", "b"},
		PackageLength:      20,
		PackageWidth:       120,
		PackageHeight:      30,
		EnablePackageNum:   1,
		Accessories: []ProductAccessory{
			{"", 0},
		},
		Goods: []ProductGoods{
			{
				GoodsAverageCost: 0,
				GoodsCurrentCost: 0,
				GoodsSKU:         "tt-sku-a-variable-1-1",
				GoodsWeight:      0,
				GoodsVariation:   nil,
			},
		},
	}
	err := ttService.CreateProduct(req)
	if err == nil {
		fmt.Println("Create product successful.")
	} else {
		t.Errorf("Create product failed, error: %s", err.Error())
	}
}

func TestProduct_ImageIsNormalized(t *testing.T) {
	type testCase struct {
		Number     int
		Product    Product
		Normalized bool
	}
	testCases := []testCase{
		{1, Product{}, false},
		{2, Product{
			SKU: "a",
			GoodsDetail: []ProductDetail{
				{GoodsSKU: "a"},
			},
			ProductImgList: []ProductImage{{ImageGroupId: "a.jpg"}},
		}, false},
		{3, Product{
			SKU: "a",
			GoodsDetail: []ProductDetail{
				{GoodsSKU: "a"},
			},
			ProductImgList: []ProductImage{{ImageGroupId: "http://a.com/a.jpg"}},
		}, true},
		{4, Product{
			SKU: "a",
			GoodsDetail: []ProductDetail{
				{GoodsSKU: "a"},
			},
			ProductImgList: []ProductImage{{ImageGroupId: "https://a.com/a.jpg"}},
		}, true},
		{5, Product{
			SKU: "a",
			GoodsDetail: []ProductDetail{
				{GoodsSKU: "a"},
			},
			ProductImgList: []ProductImage{{ImageGroupId: "https://a.com/中文.jpg"}},
		}, false},
		{6, Product{
			SKU: "a",
			GoodsDetail: []ProductDetail{
				{GoodsSKU: "a"},
			},
			ProductImgList: []ProductImage{{ImageGroupId: "https://a.com/a、.jpg"}},
		}, false},
		{7, Product{
			SKU: "a",
			GoodsDetail: []ProductDetail{
				{GoodsSKU: "a"},
			},
			ProductImgList: []ProductImage{{ImageGroupId: "https://a.com/a:jpg"}},
		}, false},
		{8, Product{
			SKU: "a",
			GoodsDetail: []ProductDetail{
				{GoodsSKU: "a"},
			},
			ProductImgList: []ProductImage{{ImageGroupId: "HTTP://a.com/a.jpg"}},
		}, true},
		{9, Product{
			SKU: "a",
			GoodsDetail: []ProductDetail{
				{GoodsSKU: "a"},
			},
			ProductImgList: []ProductImage{{ImageGroupId: "HTTP://a.com/A1.2.b.jpg"}},
		}, true},
	}
	for _, tc := range testCases {
		b := tc.Product.ImageIsNormalized()
		if b != tc.Normalized {
			t.Errorf("%d 期待：%v，实际：%v", tc.Number, tc.Normalized, b)
		}
	}
}
