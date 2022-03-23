package rust

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Struct string

type RustType string

// StructParsed reperesnts parse result of a Rust struct
type StructParsed struct {
	Comment string
	Name    string
	Fields  []FieldParsed
}

// FieldParsed reperesnts parse result of a Rust struct field
type FieldParsed struct {
	Comment string
	Name    string
	Type    TypeParsed
}

// TypeParsed reperesnts parse result of a Rust struct field type
type TypeParsed struct {
	IsOption        bool
	IsArray         bool
	IsBoxFuture     bool
	IsVariadicValue bool
	Name            string
	Core            *TypeParsed
}

func (t *TypeParsed) InnestCoreTypeName() string {
	// for {
	if t.Core == nil {
		return t.Name
	}
	//  else {
	// t = t.Core
	return t.Core.InnestCoreTypeName()
	// }
	// }
}

func (r Struct) Parse() StructParsed {
	structReg := regexp.MustCompile(`(?Us)(.*)pub struct (.*)\{(.*)\}`)
	structFinded := structReg.FindStringSubmatch(string(r))
	// fmt.Printf("%#v\n", structFinded)

	if structFinded == nil {
		logrus.WithField("structFinded", structFinded).WithField("struct", r).Info("structFinded")
		panic("not struct")
	}

	sComment, sName, sBody := strings.TrimSpace(structFinded[1]), strings.TrimSpace(structFinded[2]), strings.TrimSpace(structFinded[3])

	fmt.Printf("comment %v\nhead %#v\nbody %#v\n", sComment, sName, sBody)

	// TODO： 解析LedgerInfoWithV0不正确
	fieldReg := regexp.MustCompile(`(?Us)(.*)pub (.*): (.*),`)
	// fieldReg := regexp.MustCompile(`(?Us)(.*)(pub|)\s*(.*):\s*(.*),`)
	fieldsFinded := fieldReg.FindAllStringSubmatch(sBody, -1)
	// fmt.Printf("fieldsFinded %#v\n", fieldsFinded[0])

	fields := make([]FieldParsed, len(fieldsFinded))
	for i, field := range fieldsFinded {
		fmt.Printf("field %#v\n", field)
		fComment, fName, fType := strings.TrimSpace(field[1]), strings.TrimSpace(field[2]), RustType(field[3])
		fields[i] = FieldParsed{fComment, fName, fType.Parse()}
		logger.WithFields(logrus.Fields{
			"field":        field[0],
			"field parsed": fields[i],
		}).Debug("field parsed")
	}

	return StructParsed{sComment, sName, fields}
}

func (r RustType) Parse() (result TypeParsed) {

	inner, e := re_AngleBrackets.FindStringMatch(string(r))
	if e != nil {
		logger.WithField("rustType", r).WithError(e).Panic("failed to found inner type")
	}
	// logger.WithField("rustType", r).WithField("inner", inner).Debug("find inner")

	if inner == nil {
		result.Name = string(r)
		return
	}

	result.Name = string(r)

	trimedInner := strings.TrimPrefix(inner.String(), "<")
	trimedInner = strings.TrimSuffix(trimedInner, ">")
	coreParsed := RustType(trimedInner).Parse()
	result.Core = &coreParsed

	genericType := strings.TrimSuffix(string(r), inner.String())
	switch genericType {
	case "Option":
		result.IsOption = true
	case "Vec":
		result.IsArray = true
	case "BoxFuture":
		result.IsBoxFuture = true
	case "VariadicValue":
		result.IsVariadicValue = true
	}

	return

	// optionReg := regexp.MustCompile(`Option<(.*)>`)
	// optionMatched := optionReg.FindStringSubmatch(string(r))

	// // fmt.Printf("optionMatched %#v\n", optionMatched)
	// if len(optionMatched) > 0 {
	// 	result.IsOption = true
	// 	result.Name = optionMatched[1]
	// 	coreParsed := RustType(result.Name).Parse()
	// 	result.Core = &coreParsed
	// 	return result
	// }

	// vecReg := regexp.MustCompile(`Vec<(.*)>`)
	// vecMatched := vecReg.FindStringSubmatch(string(r))
	// // fmt.Printf("vecMatched %#v\n", vecMatched)
	// if len(vecMatched) > 0 {
	// 	result.IsArray = true
	// 	result.Name = vecMatched[1]
	// 	coreParsed := RustType(result.Name).Parse()
	// 	result.Core = &coreParsed
	// 	return result
	// }

	// result.Name = string(r)
	// return result
}
