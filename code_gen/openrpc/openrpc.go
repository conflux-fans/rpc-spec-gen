package openrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	gconfig "github.com/Conflux-Chain/rpc-gen/config"
	"github.com/Conflux-Chain/rpc-gen/parser/rust"
	"github.com/Conflux-Chain/rpc-gen/parser/rust/config"
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
	blockNumS := spec.Schema{}
	blockNumS.Type = spec.StringOrArray{"string"}
	blockNumS.Pattern = `block number`
	usetype2Schema["BlockNumber"] = blockNumS
}

// TODO: gen schema for field type
func GenSchemaByStruct(structParsed rust.StructParsed, myUseType rust.UseType,
	structsInSameFile map[string]rust.Struct) spec.Schema {

	logger.Info("GenSchemaByStruct")

	s := spec.Schema{}
	s.Title = structParsed.Comment
	s.Type = spec.StringOrArray{"object"}
	s.Properties = make(map[string]spec.Schema, len(structParsed.Fields))

	for _, field := range structParsed.Fields {
		// 先从 usetype2Schema , basetypes 匹配
		// 若都未匹配到则设置为与struct同级别usetype
		ut := findUseType(field.Type.Name, getCacheUseTypes())
		if ut == nil && config.IsBaseType(field.Type.Name) {
			tmp := rust.MustNewUseType(field.Type.Name)
			ut = &tmp
		}
		if ut == nil {
			ut = getSameLevelUseType(field.Type.Name, myUseType)
		}
		refSchema := getSchemaRef(*ut)
		refSchema.Title = field.Name
		s.Properties[field.Name] = refSchema
	}

	return s
}

func GenSchemaByEnum(enumParsed rust.EnumParsed, myUseType rust.UseType,
	structsInSameFile map[string]rust.Struct, useTypes []rust.UseType) spec.Schema {

	if _, ok := usetype2Schema[enumParsed.Name]; ok {
		return usetype2Schema[enumParsed.Name]
	}

	s := spec.Schema{}
	s.Title = enumParsed.Comment
	s.Type = spec.StringOrArray{"string"}
	s.Properties = make(map[string]spec.Schema, len(enumParsed.Fields))

	// TODO: should set enum values
	// for _, field := range enumParsed.Fields {
	// }

	return s
}

func GenSchemas(useTypes []rust.UseType) map[string]spec.Schema {

	logger.WithField("useTypes", useTypes).Debug("GenSchemas")

	for _, useType := range useTypes {
		// get code file path by usetype
		if _, ok := usetype2Schema[useType.String()]; ok {
			continue
		}

		if config.IgnoredUseTypes[useType.String()] {
			continue
		}

		meta := config.GetUseTypeMeta(useType)

		if meta == nil {
			panic(fmt.Sprintf("not found meta of %v", useType.String()))
		}

		if meta.IsBaseType() {
			continue
		}

		// 从配置中查看类型所在文件
		// 解析所在文件，拿到所有structs和use types
		// 判断structs中是否有该类型，有则生成schema
		//     先生成字段的schema
		//     再生成该struct的schema
		code, e := ioutil.ReadFile(meta.InFilePath())
		if e != nil {
			logrus.
				WithField("useType", useType).
				WithField("filePath", meta.InFilePath()).
				WithError(e).
				Panic("read file error")
		}
		// s, us := rust.FindStruct(string(code), useType.Name)

		// 获得code中的所有struct，以struct名为key，struct解析结果为value
		structs, us := rust.GetStructs(string(code))
		if _strcut, ok := structs[useType.Name]; ok {
			fieldUsetypes := getStructFieldUseTypes(useType, structs, us)
			logger.WithFields(logrus.Fields{
				"struct in useType": useType,
				"filtered us":       fieldUsetypes,
			}).Debug("filter struct field using use types")
			GenSchemas(fieldUsetypes)

			logger.Info("sdfsdfsdfsdfsdfsdf")

			usetype2Schema[useType.String()] = GenSchemaByStruct(_strcut.Parse(), useType, structs)
			continue
		}

		enums, us := rust.GetEnums(string(code))
		if _, ok := enums[useType.Name]; ok {
			// for _, u := range us {
			// 	useTypes = append(useTypes, u.Parse()...)
			// }
			// TODO: 只过滤field里边包含的useType
			// _us := []rust.UseType{}
			// for _, u := range us {
			// 	_us = append(_us, u.Parse()...)
			// }
			// logger.WithFields(logrus.Fields{
			// 	"useType": useType,
			// 	"us":      _us,
			// }).Debug("")
			// GenSchemas(_us)
			usetype2Schema[useType.String()] = GenSchemaByEnum(enums[useType.Name].Parse(), useType, structs, useTypes)
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
		// if _, err := os.Stat(p); os.IsNotExist(err) {
		// 	folder := path.Join(p, "../")
		// 	os.MkdirAll(folder, 0700)
		// }

		// e := ioutil.WriteFile(p, j, 0644)
		// if e != nil {
		// 	logger.Panic(e)
		// }
	}
}

// - 方法的参数和返回值都是schema的ref，放到methods.params和methods.result
// - 查看是否有已生成的scheme，没有则创建
func GenMethod(funcParsed rust.FuncParsed, useTypes []rust.UseType) Method {
	var method Method
	method.Summary = funcParsed.Comment
	method.Name = funcParsed.RpcMethod
	method.Params = make([]*ContentDescriptor, len(funcParsed.Params))
	for i, param := range funcParsed.Params {
		ut := findUseType(param.Type.Name, useTypes)
		if ut == nil {
			logger.WithFields(logrus.Fields{
				"Name":      param.Type.Name,
				"Use Types": useTypes,
			}).Error("not found use type")
		}

		method.Params[i] = &ContentDescriptor{
			Content: Content{
				Name:     param.Name,
				Required: !param.Type.IsOption,
				Schema:   getSchemaRef(*ut),
			},
		}
	}

	ut := findUseType(funcParsed.Return.Type.Name, useTypes)

	if ut == nil {
		logger.WithFields(logrus.Fields{
			"Name":      funcParsed.Return.Type.Name,
			"Use Types": useTypes,
		}).Error("not found use type")
	}

	method.Result = &ContentDescriptor{
		Content: Content{
			Name:     funcParsed.Return.Name,
			Required: !funcParsed.Return.Type.IsOption,
			Schema:   getSchemaRef(*ut),
		},
	}
	return method
}

func GenMethods(traitParsed rust.TraitParsed, useTypes []rust.UseType) []Method {
	var methods []Method
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
	docPath := path.Join(gconfig.GetConfig().DocTemplateRootPath, space, doc.Info.Title+".json")
	j, _ := json.MarshalIndent(doc, "", " ")
	saveFile(docPath, j)
}

func SaveDoc(doc OpenRPCSpec1, space string) {
	docPath := path.Join(gconfig.GetConfig().DocRootPath, space, doc.Info.Title+".json")
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

func findUseType(name string, useTypes []rust.UseType) *rust.UseType {
	for _, useType := range useTypes {
		if useType.Name == name {
			return &useType
		}
	}
	return nil

	// logrus.WithField("name", name).WithField("use types", useTypes).Panic("not found useType")
	// panic("not find")
}

func mustFindUseType(name string, useTypes []rust.UseType) *rust.UseType {
	ut := findUseType(name, useTypes)
	if ut == nil {
		panic("not found useType")
	}
	return ut
}

func getSameLevelUseType(name string, myUsetType rust.UseType) *rust.UseType {
	ut := myUsetType
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
func getStructFieldUseTypes(aim rust.UseType, structsPool map[string]rust.Struct, usePool []rust.Use) []rust.UseType {
	_us := []rust.UseType{}
	if _struct, ok := structsPool[aim.Name]; ok {
		// TODO: 只过滤field里边包含的useType
		_p := _struct.Parse()

		for _, field := range _p.Fields {

			fCoreType := field.Type.InnestCoreTypeName()

			if fCoreType == "NodeLockStatus" {
				logrus.Info("NodeLockStatus")
			}

			finded := false
			for _, u := range usePool {
				uItems := u.Parse()

				for _, uItem := range uItems {
					if fCoreType == uItem.Name {
						finded = true
						_us = append(_us, uItem)
						break
					}
				}
			}

			if finded {
				continue
			}

			// 如果没有找到从同文件的strcuts中找sturct
			if _, ok := structsPool[fCoreType]; ok {
				fUseType := aim
				fUseType.Name = fCoreType
				_us = append(_us, fUseType)
				continue
			}

			logrus.WithFields(
				logrus.Fields{
					"aim":             aim,
					"field core type": fCoreType,
					"is base type":    config.IsBaseType(fCoreType),
				},
			).Info("check is base type")
			if config.IsBaseType(fCoreType) {
				continue
			}

			panic(fmt.Sprintf("not find useType %v", fCoreType))

		}
	}
	return _us
}

func getCacheUseTypes() []rust.UseType {
	useTypes := []rust.UseType{}
	for k := range usetype2Schema {
		useTypes = append(useTypes, rust.MustNewUseType(k))
	}
	return useTypes
}
