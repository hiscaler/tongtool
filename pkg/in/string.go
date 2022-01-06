package in

import "strings"

// StringIn 判断 s 是否在 ss 中（忽略大小写）
func StringIn(s string, ss ...string) bool {
	if len(ss) == 0 {
		return false
	}
	for _, s2 := range ss {
		if strings.EqualFold(s, s2) {
			return true
		}
	}
	return false
}
