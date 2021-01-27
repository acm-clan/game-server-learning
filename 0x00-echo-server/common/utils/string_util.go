package utils

import "strconv"

func DumpString(s string) string {
	val := ""
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\r':
			val += "\\r"
		case '\n':
			val += "\\n"
		default:
			val += string(s[i])
		}
	}
	return val
}

func DumpStringToByte(s string) string {
	val := ""
	for i := 0; i < len(s); i++ {
		v := int(s[i])
		val += strconv.Itoa(v) + " "
	}
	return val
}
