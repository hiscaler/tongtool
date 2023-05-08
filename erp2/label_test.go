package erp2

import (
	"fmt"
	"testing"
)

func TestService_Labels(t *testing.T) {
	params := LabelsQueryParams{}
	params.PageNo = 1
	_, _, err := ttService.Labels(params)
	if err != nil {
		t.Errorf("ttService.Labels error: %s", err.Error())
	}
}

func TestService_CreateLabel(t *testing.T) {
	req := CreateLabelRequest{
		LabelName: "test1",
	}
	err := ttService.CreateLabel(req)
	if err == nil {
		fmt.Println("Create label successful.")
	} else {
		t.Errorf("Create label failed, error: %s", err.Error())
	}
}

// Label 判断是否存在
func TestService_LabelExists(t *testing.T) {
	exists, err := ttService.LabelExists("test1")
	if !exists || err != nil {
		t.Errorf("test label is not exists")
	}
}
