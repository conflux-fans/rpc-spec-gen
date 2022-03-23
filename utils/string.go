package utils

import (
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

func CamelCase2UnderScoreCase(str string) string {
	str = regexp.MustCompile(`([A-Z])`).ReplaceAllString(str, "_$1")
	str = strings.ToLower(str)
	str = strings.Trim(str, "_")
	return str
}

func UnderScoreCase2CamelCase(str string, firstUpper bool) string {

	var fun regexp2.MatchEvaluator = func(m regexp2.Match) string {
		headUpper := strings.ToUpper(m.Groups()[1].String())
		return headUpper
	}

	str, e := regexp2.MustCompile(`_([a-z])`, regexp2.None).ReplaceFunc(str, fun, -1, -1)
	if e != nil {
		panic(e)
	}

	if firstUpper {
		return strings.ToUpper(str[:1]) + str[1:]
	}
	return strings.ToLower(str[:1]) + str[1:]
}
