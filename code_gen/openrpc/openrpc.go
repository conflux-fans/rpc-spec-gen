package openrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/conflux-fans/rpc-spec-gen/config"
	"github.com/conflux-fans/rpc-spec-gen/parser/rust"
	"github.com/conflux-fans/rpc-spec-gen/utils"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var usetype2Schema = map[string]spec.Schema{}

var logger = &logrus.Logger{
	Out:   os.Stderr,
	Level: logrus.DebugLevel,
	Formatter: &prefixed.TextFormatter{
		// DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	},
}

func init() {
	// init enum2Schema for enums which custom implement json serailize
	for k, v := range customSchemas {
		usetype2Schema[k] = *v
	}
}

func GenSchemaByStruct(structParsed rust.StructParsed, defaultModPath []string) spec.Schema {

	// logger.WithField("struct parsed", utils.MustJsonPretty(structParsed)).Info("GenSchemaByStruct")

	if structParsed.Name == "CfxRpcLogFilter" {
		time.Sleep(0)
	}

	s := spec.Schema{}
	s.Title = structParsed.Comment
	s.Type = spec.StringOrArray{"object"}
	s.Properties = make(map[string]spec.Schema, len(structParsed.Fields))

	for _, field := range structParsed.Fields {
		// 生成field的schema ref
		// 	先从 usetype2Schema , basetypes 匹配
		// 	若都未匹配到则设置为与struct同级别usetype
		refSchema := genObjRefSchema(field.Type, defaultModPath)

		fNameCamleCase := utils.UnderScoreCase2CamelCase(field.Name, false)
		refSchema.Title = fNameCamleCase

		s.Properties[fNameCamleCase] = *refSchema
	}

	return s
}

func GenSchemaByEnum(enumParsed rust.EnumParsed, defaultModPath []string) spec.Schema {

	if v, ok := usetype2Schema[enumParsed.Name]; ok {
		return v
	}

	s := spec.Schema{}
	s.Title = enumParsed.Comment
	s.Properties = make(map[string]spec.Schema, len(enumParsed.Fields))

	hasTumple := false
	enums := make([]interface{}, 0)
	for _, field := range enumParsed.Fields {
		// 是否tumple
		// 如果是tumple，参数的ref生成步骤为：
		// 	先从 usetype2Schema , basetypes 匹配
		// 	若都未匹配到则设置为与enum同级别usetype
		// 否则生成enum

		if field.IsTumple() {
			hasTumple = true
			if len(field.TupleParams) > 1 {
				logger.WithField("field", field).Panic("enum tumple field should be custom")
			}

			refSchema := genObjRefSchema(field.TupleParams[0], defaultModPath)
			s.OneOf = append(s.OneOf, *refSchema)
			continue
		}
		underScoreCase := utils.CamelCase2UnderScoreCase(field.Value)
		enums = append(enums, underScoreCase)
	}

	if len(enums) == 0 {
		return s
	}

	if !hasTumple {
		s.Enum = enums
		return s
	}

	s.OneOf = append(s.OneOf, spec.Schema{
		SchemaProps: spec.SchemaProps{
			Enum: enums,
		},
	})
	return s
}

func GenSchemas(useTypes []rust.UseType) map[string]spec.Schema {

	logger.WithField("useTypes", useTypes).Debug("GenSchemas")

	for _, useType := range useTypes {
		// get code file path by usetype
		if _, ok := usetype2Schema[useType.String()]; ok {
			continue
		}

		meta, ok := rust.GetUseTypeMeta(useType)

		if !ok {
			panic(fmt.Sprintf("not found meta of %v", useType.String()))
		}

		if meta.IsBaseType() || meta.IsIgnore() {
			continue
		}

		// 从配置中查看类型所在文件
		// 解析所在文件，拿到所有structs/enums和use types
		// 判断structs中是否有该类型，有则生成schema
		//     先生成字段的schema
		//     再生成该struct的schema
		// 判断enums中是否有该类型，有则生成schema
		//     先生成子项的schema
		//     再生成该enum的schema
		code, e := ioutil.ReadFile(meta.InFilePath())
		if e != nil {
			logrus.
				WithField("useType", useType).
				WithField("filePath", meta.InFilePath()).
				WithError(e).
				Panic("read file error")
		}
		// s, us := rust.FindStruct(string(code), useType.Name)

		// 获得code中的所有struct，以struct名为key，struct为value
		structs, _ := rust.GetStructs(string(code))
		// 获得code中的所有enum，以enum名为key，enum为value
		enums, us := rust.GetEnums(string(code))

		logger.WithFields(logrus.Fields{
			"code":     string(code),
			"structs":  structs,
			"enums":    enums,
			"usetypes": us,
		}).Debug("GetStructs and Enums from code")

		if _strcut, ok := structs[useType.Name]; ok {
			fieldUsetypes := getStructFieldUseTypes(useType, us, structs, enums)
			logger.WithFields(logrus.Fields{
				"struct in useType": useType,
				"filtered us":       fieldUsetypes,
			}).Debug("filter struct field using use types")
			GenSchemas(fieldUsetypes)

			usetype2Schema[useType.String()] = GenSchemaByStruct(_strcut.Parse(), useType.ModPath)
			continue
		}

		if _enum, ok := enums[useType.Name]; ok {
			fieldUsetypes := getEnumFieldUseTypes(useType, us, enums, structs)
			logger.WithFields(logrus.Fields{
				"enum in useType": useType,
				"filtered us":     fieldUsetypes,
			}).Debug("filter struct field using use types")
			GenSchemas(fieldUsetypes)

			usetype2Schema[useType.String()] = GenSchemaByEnum(_enum.Parse(), useType.ModPath)
			continue
		}

		logrus.WithFields(logrus.Fields{
			"code":           string(code),
			"structs finded": structs,
			"use types":      us,
		}).Panicf("not found struct <%v> from code", useType.Name)
		panic("not found struct")
	}
	return usetype2Schema
}

func SaveSchemas(useTypes []rust.UseType, space string) {

	logger.WithField("useTypes", useTypes).Debug("SaveSchemas")

	schemas := GenSchemas(useTypes)
	j, _ := json.MarshalIndent(schemas, "", "  ")
	logger.WithField("schemas", string(j)).Debug("Generated schemas")

	for k, schema := range schemas {
		j, _ := json.MarshalIndent(schema, "", "  ")
		p := getSchemaSavePath(space, k)

		saveFile(p, j)
	}
}

// - 方法的参数和返回值都是schema的ref，放到methods.params和methods.result
// - 查看是否有已生成的scheme，没有则创建
func GenMethod(funcParsed rust.FuncParsed, useTypePool []rust.UseType) *Method {
	var method Method
	method.Summary = funcParsed.Comment
	method.Name = funcParsed.RpcMethod
	method.Params = make([]*ContentDescriptor, len(funcParsed.Params))
	for i, param := range funcParsed.Params {
		// ut := mustFindUseType(param.Type.InnestCoreTypeName(), useTypes)
		// method.Params[i] = getParamContentDescriptor(*ut, param)
		method.Params[i] = getParamContentDescriptor(param, useTypePool)
	}

	// ut := findUseType(funcParsed.Return.Type.InnestCoreTypeName(), useTypePool)

	// if ut == nil {
	// 	logger.WithFields(logrus.Fields{
	// 		"Func Method": funcParsed.RpcMethod,
	// 		"Name":        funcParsed.Return.Type.Name,
	// 		"Use Types":   useTypePool,
	// 	}).Panic("not found use type")
	// }

	// method.Result = getResultContentDescriptor(*ut, funcParsed.Return)
	method.Result = getResultContentDescriptor(funcParsed.Return, useTypePool)

	return &method
}

func getParamContentDescriptor(p rust.ParamParsed, useTypePool []rust.UseType) *ContentDescriptor {

	u := mustFindUseType(p.Type.InnestCoreTypeName(), useTypePool)
	refSchema := getUseTypeRefSchema(*u)

	// c := Content{
	// 	Name:   p.Name,
	// 	Schema: spec.Schema{},
	// }

	// for {
	// 	if p.Type.IsOption {
	// 		c.Required = !p.Type.IsOption
	// 		continue
	// 	}

	// 	if p.Type.IsArray {
	// 		c.Schema.Type = spec.StringOrArray{"array"}
	// 		c.Schema.Items = &spec.SchemaOrArray{Schema: &refSchema}
	// 		continue
	// 	}

	// 	c.Schema = refSchema
	// }

	c := Content{
		Name:     p.Name,
		Required: !p.Type.IsOption,
		Schema:   *genSchemaForParsedType(p.Type, refSchema),
	}

	return &ContentDescriptor{Content: c}
}

func getResultContentDescriptor(p rust.ReturnParsed, useTypePool []rust.UseType) *ContentDescriptor {
	u := mustFindUseType(p.Type.InnestCoreTypeName(), useTypePool)
	refSchema := getUseTypeRefSchema(*u)

	c := Content{
		Name:   p.Name,
		Schema: *genSchemaForParsedType(p.Type, refSchema),
	}

	return &ContentDescriptor{Content: c}

	// c := Content{
	// 	Name:   r.Name,
	// 	Schema: spec.Schema{},
	// }

	// if r.Type.IsArray {
	// 	c.Schema.Type = spec.StringOrArray{"array"}
	// 	c.Schema.Items = &spec.SchemaOrArray{Schema: &refSchema}
	// } else {
	// 	c.Schema = refSchema
	// }

	// if r.Type.IsOption {
	// 	c.Schema.Nullable = true
	// }

	// return &ContentDescriptor{Content: c}
}

// func genSchemaForParsedType(t rust.TypeParsed, coreRefSchema spec.Schema) spec.Schema {
// 	s := spec.Schema{}

func genSchemaForParsedType(t rust.TypeParsed, coreRefSchema spec.Schema) *spec.Schema {
	if t.Core == nil {
		return &coreRefSchema
	}

	coreSchema := genSchemaForParsedType(*t.Core, coreRefSchema)
	items := &spec.SchemaOrArray{Schema: coreSchema}
	if t.IsOption {
		s := coreSchema
		s.Nullable = true
		return s
	}

	if t.IsArray {
		s := &spec.Schema{}
		s.Type = spec.StringOrArray{"array"}
		s.Items = items
		return s
	}
	if t.IsVariadicValue {
	}
	return nil
}

// 	return s
// }

func GenMethods(traitParsed rust.TraitParsed, useTypes []rust.UseType) []*Method {
	var methods []*Method
	for _, funcParsed := range traitParsed.Funcs {
		methods = append(methods, GenMethod(funcParsed, useTypes))
	}
	return methods
}

// uses -> schemas
// traits -> funcs -> methods
func GenDocTempalte(trait rust.TraitParsed, useTypes []rust.UseType) OpenRPCSpec1 {

	doc := OpenRPCSpec1{}
	doc.Info.Title = trait.Name
	doc.Info.Description = trait.Comment
	doc.Methods = GenMethods(trait, useTypes)

	return doc
}

func SaveDocTemplate(doc OpenRPCSpec1, space string) {
	docPath := path.Join(config.GetConfig().DocTemplateRootPath, space, doc.Info.Title+".json")
	j, _ := json.MarshalIndent(doc, "", " ")
	saveFile(docPath, j)
}

func SaveDoc(doc OpenRPCSpec1, space string) {
	docPath := path.Join(config.GetConfig().DocRootPath, space, doc.Info.Title+".json")
	j, _ := json.MarshalIndent(doc, "", " ")
	saveFile(docPath, j)
}

func saveFile(docPath string, content []byte) {
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		folder := path.Join(docPath, "../")
		os.MkdirAll(folder, 0700)
	}

	if e := ioutil.WriteFile(docPath, content, 0644); e != nil {
		panic(e)
	}
}

func getUseType(name string, defaultModPath []string, useTypes []rust.UseType) *rust.UseType {
	for _, useType := range useTypes {
		if useType.Name == name {
			return &useType
		}
	}

	if rust.IsBaseType(name) {
		tmp := rust.MustNewUseType(name)
		return &tmp
	}

	return getSameLevelUseType(name, defaultModPath)

}

func findUseType(name string, useTypes []rust.UseType) *rust.UseType {

	for _, useType := range useTypes {
		// if useType.Name == name {
		if useType.Alias == name {
			return &useType
		}
	}

	if rust.IsBaseType(name) {
		ut := rust.MustNewUseType(name)
		return &ut
	}
	return nil
}

func mustFindUseType(name string, useTypes []rust.UseType) *rust.UseType {
	ut := findUseType(name, useTypes)
	if ut == nil {
		logger.WithFields(logrus.Fields{
			"Name":      name,
			"Use Types": useTypes,
		}).Panic("not found use type")
	}
	return ut
}

func getSameLevelUseType(name string, defaultModPath []string) *rust.UseType {
	ut := rust.UseType{}
	ut.ModPath = defaultModPath
	ut.Name = name
	ut.Alias = name
	return &ut
}

func isFieldsHasType(sp rust.StructParsed, t rust.UseType) bool {
	for _, field := range sp.Fields {
		if field.Type.Name == t.Name {
			return true
		}
	}
	return false
}

// 寻找useType对应struct的struct fields中在usetypes的匹配类型，如果没有在当前文件的所有structs中寻找
func getStructFieldUseTypes(aim rust.UseType, usePool []rust.Use, structsPool map[string]rust.Struct, enumsPool map[string]rust.Enum) []rust.UseType {
	var founds []rust.UseType
	if _struct, ok := structsPool[aim.Name]; ok {
		// 只过滤field里边包含的useType
		_p := _struct.Parse()

		for _, field := range _p.Fields {
			fCoreType := field.Type.InnestCoreTypeName()

			logrus.WithFields(
				logrus.Fields{
					"aim":             aim,
					"field core type": fCoreType,
					"is base type":    rust.IsBaseType(fCoreType),
				},
			).Info("check is base type")
			if rust.IsBaseType(fCoreType) {
				continue
			}

			// 如果没有找到从同文件的enums中找enum
			finded := findFieldUseType(aim, fCoreType, usePool, enumsPool, structsPool)
			if finded != nil {
				founds = append(founds, *finded)
				continue
			}
			panic(fmt.Sprintf("not find useType %v", fCoreType))
		}
	}
	return founds
}

// 寻找useType对应enum的enum fields中在usetypes的匹配类型，如果没有在当前文件的所有structs中寻找
func getEnumFieldUseTypes(aim rust.UseType, usePool []rust.Use, enumsPool map[string]rust.Enum, structsPool map[string]rust.Struct) []rust.UseType {
	founds := []rust.UseType{}
	if _enum, ok := enumsPool[aim.Name]; ok {
		// 只过滤field里边包含的useType
		_p := _enum.Parse()

		for _, field := range _p.Fields {

			if !field.IsTumple() {
				continue
			}

			for _, p := range field.TupleParams {

				fCoreType := p.InnestCoreTypeName()

				logrus.WithFields(
					logrus.Fields{
						"aim":             aim,
						"field core type": fCoreType,
						"is base type":    rust.IsBaseType(fCoreType),
					},
				).Info("check is base type")
				if rust.IsBaseType(fCoreType) {
					continue
				}

				// 如果没有找到从同文件的enums中找enum
				finded := findFieldUseType(aim, fCoreType, usePool, enumsPool, structsPool)
				if finded != nil {
					founds = append(founds, *finded)
					continue
				}

				panic(fmt.Sprintf("not find useType %v", fCoreType))
			}

		}
	}
	return founds
}

func findFieldUseType(aim rust.UseType, fCoreType string, usePool []rust.Use, enumsPool map[string]rust.Enum, structsPool map[string]rust.Struct) *rust.UseType {
	for _, u := range usePool {
		uItems := u.Parse()

		for _, uItem := range uItems {
			if fCoreType == uItem.Name {
				return &uItem
			}
		}
	}

	_, ok1 := enumsPool[fCoreType]
	_, ok2 := structsPool[fCoreType]
	if ok1 || ok2 {
		fUseType := aim
		fUseType.Alias = fCoreType
		fUseType.Name = fCoreType
		return &fUseType
	}

	logger.WithFields(logrus.Fields{
		"aim":             aim,
		"field core type": fCoreType,
		"use pool":        usePool,
		"enums pool":      enumsPool,
		"structs pool":    structsPool,
	}).Panic("not find field useType")

	return nil
}

func getCachedUseTypes() []rust.UseType {
	useTypes := []rust.UseType{}
	for k := range usetype2Schema {
		useTypes = append(useTypes, rust.MustNewUseType(k))
	}
	return useTypes
}

// 生成 rust.TypeParsed 的 ref schema；方法参数、返回值、结构体字段、枚举项都使用该类型
// 从内生成ref，然后到外层层剥离，生成items，array 等描述
func genObjRefSchema(_type rust.TypeParsed, defaultModPath []string) *spec.Schema {
	s := spec.Schema{}

	if _type.Core == nil {
		ut := findUseType(_type.Name, getCachedUseTypes())
		if ut == nil && rust.IsBaseType(_type.Name) {
			tmp := rust.MustNewUseType(_type.Name)
			ut = &tmp
		}
		if ut == nil {
			ut = getSameLevelUseType(_type.Name, defaultModPath)
		}

		s = getUseTypeRefSchema(*ut)
		s.Title = _type.Name
	}

	if _type.IsOption {
		s := genObjRefSchema(*_type.Core, defaultModPath)
		s.Nullable = true
		return s
	}

	if _type.IsArray {
		s.Type = spec.StringOrArray{"array"}
		s.Items = &spec.SchemaOrArray{Schema: genObjRefSchema(*_type.Core, defaultModPath)}
	}

	if _type.IsVariadicValue {
		s.OneOf = append(s.OneOf, *genObjRefSchema(*_type.Core, defaultModPath))
		s.OneOf = append(s.OneOf, spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:  spec.StringOrArray{"array"},
				Items: &spec.SchemaOrArray{Schema: genObjRefSchema(*_type.Core, defaultModPath)},
			},
		})
		s.Nullable = true
	}
	return &s

}
