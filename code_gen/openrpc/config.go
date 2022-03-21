package openrpc

import (
	"path"
	"regexp"
	"strings"

	"github.com/Conflux-Chain/rpc-gen/config"
	"github.com/Conflux-Chain/rpc-gen/parser/rust"
	"github.com/go-openapi/spec"
)

const (
	schemaRefRoot = "#/components/schemas/"
)

func getSchemaSaveRelativePath(space string, schemaFullName string) string {
	// return path.Join(strings.Join(useType.ModPath, "/"), useType.Name+".json")
	return path.Join(space, strings.Join(strings.Split(schemaFullName, "::"), "/")+".json")
}

func getSchemaSavePath(space string, schemaFullName string) string {
	return path.Join(config.GetConfig().SchemaRootPath, getSchemaSaveRelativePath(space, schemaFullName))
}

func getUseTypeRefSchema(useType rust.UseType) spec.Schema {
	s := spec.Schema{}
	// schemaName := strings.ReplaceAll(getSchemaSaveRelativePath(useType), "/", "_")

	// TODO: 如果ModdPath为空，则直接使用useType.Name

	schemaName := strings.Join(useType.ModPath, "__") + "__" + useType.Name
	s.Ref = spec.MustCreateRef(schemaRefRoot + schemaName)
	return s
}

func parseSchemaRefToUseType(ref string) rust.UseType {

	fullUseType := strings.TrimPrefix(ref, schemaRefRoot)
	matchs := regexp.MustCompile(`(.*)__(.*)`).FindStringSubmatch(fullUseType)

	if len(matchs) != 3 {
		logger.WithField("ref", ref).Debug("parse to a base type")
		return rust.UseType{
			Name: fullUseType,
		}
	}

	fullName := matchs[1]

	return rust.UseType{
		ModPath: strings.Split(fullName, "__"),
		Name:    matchs[2],
	}
}
