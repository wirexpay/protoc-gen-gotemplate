package pgghelpers

import (
	"flag"
	"fmt"
	"strings"

	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	ggdescriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"google.golang.org/protobuf/compiler/protogen"
)

type Parameters struct {
	TemplateDir       string
	DestinationDir    string
	Files             string
	Debug             bool
	All               bool
	SinglePackageMode bool
	FileMode          bool
}

var (
	Flags  flag.FlagSet
	params = Parameters{}
)

func init() {
	Flags.StringVar(&params.TemplateDir, "template_dir", "", "path to look for templates")
	Flags.StringVar(&params.DestinationDir, "destination_dir", "", "base path to write output")
	Flags.StringVar(&params.Files, "files", "", "files to process, with extensions")
	Flags.BoolVar(&params.Debug, "debug", false, "if 'true', `protoc` will generate a more verbose output")
	Flags.BoolVar(&params.All, "all", false, "if 'true', protobuf files without `Service` will also be parsed")
	Flags.BoolVar(&params.SinglePackageMode, "single_package_mode", false, "if 'true', `protoc` won't accept multiple packages to be compiled at once ('!= from `all`'), but will support `Message` lookup across the imported protobuf dependencies")
	Flags.BoolVar(&params.FileMode, "file_mode", false, "")
}

func ParseParams(plugin *protogen.Plugin, file *protogen.File) {
	tmplMap := make(map[string]*plugin_go.CodeGeneratorResponse_File)
	concatOrAppend := func(f *plugin_go.CodeGeneratorResponse_File) {
		if val, ok := tmplMap[f.GetName()]; ok {
			*val.Content += f.GetContent()
		} else {
			tmplMap[f.GetName()] = f
			g := plugin.NewGeneratedFile(f.GetName(), file.GoImportPath)
			g.P(f.GetContent())
		}
	}

	allFiles := strings.Split(params.Files, " ")
	allowedFile := false
	if len(allFiles) > 0 {
		for _, f := range allFiles {
			if strings.Contains(file.Proto.GetName(), f) {
				allowedFile = true
				break
			}
		}
	} else {
		allowedFile = true
	}

	if !allowedFile {
		return
	}

	if params.SinglePackageMode {
		registry = ggdescriptor.NewRegistry()
		SetRegistry(registry)
		if err := registry.Load(plugin.Request); err != nil {
			plugin.Error(fmt.Errorf("registry: failed to load the request: %w", err))
		}
	}

	// Generate the encoders
	if params.All {
		if params.SinglePackageMode {
			if _, err := registry.LookupFile(file.Proto.GetName()); err != nil {
				plugin.Error(fmt.Errorf("registry: failed to lookup file %q: %w", file.Proto.GetName(), err))
			}
		}
		encoder := NewGenericTemplateBasedEncoder(params.TemplateDir, file, params.Debug, params.DestinationDir)
		for _, tmpl := range encoder.Files() {
			concatOrAppend(tmpl)
		}

		return
	}

	if params.FileMode {
		if s := file.Proto.GetService(); len(s) > 0 {
			encoder := NewGenericTemplateBasedEncoder(params.TemplateDir, file, params.Debug, params.DestinationDir)
			for _, tmpl := range encoder.Files() {
				concatOrAppend(tmpl)
			}
		}

		return
	}

	encoder := NewGenericTemplateBasedEncoder(params.TemplateDir, file, params.Debug, params.DestinationDir)
	for _, tmpl := range encoder.Files() {
		concatOrAppend(tmpl)
	}
}
