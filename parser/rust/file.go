package rust

import (
	"fmt"
	"regexp"

	"github.com/dlclark/regexp2"
	"github.com/sirupsen/logrus"
)

// TODO: 移除注释掉的函数，如 cfx.rs最后的几个注释掉的函数

func FindStruct(content string, structName string) (Struct, []Use) {
	var re = regexp2.MustCompile(fmt.Sprintf(`\/\/\/(?:.(?!\/\/\/))+pub struct %v {.*?}|pub struct %v {.*?}`, structName, structName), regexp2.Multiline|regexp2.Singleline)
	matched, e := re.FindStringMatch(content)
	if e != nil {
		panic(e)
	}
	if matched == nil {
		logrus.WithFields(logrus.Fields{
			"structName": structName,
			"content:":   content,
		}).Panic("can't find struct")
	}

	s := Struct(matched.String())
	uses := GetUses(content)
	return s, uses
}

func GetUses(content string) []Use {
	reg := regexp.MustCompile(`(?mUs)use .*;`)
	finds := reg.FindAllString(content, -1)
	// fmt.Printf("useFinded %v\n", finds)

	uses := []Use{}
	for _, use := range finds {
		uses = append(uses, Use(use))
	}
	return uses
}

func GetTraits(content string) ([]Trait, []Use) {
	reg := regexp.MustCompile(`(?mUs)(\/\/\/.*\n|)#\[rpc\(.*\)\]\npub trait .* \{[\s\S]*}`)
	finds := reg.FindAllString(string(content), -1)
	// fmt.Printf("traitRegFinded len %v, %v\n", len(finds), finds)

	traits := []Trait{}
	for _, trait := range finds {
		traits = append(traits, Trait(trait))
	}
	return traits, GetUses(content)
}

func GetStructs(content string) (map[string]Struct, []Use) {
	// var re = regexp2.MustCompile(`\/\/\/(?:.(?!\/\/\/))+pub struct ([^\{]*) \{.*?}|pub struct ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	var re = regexp2.MustCompile(`\/\/\/[^{}]+pub struct ([^\{]*) \{.*?}|pub struct ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	ss, uses := getStructsOrEnums(content, re)
	structs := make(map[string]Struct)
	for k, v := range ss {
		structs[k] = Struct(v)
	}
	return structs, uses
}

func GetEnums(content string) (map[string]Enum, []Use) {
	// var re = regexp2.MustCompile(`\/\/\/(?:.(?!\/\/\/))+pub enum ([^\{]*) \{.*?}|pub enum ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	var re = regexp2.MustCompile(`\/\/\/[^{}]+pub enum ([^\{]*) \{.*?}|pub enum ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	es, uses := getStructsOrEnums(content, re)
	enums := make(map[string]Enum)
	for k, v := range es {
		enums[k] = Enum(v)
	}
	return enums, uses
}

func getStructsOrEnums(content string, re *regexp2.Regexp) (map[string]string, []Use) {
	m, e := re.FindStringMatch(content)
	if e != nil {
		panic(e)
	}

	us := GetUses(content)
	if m == nil {
		return nil, us
	}

	structOrEnums := make(map[string]string)

	for m != nil {
		structOrEnums[getStructName(m)] = m.String()
		m, e = re.FindNextMatch(m)
		if e != nil {
			panic(e)
		}
	}

	return structOrEnums, us
}

func getStructName(m *regexp2.Match) string {
	if len(m.Groups()) < 3 {
		panic("can't get struct name")
	}
	name := m.Groups()[1].Capture.String()
	if name == "" {
		name = m.Groups()[2].Capture.String()
	}
	return name
}
