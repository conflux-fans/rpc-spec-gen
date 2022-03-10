package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type RustStruct string

type RustType string

// RustStructParsed reperesnts parse result of a Rust struct
type RustStructParsed struct {
	Comment string
	Name    string
	Fields  []RustFieldParsed
}

// RustFieldParsed reperesnts parse result of a Rust struct field
type RustFieldParsed struct {
	Comment string
	Name    string
	Type    RustTypeParsed
}

// RustTypeParsed reperesnts parse result of a Rust struct field type
type RustTypeParsed struct {
	IsOption bool
	IsArray  bool
	Name     string
	Core     *RustTypeParsed
}

func (r RustStruct) Parse() RustStructParsed {
	structReg := regexp.MustCompile(`(?Us)(.*)pub struct (.*)\{(.*)\}`)
	structFinded := structReg.FindStringSubmatch(string(r))
	// fmt.Printf("%#v\n", structFinded)

	sComment, sName, sBody := strings.TrimSpace(structFinded[1]), strings.TrimSpace(structFinded[2]), strings.TrimSpace(structFinded[3])

	fmt.Printf("comment %v\nhead %#v\nbody %#v\n", sComment, sName, sBody)

	fieldReg := regexp.MustCompile(`(?Us)(.*)pub (.*): (.*),`)
	fieldsFinded := fieldReg.FindAllStringSubmatch(sBody, -1)
	// fmt.Printf("fieldsFinded %#v\n", fieldsFinded[0])

	fields := make([]RustFieldParsed, len(fieldsFinded))
	for i, field := range fieldsFinded {
		fmt.Printf("field %#v\n", field)
		fComment, fName, fType := strings.TrimSpace(field[1]), strings.TrimSpace(field[2]), RustType(field[3])
		fields[i] = RustFieldParsed{fComment, fName, fType.Parse()}
	}

	return RustStructParsed{sComment, sName, fields}
}

func (r RustType) Parse() (result RustTypeParsed) {

	optionReg := regexp.MustCompile(`Option<(.*)>`)
	optionMatched := optionReg.FindStringSubmatch(string(r))

	fmt.Printf("optionMatched %#v\n", optionMatched)
	if len(optionMatched) > 0 {
		result.IsOption = true
		result.Name = optionMatched[1]
		coreParsed := RustType(result.Name).Parse()
		result.Core = &coreParsed
		return result
	}

	vecReg := regexp.MustCompile(`Vec<(.*)>`)
	vecMatched := vecReg.FindStringSubmatch(string(r))
	fmt.Printf("vecMatched %#v\n", vecMatched)
	if len(vecMatched) > 0 {
		result.IsArray = true
		result.Name = vecMatched[1]
		coreParsed := RustType(result.Name).Parse()
		result.Core = &coreParsed
		return result
	}

	result.Name = string(r)
	return result
}
