package rust

type TraitsFile string

type traitsFileParsedMiddle struct {
	traits []Trait
	uses   []Use
}

type TraitsFileParsed struct {
	Traits []TraitParsed
	Uses   []UseType
}

func (rf TraitsFile) Parse() TraitsFileParsed {
	resultM := traitsFileParsedMiddle{}

	// traitsReg := regexp.MustCompile(`(?mUs)(\/\/\/.*\n|)#\[rpc\(.*\)\]\npub trait .* \{[\s\S]*}`)
	// traitsFinded := traitsReg.FindAllString(string(rf), -1)
	// fmt.Printf("traitRegFinded len %v, %v\n", len(traitsFinded), traitsFinded)

	// resultM.traits = make([]Trait, len(traitsFinded))
	// for i, trait := range traitsFinded {
	// 	resultM.traits[i] = Trait(trait)
	// }

	resultM.traits, resultM.uses = GetTraits(string(rf))
	// resultM.uses = GetUses(string(rf))
	// useReg := regexp.MustCompile(`(?mUs)use .*;`)
	// useFinded := useReg.FindAllString(string(rf), -1)
	// fmt.Printf("useFinded %v\n", useFinded)
	// for _, use := range useFinded {
	// 	resultM.uses = append(resultM.uses, Use(use))
	// }

	// useReg := regexp.MustCompile(`(?mUs)use (.*);`)
	// useFinded := useReg.FindAllStringSubmatch(string(rf), -1)
	// fmt.Printf("useFinded %v\n", useFinded)

	// resultM.uses = make([]Use, len(useFinded))
	// for i, use := range useFinded {
	// 	resultM.uses[i] = Use(use[1])
	// }

	return resultM.parseCore()
}

func (rm traitsFileParsedMiddle) parseCore() TraitsFileParsed {
	result := TraitsFileParsed{}
	result.Traits = make([]TraitParsed, len(rm.traits))
	for i, trait := range rm.traits {
		result.Traits[i] = trait.Parse()
	}

	result.Uses = make([]UseType, 0)
	for _, use := range rm.uses {
		result.Uses = append(result.Uses, use.Parse()...)
	}
	return result
}
