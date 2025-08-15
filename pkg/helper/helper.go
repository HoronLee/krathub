package helper

import (
	"regexp"
)

// IsEmail 判断字符串是否为邮箱
func IsEmail(s string) bool {
	re := regexp.MustCompile(`^[\w\.-]+@[\w\.-]+\.\w+$`)
	return re.MatchString(s)
}

// IsPhone 判断字符串是否为手机号（以中国大陆手机号为例）
func IsPhone(s string) bool {
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(s)
}
