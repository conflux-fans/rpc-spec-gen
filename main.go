package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/Conflux-Chain/rpc-gen/code_gen/openrpc"
	"github.com/Conflux-Chain/rpc-gen/parser/rust"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var config = struct {
	RustRootPath string
}{
	RustRootPath: "/Users/wangdayong/myspace/mywork/conflux-rust",
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
// X 生成 open rpc 方法描述文件
func main() {
	// traitsFiles := []string{"/Users/wangdayong/myspace/mywork/conflux-rust/client/src/rpc/traits/cfx_space/pos.rs"}
	// parseds := []rust.TraitsFileParsed{}
	// for _, traitsFile := range traitsFiles {
	// 	traits, _ := ioutil.ReadFile(traitsFile)
	// 	parseds = append(parseds, rust.TraitsFile(traits).Parse())

	// 	j, _ := json.MarshalIndent(rust.TraitsFile(traits).Parse(), "", " ")
	// 	logrus.Info("traitsFile parsed result: ", string(j))
	// }
	// // return

	traitsFile := "/Users/wangdayong/myspace/mywork/conflux-rust/client/src/rpc/traits/cfx_space/pos_test.rs"

	space := path.Join(traitsFile, "..")[len("/Users/wangdayong/myspace/mywork/conflux-rust/client/src/rpc/traits/"):]

	traits, _ := ioutil.ReadFile(traitsFile)
	parsed := rust.TraitsFile(traits).Parse()

	j, _ := json.MarshalIndent(rust.TraitsFile(traits).Parse(), "", " ")
	logrus.Info("traitsFile parsed result: ", string(j))

	// var docs []openrpc.OpenRPCSpec1
	for _, trait := range parsed.Traits {
		// schemas := openrpc.GenSchemas(parsed.Uses)
		// j, _ := json.MarshalIndent(schemas, "", " ")
		// logger.Printf("schemas.json\n%s\n", j)

		openrpc.SaveSchemas(parsed.Uses, space)

		// if trait == nil {
		// 	logger.Printf("trait is nil")
		// 	os.Exit(1)
		// }
		doc := openrpc.GenDocTempalte(trait, parsed.Uses)
		openrpc.SaveDocTemplate(doc, space)

		doc = openrpc.CompleteDoc(doc, space)
		openrpc.SaveDoc(doc, space)

		// j, _ = json.MarshalIndent(doc, "", " ")
		// logger.Printf("%v.json\n%s\n", doc.Info.Title, j)

		// docPath := path.Join(gconfig.GetConfig().DocRootPath, space, trait.Name+".json")

		// if _, err := os.Stat(docPath); os.IsNotExist(err) {
		// 	folder := path.Join(docPath, "../")
		// 	os.MkdirAll(folder, 0700)
		// }

		// if e := ioutil.WriteFile(docPath, j, 0644); e != nil {
		// 	panic(e)
		// }
	}
}
