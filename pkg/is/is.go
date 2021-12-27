package is

import (
	"strings"
)

// Number 判断是否为数字
func Number(s string) bool {
	isNotDigit := func(c rune) bool { return c < '0' || c > '9' }
	return strings.IndexFunc(s, isNotDigit) == -1
}
