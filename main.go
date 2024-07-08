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
	//g := generator.New()
	//
	//data, err := ioutil.ReadAll(os.Stdin)
	//if err != nil {
	//	g.Error(err, "reading input")
	//}
	//
	//if err = proto.Unmarshal(data, g.Request); err != nil {
	//	g.Error(err, "parsing input proto")
	//}
	//
	//if len(g.Request.FileToGenerate) == 0 {
	//	g.Fail("no files to generate")
	//}
	//
	//g.CommandLineParameters(g.Request.GetParameter())
	//
	//pgghelpers.ParseParams(g)
	//
	//// Generate the protobufs
	//g.GenerateAllFiles()
	//
	//data, err = proto.Marshal(g.Response)
	//if err != nil {
	//	g.Error(err, "failed to marshal output proto")
	//}
	//
	//_, err = os.Stdout.Write(data)
	//if err != nil {
	//	g.Error(err, "failed to write output proto")
	//}
}
