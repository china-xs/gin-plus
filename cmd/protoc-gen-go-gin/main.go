// Package protoc_gen_gin
// @file      : main.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/7 17:21
// @Description: gin 自动文件
package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	showVersion = flag.Bool("version", false, "print the version and exit")
	omitempty   = flag.Bool("omitempty", true, "omit if google.api is empty")
)

const release = "v1.0.1"

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-gin %v\n", release)
		return
	}
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f, *omitempty)
		}
		return nil
	})
}

const deprecationComment = "// Deprecated: Do not use."
