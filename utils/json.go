package utils

import "encoding/json"

func MustJsonPretty(v interface{}) string {
	j, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		panic(e)
	}
	return string(j)
}
