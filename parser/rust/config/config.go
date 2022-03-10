package config

type RustUseTypeMeta struct {
	isBaseType bool
	file       string
}

func (r *RustUseTypeMeta) IsBaseType() bool {
	return r.isBaseType
}

func (r *RustUseTypeMeta) InFilePath() string {
	return r.file
}

// Struct exist in rust file path
var RustUseTypeMetas map[string]RustUseTypeMeta = map[string]RustUseTypeMeta{
	"cfx_types::H160": {true, ""},
	"cfx_types::H256": {true, ""},
	"cfx_types::U256": {true, ""},
}
