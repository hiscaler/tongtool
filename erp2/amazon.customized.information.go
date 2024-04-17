package erp2

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hiscaler/gox/archivex/zipx"
	"github.com/hiscaler/gox/filex"
	"github.com/hiscaler/gox/inx"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 亚马逊定制信息处理

const (
	imageCustomization  = "ImageCustomization"
	optionCustomization = "OptionCustomization"
	fontCustomization   = "FontCustomization"
	colorCustomization  = "ColorCustomization"
	textCustomization   = "TextCustomization"
)

type PageContainerCustomization struct {
	Type     string                          `json:"type"`
	Children []PreviewContainerCustomization `json:"children"`
}

type PreviewContainerCustomization struct {
	Type      string                       `json:"type"`
	Name      string                       `json:"name"`
	Label     string                       `json:"label"`
	Children  []FlatContainerCustomization `json:"children"`
	BaseImage struct {
		ImageURL string `json:"imageUrl"`
	} `json:"baseImage"`
	Snapshot struct {
		ImageName   string `json:"imageName"`
		ImageBase64 string `json:"imageBase64"`
	}
}

type FlatContainerCustomization struct {
	Type     string                   `json:"type"`
	Children []map[string]interface{} `json:"children"`
}

type Customization struct {
	Type     string        `json:"type"`
	Children []interface{} `json:"children,omitempty"`
}

type ContainerCustomization struct {
	Type     string      `json:"customizationType"`
	Name     string      `json:"name"`
	Label    string      `json:"label"`
	Children interface{} `json:"children"`
}

type ImageCustomization struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Image struct {
		ImageName     string `json:"imageName"`
		BuyerFilename string `json:"buyerFilename"`
	} `json:"image"`
}

// OptionCustomization 下拉
type OptionCustomization struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	DisplayValue string `json:"displayValue"`
}

type FontCustomization struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	Label         string `json:"label"`
	FontSelection struct {
		Family string `json:"family"`
	} `json:"fontSelection"`
}

type ColorCustomization struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	Label          string `json:"label"`
	ColorSelection struct {
		Name       string `json:"name"`
		Value      string `json:"value"`
		ColorModel string `json:"colorModel"`
	} `json:"colorSelection"`
}

type TextCustomization struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Label      string `json:"label"`
	InputValue string `json:"inputValue"`
}

// AmazonCustomizationInformation 亚马逊定制信息
type AmazonCustomizationInformation struct {
	OrderId           string                     `json:"orderId"`
	OrderItemId       string                     `json:"orderItemId"`
	MerchantId        string                     `json:"merchantId"`
	MarketplaceId     string                     `json:"marketplaceId"`
	ASIN              string                     `json:"asin"`
	Title             string                     `json:"title"`
	Quantity          int                        `json:"quantity"`
	CustomizationData PageContainerCustomization `json:"customizationData"`
	Version3          struct {
		CustomizationInfo struct {
			Surfaces []struct {
				Areas []map[string]any `json:"areas"`
			} `json:"surfaces"`
		} `json:"customizationInfo"`
	} `json:"version3.0"`
}

type AmazonCustomizationInformationParser struct {
	zipFile           string
	jsonText          string
	SnapshotImageName string
	SnapshotImage     string            // Image is base64 format
	Images            map[string]string // Image is base64 format
	Text              string
	LabeledValues     map[string]string
	Version3          map[string]string
}

func NewAmazonCustomizationInformationParser() *AmazonCustomizationInformationParser {
	return &AmazonCustomizationInformationParser{}
}

func (parser *AmazonCustomizationInformationParser) SetZipFile(file string) *AmazonCustomizationInformationParser {
	parser.zipFile = file
	return parser
}

func toImageBase64(name string) (base64Encoding string, err error) {
	if !filex.Exists(name) {
		err = fmt.Errorf("%s is invalid name", name)
		return
	}
	b, err := os.ReadFile(name)
	if err != nil {
		return
	}

	switch http.DetectContentType(b) {
	case "image/jpeg":
		base64Encoding = "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding = "data:image/png;base64,"
	default:
		return "", errors.New("unsupported image type")
	}
	base64Encoding += base64.StdEncoding.EncodeToString(b)
	return
}

func isValidCustomizationType(typ string) bool {
	if inx.StringIn(typ, imageCustomization, fontCustomization, colorCustomization, optionCustomization, textCustomization) {
		return true
	}
	return false
}

func read(customizations []map[string]interface{}) (labeledValues map[string]string, lines []string, images []string) {
	labeledValues = make(map[string]string)
	lines = make([]string, 0)
	for _, customization := range customizations {
		s, ok := customization["type"]
		typ := s.(string)
		if !ok {
			return
		}

		if !isValidCustomizationType(typ) {
			rawChildren, ok := customization["children"]
			if !ok {
				continue
			}
			children := make([]map[string]interface{}, 0)
			for _, child := range rawChildren.([]interface{}) {
				children = append(children, child.(map[string]interface{}))
			}
			v1, v2, v3 := read(children)
			for k, v := range v1 {
				labeledValues[k] = v
			}
			lines = append(lines, v2...)
			images = append(images, v3...)
		}

		jsonData, err := json.Marshal(customization)
		if err != nil {
			return
		}

		label := ""
		value := ""
		switch typ {
		case optionCustomization:
			var d OptionCustomization
			err = json.Unmarshal(jsonData, &d)
			if err != nil {
				return
			}
			label = d.Label
			value = d.DisplayValue

		case textCustomization:
			var d TextCustomization
			err = json.Unmarshal(jsonData, &d)
			if err != nil {
				return
			}
			label = d.Label
			value = d.InputValue

		case imageCustomization:
			var d ImageCustomization
			err = json.Unmarshal(jsonData, &d)
			if err != nil {
				return
			}
			label = d.Label
			value = d.Image.ImageName
			if value != "" {
				images = append(images, value)
			}

		case fontCustomization:
			var d FontCustomization
			err = json.Unmarshal(jsonData, &d)
			if err != nil {
				return
			}
			label = d.Label
			value = d.FontSelection.Family

		case colorCustomization:
			var d ColorCustomization
			err = json.Unmarshal(jsonData, &d)
			if err != nil {
				return
			}
			label = d.Label
			value = fmt.Sprintf("%s (%s)", d.ColorSelection.Name, d.ColorSelection.Value)
		}

		if label != "" {
			lines = append(lines, fmt.Sprintf("%s:%s", label, value))
			labeledValues[label] = value
		}
	}

	return
}

// Parse 解析压缩文件内容
// 文件名格式：702-2644781-4722617_82729974619961 [订单号_订单商品 ID]，json 文件为 82729974619961
func (parser *AmazonCustomizationInformationParser) Parse() (*AmazonCustomizationInformationParser, error) {
	if parser.zipFile == "" {
		return parser, errors.New("zip 文件路径不能为空")
	}

	if !filex.Exists(parser.zipFile) {
		return parser, fmt.Errorf("%s 不存在", parser.zipFile)
	}
	zipFilename := strings.Replace(filepath.Base(parser.zipFile), ".zip", "", -1)
	dst := filepath.Join(os.TempDir(), zipFilename)
	err := zipx.UnCompress(parser.zipFile, dst)
	if err != nil {
		return parser, err
	}
	defer os.RemoveAll(dst)

	jsonFile := ""
	index := strings.LastIndex(zipFilename, "_")
	if index == -1 {
		return parser, errors.New("无法获取到 JSON 文件名。")
	}

	jsonFile = filepath.Join(dst, zipFilename[index+1:]+".json")
	if !filex.Exists(jsonFile) {
		return parser, fmt.Errorf("%s 不存在。", jsonFile)
	}
	b, err := os.ReadFile(jsonFile)
	if err != nil {
		return parser, err
	}
	parser.jsonText = string(b)

	var ci AmazonCustomizationInformation
	err = json.Unmarshal(b, &ci)
	if err != nil {
		return parser, err
	}

	if len(ci.CustomizationData.Children) != 1 {
		return parser, errors.New("无效的 JSON")
	}

	if len(ci.Version3.CustomizationInfo.Surfaces) != 0 {
		labeledValue := make(map[string]string)
		areas := ci.Version3.CustomizationInfo.Surfaces[0].Areas
		for _, area := range areas {
			typ, ok := area["customizationType"]
			if !ok {
				continue
			}
			label := ""
			if v, ok := area["label"]; ok {
				label = v.(string)
			}
			value := ""
			switch typ {
			case "TextPrinting":
				if v, ok := area["text"]; ok {
					value = v.(string)
				}

			case "Options":
				if v, ok := area["optionValue"]; ok {
					value = v.(string)
				}
			}
			if label == "" {
				continue
			}
			labeledValue[label] = value
		}
		parser.Version3 = labeledValue
	}

	previewContainerCustomizationData := ci.CustomizationData.Children[0]

	var imageBase64String string

	imageName := previewContainerCustomizationData.Snapshot.ImageName
	if imageName != "" {
		parser.SnapshotImageName = imageName
		imageBase64String, err = toImageBase64(filepath.Join(dst, imageName))
		if err != nil {
			return parser, err
		}
		parser.SnapshotImage = imageBase64String
	}
	lines := make([]string, 0)
	images := make(map[string]string)
	labeledValues := make(map[string]string)
	for _, c := range previewContainerCustomizationData.Children {
		v1, v2, v3 := read(c.Children)
		for k, v := range v1 {
			labeledValues[k] = v
		}
		lines = append(lines, v2...)
		for _, img := range v3 {
			imageBase64String, err = toImageBase64(filepath.Join(dst, img))
			if err != nil {
				return parser, err
				// images[img] = img
			} else {
				images[img] = imageBase64String
			}
		}
	}

	parser.Text = strings.Join(lines, "\n")
	parser.Images = images
	parser.LabeledValues = labeledValues

	return parser, nil
}

func (parser *AmazonCustomizationInformationParser) Reset() *AmazonCustomizationInformationParser {
	parser.zipFile = ""
	parser.jsonText = ""
	parser.SnapshotImageName = ""
	parser.SnapshotImage = ""
	parser.Images = make(map[string]string)
	parser.Text = ""
	parser.LabeledValues = make(map[string]string)
	parser.Version3 = make(map[string]string)
	return parser
}
