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

func (tf TraitsFile) Parse() TraitsFileParsed {
	resultM := traitsFileParsedMiddle{}
	resultM.traits, resultM.uses = NewSouceCode(string(tf)).GetTraits()
	return resultM.parseCore()
}

func (m traitsFileParsedMiddle) parseCore() TraitsFileParsed {
	result := TraitsFileParsed{}
	result.Traits = make([]TraitParsed, len(m.traits))
	for i, trait := range m.traits {
		result.Traits[i] = trait.Parse()
	}

	result.Uses = make([]UseType, 0)
	for _, use := range m.uses {
		result.Uses = append(result.Uses, use.Parse()...)
	}
	return result
}
