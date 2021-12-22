package in

import "strings"

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
