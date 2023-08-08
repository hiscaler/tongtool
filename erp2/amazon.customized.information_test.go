package erp2

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAmazonCustomizationInformationParser_Parse(t *testing.T) {
	parser := NewAmazonCustomizationInformationParser()
	dirs, err := os.ReadDir("../test/data")
	if err != nil {
		t.Fatalf("os.ReadDir error: %s", err.Error())
	}

	for _, d := range dirs {
		t.Log("DirEntry", d.Type(), d.Name(), d.IsDir())
		if !d.IsDir() {
			parser.Reset()
			// fi, _ := d.Info()
			// t.Logf("file info = %#v", fi)
			// fi.Sys()
			_, err = parser.SetZipFile(filepath.Join("../test/data", d.Name())).Parse()
			if err != nil {
				t.Errorf("%s parse error: %s", d.Name(), err.Error())
			} else {
				t.Logf(`%s
text = %s 
images = %#v`, d.Name(), parser.Text, parser.Images)
			}
		}
	}
}
