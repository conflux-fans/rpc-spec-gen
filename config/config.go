package config

type Config struct {
	RustRootPath        string
	SchemaRootPath      string
	DocRootPath         string
	DocTemplateRootPath string
}

var config Config

func init() {
	config = Config{
		RustRootPath:        "/Users/dayong/myspace/mywork/conflux-rust",
		SchemaRootPath:      "./output/schemas",
		DocRootPath:         "./output/doc",
		DocTemplateRootPath: "./output/doc_template",
	}
}

func GetConfig() Config {
	return config
}
