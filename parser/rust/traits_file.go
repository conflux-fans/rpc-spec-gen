package parser

import (
	"fmt"
	"regexp"
)

type RustTraitsFile string

type rustTraitsFileParsedMiddle struct {
	traits []Trait
	uses   []Use
}

type RustTraitsFileParsed struct {
	Traits []TraitParsed
	Uses   []UseParsed
}

func (rf RustTraitsFile) Parse() RustTraitsFileParsed {
	resultM := rustTraitsFileParsedMiddle{}

	traitsReg := regexp.MustCompile(`(?mUs)(\/\/\/.*\n|)#\[rpc\(.*\)\]\npub trait .* \{[\s\S]*}`)
	traitsFinded := traitsReg.FindAllString(string(rf), -1)
	fmt.Printf("traitRegFinded len %v, %v\n", len(traitsFinded), traitsFinded)

	resultM.traits = make([]Trait, len(traitsFinded))
	for i, trait := range traitsFinded {
		resultM.traits[i] = Trait(trait)
	}

	useReg := regexp.MustCompile(`(?mUs)use (.*);`)
	useFinded := useReg.FindAllString(string(rf), -1)
	fmt.Printf("useFinded %v\n", useFinded)

	resultM.uses = make([]Use, len(useFinded))
	for i, use := range useFinded {
		resultM.uses[i] = Use(use)
	}

	return resultM.parseCore()
}

func (rm rustTraitsFileParsedMiddle) parseCore() RustTraitsFileParsed {
	result := RustTraitsFileParsed{}
	result.Traits = make([]TraitParsed, len(rm.traits))
	for i, trait := range rm.traits {
		result.Traits[i] = trait.Parse()
	}
	result.Uses = make([]UseParsed, len(rm.uses))
	for i, use := range rm.uses {
		result.Uses[i] = use.Parse()
	}
	return result
}
