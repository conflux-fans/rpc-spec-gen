package main

import (
	"io/ioutil"

	"github.com/Conflux-Chain/rpc-gen/code_gen/openrpc"
	"github.com/Conflux-Chain/rpc-gen/parser"
)

// V rpc 方法解析 参数 和 返回值
// V 根据 rpc 类型 <-> 对应的 rust 文件路径 寻找类型代码
// V 解析出字段类型
// X 生成 open rpc 方法描述文件
func main() {
	traitsFiles := []string{""}
	parseds := []parser.RustTraitsFileParsed{}
	for _, traitsFile := range traitsFiles {
		traits, _ := ioutil.ReadFile(traitsFile)
		parseds = append(parseds, parser.RustTraitsFile(traits).Parse())
	}

	openrpc.GenDoc(parseds)

}
