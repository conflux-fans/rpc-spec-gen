package openrpc

import (
	"encoding/json"
	"io/ioutil"

	rustconfig "github.com/Conflux-Chain/rpc-gen/parser/rust/config"
	"github.com/go-openapi/spec"
)

func CompleteDoc(doc OpenRPCSpec1, space string) OpenRPCSpec1 {

	// 填充components
	if doc.Components == nil {
		doc.Components = &Components{}
	}

	if doc.Components.Schemas == nil {
		doc.Components.Schemas = make(map[string]spec.Schema)
	}

	for mk, m := range doc.Methods {
		for _, p := range m.Params {
			// schemas = append(schemas, p.Schema)
			// 修改shcema ref 为 name
			// useType := parseSchemaRefToUseType(p.Schema.Ref.String())
			// p.Schema.Ref = spec.MustCreateRef(schemaRefRoot + useType.Name)

			newSchema := recrusiveFillSchemas(p.Schema.Ref.String(), space, doc.Components.Schemas)
			p.Schema.Ref = spec.MustCreateRef(newSchema)
		}

		if m.Result != nil {
			newSchema := recrusiveFillSchemas(m.Result.Schema.Ref.String(), space, doc.Components.Schemas)
			doc.Methods[mk].Result.Schema.Ref = spec.MustCreateRef(newSchema)
		}
	}
	return doc

	// // 查找所有schema
	// var schemas []spec.Schema
	// for _, m := range doc.Methods {
	// 	for _, p := range m.Params {
	// 		schemas = append(schemas, p.Schema)
	// 		// 修改shcema ref 为 name
	// 		useType := parseSchemaRefToUseType(p.Schema.Ref.String())
	// 		p.Schema.Ref = spec.MustCreateRef(schemaRefRoot + useType.Name)
	// 	}
	// 	schemas = append(schemas, m.Result.Schema)
	// }

	// // 填充components
	// if doc.Components == nil {
	// 	doc.Components = &Components{}
	// }

	// if doc.Components.Schemas == nil {
	// 	doc.Components.Schemas = make(map[string]spec.Schema)
	// }

	// for _, schema := range schemas {
	// 	useType := parseSchemaRefToUseType(schema.Ref.String())

	// 	if rustconfig.GetUseTypeMeta(useType).IsBaseType() {
	// 		doc.Components.Schemas[useType.Name] = basetypeSchemas[useType.Name]
	// 		continue
	// 	}

	// 	savePath := getSchemaSavePath(space, useType.String())
	// 	content, e := ioutil.ReadFile(savePath)

	// 	if e != nil {
	// 		panic(e)
	// 	}

	// 	tmp := spec.Schema{}
	// 	if e := json.Unmarshal(content, &tmp); e != nil {
	// 		panic(e)
	// 	}

	// 	doc.Components.Schemas[useType.Name] = tmp
	// }
	// return doc
}

// 填充 component schemas， 递归查找每个schema的field并将其schema填充到targets中， 最后返回值为该 refWithFullname 对应的 ref with name 值
func recrusiveFillSchemas(refWithFullname string, space string, targets map[string]spec.Schema) string {

	useType := parseSchemaRefToUseType(refWithFullname)

	refWithName := spec.MustCreateRef(schemaRefRoot + useType.Name)

	if _, ok := targets[useType.Name]; ok {
		return refWithName.String()
	}

	meta := rustconfig.GetUseTypeMeta(useType)
	if meta == nil {
		panic("not found useType meta of " + useType.String())
	}

	if meta.IsBaseType() {
		targets[useType.Name] = basetypeSchemas[useType.Name]
		return refWithName.String()
	}

	savePath := getSchemaSavePath(space, useType.String())
	content, e := ioutil.ReadFile(savePath)

	if e != nil {
		panic(e)
	}

	schema := spec.Schema{}
	if e := json.Unmarshal(content, &schema); e != nil {
		panic(e)
	}

	// 设置 fields ref 为 name
	for k, field := range schema.Properties {
		_refWithName := recrusiveFillSchemas(field.Ref.String(), space, targets)
		field.Ref = spec.MustCreateRef(_refWithName)
		schema.Properties[k] = field
	}
	// 添加 fields schema 到 targets
	targets[useType.Name] = schema
	return refWithName.String()
}
