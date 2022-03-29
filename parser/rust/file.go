package rust

import (
	"fmt"
	"regexp"

	"github.com/dlclark/regexp2"
	"github.com/sirupsen/logrus"
)

// TODO: 移除注释掉的函数，如 cfx.rs最后的几个注释掉的函数
// TODO: 移除掉mod test (/Users/wangdayong/myspace/mywork/conflux-rust/client/src/rpc/types/eth/sync.rs)

type SourceCode string

func (sc SourceCode) FindStruct(structName string) (Struct, []Use) {
	content := string(sc)
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
	uses := sc.GetUses()
	return s, uses
}

func (sc SourceCode) GetUses() []Use {
	content := string(sc)
	reg := regexp.MustCompile(`(?mUs)use .*;`)
	finds := reg.FindAllString(content, -1)
	// fmt.Printf("useFinded %v\n", finds)

	uses := []Use{}
	for _, use := range finds {
		uses = append(uses, Use(use))
	}
	return uses
}

func (sc SourceCode) GetTraits() ([]Trait, []Use) {
	content := string(sc)
	reg := regexp.MustCompile(`(?mUs)(\/\/\/.*\n|)#\[rpc\(.*\)\]\npub trait .* \{[\s\S]*}`)
	finds := reg.FindAllString(string(content), -1)
	// fmt.Printf("traitRegFinded len %v, %v\n", len(finds), finds)

	traits := []Trait{}
	for _, trait := range finds {
		traits = append(traits, Trait(trait))
	}
	return traits, sc.GetUses()
}

func (sc SourceCode) GetStructs() (map[string]Struct, []Use) {
	// content := string(sc)
	// var re = regexp2.MustCompile(`\/\/\/(?:.(?!\/\/\/))+pub struct ([^\{]*) \{.*?}|pub struct ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	var re = regexp2.MustCompile(`\/\/\/[^{}]+pub struct ([^\{]*) \{.*?}|pub struct ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	ss, uses := sc.getStructsOrEnums(re)
	structs := make(map[string]Struct)
	for k, v := range ss {
		structs[k] = Struct(v)
	}
	return structs, uses
}

func (sc SourceCode) GetEnums() (map[string]Enum, []Use) {
	// content := string(sc)
	// var re = regexp2.MustCompile(`\/\/\/(?:.(?!\/\/\/))+pub enum ([^\{]*) \{.*?}|pub enum ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	var re = regexp2.MustCompile(`\/\/\/[^{}]+pub enum ([^\{]*) \{.*?}|pub enum ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	es, uses := sc.getStructsOrEnums(re)
	enums := make(map[string]Enum)
	for k, v := range es {
		enums[k] = Enum(v)
	}
	return enums, uses
}

func (sc SourceCode) GetDefineTypes() (map[string]RustType, []Use) {
	content := string(sc)
	var re = regexp2.MustCompile(`^\/\/\/[^;]*?pub type (.*?) = (.*?);|pub type (.*?) = (.*?);`, regexp2.Multiline|regexp2.Singleline)
	m, e := re.FindStringMatch(content)
	if e != nil {
		panic(e)
	}

	us := sc.GetUses()
	if m == nil {
		return nil, us
	}

	defineTypes := make(map[string]RustType)
	for m != nil {
		name, define := getDefineType(m)
		defineTypes[name] = RustType(define)
		m, e = re.FindNextMatch(m)
		if e != nil {
			panic(e)
		}
	}

	return defineTypes, us
}

func (sc SourceCode) getStructsOrEnums(re *regexp2.Regexp) (map[string]string, []Use) {
	content := string(sc)
	m, e := re.FindStringMatch(content)
	if e != nil {
		panic(e)
	}

	us := sc.GetUses()
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

func getDefineType(m *regexp2.Match) (string, string) {
	if len(m.Groups()) < 5 {
		panic("can't get struct name")
	}
	name := m.Groups()[1].Capture.String()
	if name != "" {
		define := m.Groups()[2].Capture.String()
		return name, define
	}
	name = m.Groups()[3].Capture.String()
	if name != "" {
		define := m.Groups()[2].Capture.String()
		return name, define
	}
	panic("can't get struct name and type define")
}
