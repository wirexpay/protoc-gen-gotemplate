package main // import "github.com/wirexpay/protoc-gen-gotemplate"

import (
	pgghelpers "github.com/wirexpay/protoc-gen-gotemplate/helpers"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
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
