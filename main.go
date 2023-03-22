package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/conflux-fans/rpc-spec-gen/code_gen/openrpc"
	"github.com/conflux-fans/rpc-spec-gen/parser/rust"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var config = struct {
	RustRootPath  string
	TraitRootPath string
}{
	RustRootPath:  "/Users/dayong/myspace/mywork/conflux-rust/",
	TraitRootPath: "client/src/rpc/traits/",
}

var logger = &logrus.Logger{
	Out:   os.Stdout,
	Level: logrus.DebugLevel,
	Formatter: &prefixed.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	},
}

// V rpc 方法解析 参数 和 返回值
// V 根据 rpc 类型 <-> 对应的 rust 文件路径 寻找类型代码
// V 解析出字段类型
// V 生成 open rpc 方法描述文件
func main() {
	// traitsFile := "/Users/wangdayong/myspace/mywork/conflux-rust/client/src/rpc/traits/cfx_space/cfx.rs"
	traitsFile := "/Users/dayong/myspace/mywork/conflux-rust/client/src/rpc/traits/cfx_space/cfx.rs"

	// traitsFile = "/Users/wangdayong/myspace/mywork/conflux-rust/client/src/rpc/traits/eth_space/eth.rs"

	space := path.Join(traitsFile, "..")[len(config.RustRootPath+config.TraitRootPath):]

	traits, err := ioutil.ReadFile(traitsFile)
	if err != nil {
		panic(err)
	}
	parsed := rust.TraitsFile(traits).Parse()

	j, _ := json.MarshalIndent(rust.TraitsFile(traits).Parse(), "", " ")
	logger.Info("traitsFile parsed result: ", string(j))

	for _, trait := range parsed.Traits {
		openrpc.SaveSchemas(parsed.Uses, space)
		doc := openrpc.GenDocTempalte(trait, parsed.Uses)
		openrpc.SaveDocTemplate(doc, space)

		doc = openrpc.CompleteDoc(doc, space, trait.Name)
		openrpc.SaveDoc(doc, space, trait.Name)
	}
}
