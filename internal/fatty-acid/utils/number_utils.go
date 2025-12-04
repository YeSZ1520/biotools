package utils

import (
	"strconv"
)

// isInteger 判断字符串是否表示一个整数
func IsInteger(s string) bool {
    _, err := strconv.Atoi(s)
    return err == nil
}