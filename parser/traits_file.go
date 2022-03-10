package parser

import (
	"fmt"
	"regexp"
)

type RustTraitsFile string

type rustTraitsFileParsedMiddle struct {
	traits []RustTrait
	uses   []RustUse
}

type RustTraitsFileParsed struct {
	Traits []RustTraitParsed
	Uses   []RustUseParsed
}

func (rf RustTraitsFile) Parse() RustTraitsFileParsed {
	resultM := rustTraitsFileParsedMiddle{}

	traitsReg := regexp.MustCompile(`(?mUs)(\/\/\/.*\n|)#\[rpc\(.*\)\]\npub trait .* \{[\s\S]*}`)
	traitsFinded := traitsReg.FindAllString(string(rf), -1)
	fmt.Printf("traitRegFinded len %v, %v\n", len(traitsFinded), traitsFinded)

	resultM.traits = make([]RustTrait, len(traitsFinded))
	for i, trait := range traitsFinded {
		resultM.traits[i] = RustTrait(trait)
	}

	useReg := regexp.MustCompile(`(?mUs)use (.*);`)
	useFinded := useReg.FindAllString(string(rf), -1)
	fmt.Printf("useFinded %v\n", useFinded)

	resultM.uses = make([]RustUse, len(useFinded))
	for i, use := range useFinded {
		resultM.uses[i] = RustUse(use)
	}

	return resultM.parseCore()
}

func (rm rustTraitsFileParsedMiddle) parseCore() RustTraitsFileParsed {
	result := RustTraitsFileParsed{}
	result.Traits = make([]RustTraitParsed, len(rm.traits))
	for i, trait := range rm.traits {
		result.Traits[i] = trait.Parse()
	}
	result.Uses = make([]RustUseParsed, len(rm.uses))
	for i, use := range rm.uses {
		result.Uses[i] = use.Parse()
	}
	return result
}
