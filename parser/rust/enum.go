package rust

import (
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Enum string

type EnumParsed struct {
	Comment string
	Name    string
	Fields  []EnumItemsParsed
}

type EnumItemsParsed struct {
	Comment     string
	TupleParams []TypeParsed
	Value       string
}

func (e Enum) Parse() EnumParsed {
	re := regexp.MustCompile(`(?Us)(.*)pub enum (.*)\{(.*)\}`)
	finds := re.FindStringSubmatch(string(e))

	if finds == nil {
		logrus.WithField("enum Finded", finds).WithField("enum", e).Panic("not enum")
		panic("not struct")
	}

	sComment, sName, sBody := strings.TrimSpace(finds[1]), strings.TrimSpace(finds[2]), strings.TrimSpace(finds[3])

	// fmt.Printf("comment %v\nhead %#v\nbody %#v\n", sComment, sName, sBody)

	iRe := regexp.MustCompile(`(?Um)(\/\/\/.*|^)\s*(\w+|\w+\((.*)\)),`)
	iFinds := iRe.FindAllStringSubmatch(sBody, -1)

	items := make([]EnumItemsParsed, len(iFinds))
	for idx, item := range iFinds {
		// fmt.Printf("item %#v\n", item)
		iComment, iValue := strings.TrimSpace(item[1]), strings.TrimSpace(item[2])

		var tupleParams []TypeParsed
		rawParams := item[3]
		if rawParams != "" {
			sParams := strings.Split(rawParams, ",")
			for _, rawParam := range sParams {
				paramType := RustType(rawParam).Parse()
				tupleParams = append(tupleParams, paramType)
			}
		}

		items[idx] = EnumItemsParsed{
			TupleParams: tupleParams,
			Value:       iValue,
			Comment:     iComment,
		}
	}

	return EnumParsed{sComment, sName, items}
}

func (e *EnumItemsParsed) IsTumple() bool {
	return len(e.TupleParams) > 0
}
