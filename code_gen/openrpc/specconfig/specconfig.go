package specconfig

import "github.com/conflux-fans/rpc-spec-gen/code_gen/openrpc/types"

type SpecConfig struct {
	Space string
	// TraitName string
	Info    *types.Info
	Methods map[string]*MethodConfig
	Servers []*types.Server
}

type MethodConfig struct {
	Summary     string
	ParamNames  []string
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

// space: ethspace or cfxspace
// traitName: trait name, such as "Eth","Parity"
func GetSpecConfig(space string) (*SpecConfig, bool) {
	v, ok := specConfigs[space]
	return v, ok
}
