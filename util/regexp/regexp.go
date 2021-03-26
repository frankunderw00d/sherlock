package regexp

import (
	"regexp"
)

// 正则匹配，str 为正则表达式，v 为值
func Match(str, v string) bool {
	r := regexp.MustCompile(str)
	return r.MatchString(v)
}
