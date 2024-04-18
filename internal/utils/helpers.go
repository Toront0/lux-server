package utils

import "unicode"

func ReplaceDoubleQuotesWithSingle(str string) string {

	res := ""


	for _, s := range str {

		if s == '"' {
			res += "'"
		} else {
			res += string(s)
		}

	}

	return res
}

func CamelCaseToSnakeCase(str string) string {
	res := ""

	for _, char := range str {

		if v := unicode.IsUpper(char); v {
			res += "_" + string(unicode.ToLower(char))
		} else {
			res += string(char)
		}

	}


	return res
}