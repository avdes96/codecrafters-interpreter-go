package utils

import "strconv"

func Stringify(val any) string {
	if val == nil {
		return "nil"
	} else if v, ok := val.(bool); ok {
		return strconv.FormatBool(v)
	} else if v, ok := val.(float64); ok {
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return val.(string)
}
