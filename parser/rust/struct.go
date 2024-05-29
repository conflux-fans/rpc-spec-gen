package rust

import (
	"regexp"
	"strings"

	"github.com/conflux-fans/rpc-spec-gen/utils"
	"github.com/sirupsen/logrus"
)

type Struct string

type RustType string

// StructParsed reperesnts parse result of a Rust struct
type StructParsed struct {
	Comment string
	Name    string
	// TODO: 还没有解析
	IsDriveSerialize bool
	Fields           []FieldParsed
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
	if t.Core == nil {
		return t.Name
	}
	return t.Core.InnestCoreTypeName()
}

// Default define name for params if param name is not specified
func (t *TypeParsed) DefaultDefineName() string {
	core := strings.ToLower(t.InnestCoreTypeName())
	if t.IsArray {
		return core + "s"
	}
	return core
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
	sComment = utils.CleanComment(sComment)

	// logrus.WithField("comment %v\nhead %#v\nbody %#v\n", sComment, sName, sBody)
	logrus.WithField("comment", sComment).WithField("head", sName).WithField("body", sBody).Debug("split struct content")

	// TODO： 解析LedgerInfoWithV0不正确
	// fieldReg := regexp.MustCompile(`(?Us)(.*)(?:pub|) (.*): (.*),`)
	fieldReg := regexp.MustCompile(`(?Us)(.*)(?:pub) (.*): (.*),|(.*)(.*): (.*),`)
	// fieldReg := regexp.MustCompile(`(?Us)(.*)(pub|)\s*(.*):\s*(.*),`)
	fieldsFinded := fieldReg.FindAllStringSubmatch(sBody, -1)
	// fmt.Printf("fieldsFinded %#v\n", fieldsFinded[0])

	fields := make([]FieldParsed, len(fieldsFinded))
	for i, field := range fieldsFinded {
		// fmt.Printf("field %#v\n", field)
		logrus.WithField("field", field).Debug("find struct field")
		fComment, fName, fType := strings.TrimSpace(field[1]), strings.TrimSpace(field[2]), RustType(field[3])
		if fComment == "" && fName == "" && fType == "" {
			fComment, fName, fType = strings.TrimSpace(field[4]), strings.TrimSpace(field[5]), RustType(field[6])
		}
		fComment = utils.CleanComment(fComment)

		fields[i] = FieldParsed{fComment, fName, fType.Parse()}
		logrus.WithFields(logrus.Fields{
			"field":        field[0],
			"field parsed": fields[i],
		}).Debug("field parsed")
	}

	return StructParsed{sComment, sName, false, fields}
}

func (r RustType) Parse() (result TypeParsed) {

	inner, e := re_AngleBrackets.FindStringMatch(string(r))
	if e != nil {
		logrus.WithField("rustType", r).WithError(e).Panic("failed to found inner type")
	}
	// logrus.WithField("rustType", r).WithField("inner", inner).Debug("find inner")

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
	case "VecDeque":
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
