package types

import (
	"strconv"
)

// Int64ToString int转换为string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}
