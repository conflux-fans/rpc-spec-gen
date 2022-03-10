package rust

import (
	"fmt"

	"github.com/dlclark/regexp2"
)

func FindStruct(content string, structName string) Struct {
	var re = regexp2.MustCompile(fmt.Sprintf(`\/\/\/(?:.(?!\/\/\/))+pub struct %v {.*?}`, structName), regexp2.Multiline|regexp2.Singleline)
	matched, e := re.FindStringMatch(content)
	if e != nil {
		panic(e)
	}
	return Struct(matched.String())
}
