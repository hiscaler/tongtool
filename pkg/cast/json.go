package cast

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func ToJson(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("%+v", i)
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", i)
	}

	return buf.String()
}
