package specconfig

type SpecConfig struct {
	Space   string
	Methods map[string]MethodConfig
}

type MethodConfig struct {
	Summary     string
	ResultName  string
	Description string
}

var specConfigs map[string]*SpecConfig = make(map[string]*SpecConfig)

func init() {
	cs := []SpecConfig{
		getEthSpaceSpecConfig(),
	}

	for _, c := range cs {
		specConfigs[c.Space] = &c
	}
}

func GetSpecConfig(space string) *SpecConfig {
	return specConfigs[space]
}
