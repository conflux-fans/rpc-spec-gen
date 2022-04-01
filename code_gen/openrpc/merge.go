package openrpc

import (
	"encoding/json"

	"github.com/conflux-fans/rpc-spec-gen/code_gen/openrpc/specconfig"
	"github.com/conflux-fans/rpc-spec-gen/parser/rust"
	"github.com/conflux-fans/rpc-spec-gen/utils"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

func CompleteDoc(doc OpenRPCSpec1, space string) OpenRPCSpec1 {

	if doc.OpenRPC == "" {
		doc.OpenRPC = "1.2.6"
	}

	if doc.Info.Version == "" {
		doc.Info.Version = "0.1.0"
	}

	if doc.Components == nil {
		doc.Components = &Components{}
	}

	if doc.Components.Schemas == nil {
		doc.Components.Schemas = make(map[string]*spec.Schema)
	}

	// 根据配置替换 method summary 和 result name
	specConfig := specconfig.GetSpecConfig(space)
	for _, m := range doc.Methods {
		methodConfig := specConfig.Methods[m.Name]
		m.Summary = methodConfig.Summary
		m.Description = methodConfig.Description
		m.Result.Name = methodConfig.ResultName
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

	useType := specconfig.ParseSchemaRefToUseType(schema.Ref.String())

	meta, ok := rust.GetUseTypeMeta(useType)
	if !ok {
		logger.WithFields(logrus.Fields{
			"useType":    useType.String(),
			"schema ref": schema.Ref.String(),
		}).Panic("meta is nil")
	}
	if meta.IsBaseType() {
		comp.Schemas[useType.Name] = specconfig.MustGetBasetypeSchemasByUseType(useType)
		return
	}

	comp.Schemas[useType.Name] = specconfig.MustLoadSchema(space, useType)
}

func convet2RefWithName(refWithFullname string) spec.Ref {
	useType := specconfig.ParseSchemaRefToUseType(refWithFullname)
	refWithName := spec.MustCreateRef(specconfig.SchemaRefRoot + useType.Name)
	return refWithName
}

func setSchemaRefBeName(s *spec.Schema) {
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
		useType := specconfig.ParseSchemaRefToUseType(s.Ref.String())
		realSchema := specconfig.MustLoadSchema(space, useType)

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

// 不会递归查找子项schema
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
