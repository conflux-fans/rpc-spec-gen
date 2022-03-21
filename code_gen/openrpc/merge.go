package openrpc

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Conflux-Chain/rpc-gen/parser/rust"
	"github.com/Conflux-Chain/rpc-gen/parser/rust/config"
	"github.com/Conflux-Chain/rpc-gen/utils"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

func CompleteDoc(doc OpenRPCSpec1, space string) OpenRPCSpec1 {
	if doc.Components == nil {
		doc.Components = &Components{}
	}

	if doc.Components.Schemas == nil {
		doc.Components.Schemas = make(map[string]*spec.Schema)
	}

	// 递归查找所有schema及子项相关schema
	schemas := getDocAllSchemas(doc, space)
	for _, s := range schemas {
		schemas = append(schemas, getRelatedSchemas(*s, space)...)
	}

	// 填充components.schemas
	for _, schema := range schemas {
		fillComponent(doc.Components, schema, space)
	}

	j, _ := json.MarshalIndent(doc, "", "  ")
	logger.WithField("doc", string(j)).Print("components completed")

	// 重置schema的ref为name
	schemas = getDocAllSchemas(doc, space)
	for _, schema := range schemas {
		setSchemaRefBeName(schema)
	}
	return doc
}

func fillComponent(comp *Components, schema *spec.Schema, space string) {
	if schema.Ref.String() == "" {
		return
	}

	useType := parseSchemaRefToUseType(schema.Ref.String())

	meta := config.GetUseTypeMeta(useType)
	if meta == nil {
		logger.WithFields(logrus.Fields{
			"useType":    useType.String(),
			"schema ref": schema.Ref.String(),
		}).Panic("meta is nil")
	}
	if meta.IsBaseType() {
		comp.Schemas[useType.Name] = mustGetBasetypeSchemasByUseType(useType)
		return
	}

	comp.Schemas[useType.Name] = mustLoadSchema(space, useType)
}

// // 填充 component schemas， 递归查找每个schema的field并将其schema填充到targets中， 最后返回值为该 refWithFullname 对应的 ref with name 值
// func recrusiveFillSchemas(refWithFullname string, space string, targets map[string]spec.Schema) spec.Ref {

// 	// logger.WithField("refWithFullname", refWithFullname).Debug("recrusiveFillSchemas")
// 	// if refWithFullname == "" {
// 	// 	time.Sleep(0)
// 	// }

// 	useType := parseSchemaRefToUseType(refWithFullname)

// 	refWithName := spec.MustCreateRef(schemaRefRoot + useType.Name)

// 	if _, ok := targets[useType.Name]; ok {
// 		return refWithName
// 	}

// 	meta := rustconfig.GetUseTypeMeta(useType)
// 	if meta == nil {
// 		panic("not found useType meta of " + useType.String())
// 	}

// 	if meta.IsBaseType() {
// 		targets[useType.Name] = basetypeSchemas[useType.Name]
// 		return refWithName
// 	}

// 	// schema := mustLoadSchema(space, useType)

// 	// // 设置 fields ref 为 name
// 	// for k, field := range schema.Properties {
// 	// 	if field.Ref.String() != "" {
// 	// 		_refWithName := recrusiveFillSchemas(field.Ref.String(), space, targets)
// 	// 		field.Ref = spec.MustCreateRef(_refWithName)
// 	// 		schema.Properties[k] = field
// 	// 	}
// 	// }

// 	// 添加 fields schema 到 targets
// 	// targets[useType.Name] = schema
// 	return refWithName
// }

func convet2RefWithName(refWithFullname string) spec.Ref {
	useType := parseSchemaRefToUseType(refWithFullname)
	refWithName := spec.MustCreateRef(schemaRefRoot + useType.Name)
	return refWithName
}

func setSchemaRefBeName(s *spec.Schema) {
	// s.Items
	// s.OneOf
	// s.Properties
	if s == nil {
		return
	}

	if s.Ref.String() != "" {
		refWithName := convet2RefWithName(s.Ref.String())
		s.Ref = refWithName
	}

	if s.Items != nil {
		setSchemaRefBeName(s.Items.Schema)
		for i := range s.Items.Schemas {
			setSchemaRefBeName(&s.Items.Schemas[i])
		}
	}

	if s.OneOf != nil {
		for i := range s.OneOf {
			setSchemaRefBeName(&s.OneOf[i])
		}
	}

	if s.Properties != nil {
		for k, v := range s.Properties {
			setSchemaRefBeName(&v)
			s.Properties[k] = v
		}
	}
}

// FIXME: s 为 指针会导致append(schemas, result...)后的所有值被最后一个值覆盖
// 会递归查找所有子项的schema
func getRelatedSchemas(s spec.Schema, space string) []*spec.Schema {

	var schemas []*spec.Schema = []*spec.Schema{&s}

	// 当为具体schema时
	if s.Items != nil {
		schemas = append(schemas, getRelatedSchemas(*s.Items.Schema, space)...)
		for i := range s.Items.Schemas {
			schemas = append(schemas, getRelatedSchemas(s.Items.Schemas[i], space)...)
		}
	}

	if s.OneOf != nil {
		for i := range s.OneOf {
			schemas = append(schemas, getRelatedSchemas(s.OneOf[i], space)...)
		}
	}

	if s.Properties != nil {
		for k, v := range s.Properties {
			result := getRelatedSchemas(v, space)
			j, _ := json.MarshalIndent(result, "", "  ")
			logger.WithFields(logrus.Fields{
				"schema": k,
				"result": string(j),
			}).Debug("getRelatedSchemas")
			schemas = append(schemas, result...)
			j, _ = json.MarshalIndent(schemas, "", "  ")
			logger.WithField("schemas", string(j)).Debug("append schemas")
		}
	}

	// 当为引用时，根据ref找到usetype，然后loadschema，然后 get related schemas
	if s.Ref.String() != "" {
		useType := parseSchemaRefToUseType(s.Ref.String())
		realSchema := mustLoadSchema(space, useType)

		if realSchema == nil {
			logger.WithField("useType", useType).Panic("not found schema")
		}

		logger.WithFields(logrus.Fields{
			"space":      useType.String(),
			"realSchema": utils.MustJsonPretty(realSchema),
			"result":     getRelatedSchemas(*realSchema, space),
		}).Debug("getRelatedSchemas")
		schemas = append(schemas, getRelatedSchemas(*realSchema, space)...)
	}

	return schemas
}

// 不会递归查找
func getDocAllSchemas(doc OpenRPCSpec1, space string) []*spec.Schema {
	// 查找所有schema
	var schemas []*spec.Schema
	for i := range doc.Methods {
		for _, p := range doc.Methods[i].Params {
			schemas = append(schemas, &p.Schema)
		}

		if doc.Methods[i].Result != nil {
			schemas = append(schemas, &doc.Methods[i].Result.Schema)
		}
	}

	for k, _ := range doc.Components.Schemas {
		schemas = append(schemas, doc.Components.Schemas[k])
	}
	return schemas
}

func mustLoadSchema(space string, useType rust.UseType) *spec.Schema {

	if config.IsBaseType(useType.String()) {
		return mustGetBasetypeSchemasByUseType(useType)
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

	j, _ := json.MarshalIndent(schema, "", "  ")

	logger.WithField("schema", string(j)).Debug("load schema")
	return &schema
}
