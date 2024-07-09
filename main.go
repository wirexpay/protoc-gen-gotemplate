package main // import "moul.io/protoc-gen-gotemplate"

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
	pgghelpers "moul.io/protoc-gen-gotemplate/helpers"
)

func main() {
	protogen.Options{
		ParamFunc: pgghelpers.Flags.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}

			pgghelpers.ParseParams(plugin, file)
		}

		return nil
	})
}
