package plugin

import "github.com/gogo/protobuf/protoc-gen-gogo/generator"

type plugin struct{}

func NewPlugin(useGogoImport bool) generator.Plugin {
	return nil
}
