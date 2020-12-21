package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// StringToInt 字符串转数字
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		logger.LogError(err)
	}
	return i
}
