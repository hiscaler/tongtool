package erp2

import (
	"fmt"
	"github.com/hiscaler/tongtool/pkg/cast"
	"testing"
)

func TestService_Products(t *testing.T) {
	_, ttService := newTestTongTool()
	params := ProductQueryParams{
		ProductType: ProductTypeNormal,
	}
	_, _, err := ttService.Products(params)
	if err != nil {
		t.Errorf("ttService.Products error: %s", err.Error())
	}
}

func TestService_ProductByNormalType(t *testing.T) {
	_, ttService := newTestTongTool()
	typ := ProductTypeNormal
	sku := "Lillian201309130002"
	isAlias := false
	product, err := ttService.Product(typ, sku, isAlias)
	if err == nil {
		fmt.Println("sku is ", product.SKU)
		fmt.Println(cast.ToJson(product))
	} else {
		t.Error(err.Error())
	}
	exists := ttService.ProductExists(typ, sku, isAlias)
	if !exists {
		t.Errorf("sku %s is not exists.", sku)
	}
}

// 变体商品查询
func TestService_ProductByVariableType(t *testing.T) {
	_, ttService := newTestTongTool()
	typ := ProductTypeVariable
	sku := "00145_2"
	isAlias := false
	product, err := ttService.Product(typ, sku, isAlias)
	if err == nil {
		fmt.Println("sku is ", product.SKU)
		fmt.Println(cast.ToJson(product))
	} else {
		t.Error(err.Error())
	}
	exists := ttService.ProductExists(typ, sku, isAlias)
	if !exists {
		t.Errorf("sku %s is not exists.", sku)
	}
}

func TestService_CreateProduct(t *testing.T) {
	ttInstance, ttService := newTestTongTool()
	labels := make([]Label, 0)
	labels = append(labels, Label{SKULabel: "tt-sku-a"})
	req := CreateProductRequest{
		ProductCode:          "tt-sku-a",
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
		ProductLabelIds:    []string{},
		PackageLength:      20,
		PackageWidth:       120,
		PackageHeight:      30,
	}
	req.MerchantId = ttInstance.MerchantId
	err := ttService.CreateProduct(req)
	if err == nil {
		fmt.Println("Create product successful.")
	} else {
		t.Errorf("Create product failed, error: %s", err.Error())
	}
}
