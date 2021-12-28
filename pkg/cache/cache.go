package cache

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// GenerateKey 生成 Cache 键名
func GenerateKey(values ...interface{}) string {
	var sb strings.Builder
	for _, value := range values {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String:
			sb.WriteString(value.(string))
		case reflect.Bool:
			sb.WriteString(strconv.FormatBool(v.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			sb.WriteString(fmt.Sprintf("%d", value))
		case reflect.Float32, reflect.Float64:
			sb.WriteString(strconv.FormatFloat(v.Float(), 'f', 2, 64))
		case reflect.Map:
			mapValue := value.(map[string]interface{})
			keys := make([]string, len(mapValue))
			i := 0
			for k, _ := range mapValue {
				keys[i] = k
				i++
			}
			sort.Strings(keys)
			for _, k := range keys {
				sb.WriteString(GenerateKey(k, mapValue[k]))
			}
		case reflect.Slice, reflect.Array:
			interfaces := make([]interface{}, 0)
			for i := 0; i < v.Len(); i++ {
				interfaces = append(interfaces, v.Index(i).Interface())
			}
			sb.WriteString(GenerateKey(interfaces...))
		case reflect.Struct:
			kv := map[string]interface{}{}
			t := reflect.TypeOf(value)
			for k := 0; k < t.NumField(); k++ {
				kv[t.Field(k).Name] = v.Field(k).Interface()
			}
			sb.WriteString(GenerateKey(kv))
		default:
			sb.WriteString(v.String())
		}
	}
	return sb.String()
}
