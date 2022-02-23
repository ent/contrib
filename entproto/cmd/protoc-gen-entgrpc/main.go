// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"embed"
	"entgo.io/ent/schema/field"
	"flag"
	"fmt"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/bionicstork/bionicstork/pkg/entproto"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	entSchemaPath     *string
	entSchemaTarget   *string
	idType   		  *string
	snake         = gen.Funcs["snake"].(func(string) string)
	status        = protogen.GoImportPath("google.golang.org/grpc/status")
	codes         = protogen.GoImportPath("google.golang.org/grpc/codes")
)

func main() {
	var flags flag.FlagSet
	entSchemaPath = flags.String("schema_path", "", "ent schema path")
	entSchemaTarget = flags.String("schema_target", "", "ent schema target (--target in ent cmdline)")
	idType = flags.String("idtype", "", "type of the tables primary key")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plg *protogen.Plugin) error {
		idConfigType := IDType(field.TypeInt)
		if *idType != "" {
			err := idConfigType.Set(*idType)
			if err != nil {
				log.Fatal(err)
			}
		}

		g, err := entc.LoadGraph(*entSchemaPath, &gen.Config{Target: *entSchemaTarget, IDType: &field.TypeInfo{Type: field.Type(idConfigType)}})
		if err != nil {
			return err
		}
		for _, f := range plg.Files {
			if !f.Generate {
				continue
			}
			if err := processFile(plg, f, g); err != nil {
				return err
			}
		}
		return nil
	})
}

// processFile generates service implementations from all services defined in the file.
func processFile(gen *protogen.Plugin, file *protogen.File, graph *gen.Graph) error {
	if len(file.Services) == 0 {
		return nil
	}
	for _, s := range file.Services {
		sg, err := newServiceGenerator(gen, file, graph, s)
		if err != nil {
			return err
		}
		if err := sg.generate(); err != nil {
			return err
		}
	}
	toProto, err := newToProtoGenerator(gen, file, graph)
	if err != nil {
		return err
	}
	if len(toProto.Nodes) > 0 {
		toProtoTmpl, err := toProto.Nodes[0].generateTemplate("to_proto")
		if err != nil {
			return err
		}
		if err := toProtoTmpl.ExecuteTemplate(toProto.Nodes[0], "to_proto_func", toProto); err != nil {
			return fmt.Errorf("template execution failed: %w", err)
		}
	}
	return nil
}

func newServiceGenerator(plugin *protogen.Plugin, file *protogen.File, graph *gen.Graph, service *protogen.Service) (*serviceGenerator, error) {
	adapter, err := entproto.LoadAdapter(graph, "")
	if err != nil {
		return nil, err
	}
	typ, err := extractEntTypeName(service, graph)
	if err != nil {
		return nil, err
	}
	filename := file.GeneratedFilenamePrefix + "_" + snake(service.GoName) + ".go"
	g := plugin.NewGeneratedFile(filename, file.GoImportPath)
	fieldMap, err := adapter.FieldMap(typ.Name)
	if err != nil {
		return nil, err
	}

	pkgPath, err := PkgPath(DefaultConfig, graph.Config.Target)
	if err != nil {
		log.Fatalln(err)
	}
	return &serviceGenerator{
		GeneratedFile: g,
		EntPackage:    protogen.GoImportPath(pkgPath),
		File:          file,
		Service:       service,
		EntType:       typ,
		FieldMap:      fieldMap,
	}, nil
}

func (g *serviceGenerator) generate() error {
	tmpl, err := g.generateTemplate("service")
	if err != nil {
		return err
	}
	if err := tmpl.ExecuteTemplate(g, "service", g); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}
	return nil
}

func (g *serviceGenerator) generateTemplate(name string) (*gen.Template, error) {
	tmpl, err := gen.NewTemplate(name).
		Funcs(template.FuncMap{
			"ident":        g.QualifiedGoIdent,
			"entIdent":     g.entIdent,
			"newConverter": g.newConverter,
			"unquote":      strconv.Unquote,
			"qualify": func(pkg, ident string) string {
				return g.QualifiedGoIdent(protogen.GoImportPath(pkg).Ident(ident))
			},
			"statusErr": func(code, msg string) string {
				return fmt.Sprintf("%s(%s, %q)",
					g.QualifiedGoIdent(status.Ident("Error")),
					g.QualifiedGoIdent(codes.Ident(code)),
					msg,
				)
			},
			"statusErrf": func(code, format string, args ...string) string {
				return fmt.Sprintf("%s(%s, %s, %s)",
					g.QualifiedGoIdent(status.Ident("Errorf")),
					g.QualifiedGoIdent(codes.Ident(code)),
					strconv.Quote(format),
					strings.Join(args, ","),
				)
			},
			"method": func(m *protogen.Method) *methodInput {
				return &methodInput{
					G:      g,
					Method: m,
				}
			},
		}).
		ParseFS(templates, "template/*.tmpl")
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func newToProtoGenerator(plugin *protogen.Plugin, file *protogen.File, graph *gen.Graph) (toProto toProtoGenerator, err error) {
	adapter, err := entproto.LoadAdapter(graph, "")
	if err != nil {
		return toProto, err
	}

	pkgPath, err := PkgPath(DefaultConfig, graph.Config.Target)
	if err != nil {
		log.Fatalln(err)
	}
	filename := file.GeneratedFilenamePrefix + "_to_proto.go"
	g := plugin.NewGeneratedFile(filename, file.GoImportPath)

	for _, nodeType := range graph.Nodes {
		msg, err := entproto.ExtractMessageAnnotation(nodeType)
		if err != nil || !msg.Generate {
			continue
		}
		fieldMap, err := adapter.FieldMap(nodeType.Name)
		if err != nil {
			return toProto, err
		}

		toProto.Nodes = append(toProto.Nodes, &serviceGenerator{
			GeneratedFile: g,
			File:          file,
			EntPackage:    protogen.GoImportPath(pkgPath),
			EntType:       nodeType,
			FieldMap:      fieldMap,
		})
	}
	return toProto, nil
}

type (
	serviceGenerator struct {
		*protogen.GeneratedFile
		EntPackage protogen.GoImportPath
		File       *protogen.File
		Service    *protogen.Service
		EntType    *gen.Type
		FieldMap   entproto.FieldMap
	}
	toProtoGenerator struct {
		Nodes []*serviceGenerator
	}
	methodInput struct {
		G      *serviceGenerator
		Method *protogen.Method
	}
)

//go:embed template/*
var templates embed.FS

func extractEntTypeName(s *protogen.Service, g *gen.Graph) (*gen.Type, error) {
	typeName := strings.TrimSuffix(s.GoName, "Service")
	for _, gt := range g.Nodes {
		if gt.Name == typeName {
			return gt, nil
		}
	}
	return nil, fmt.Errorf("entproto: type %q of service %q not found in graph", typeName, s.GoName)
}

func (g *serviceGenerator) entIdent(subpath string, ident string) protogen.GoIdent {
	ip := path.Join(string(g.EntPackage), subpath)
	return protogen.GoImportPath(ip).Ident(ident)
}

var DefaultConfig = &packages.Config{Mode: packages.NeedName}

// Note: Copied from ent/cmd/internal/base/packages.go
// PkgPath returns the Go package name for the given target path.
// Even if the existing path does not exist yet in the filesystem.
//
// If base.Config is nil, DefaultConfig will be used to load base.
func PkgPath(config *packages.Config, target string) (string, error) {
	if config == nil {
		config = DefaultConfig
	}
	pathCheck, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	var parts []string
	if _, err := os.Stat(pathCheck); os.IsNotExist(err) {
		parts = append(parts, filepath.Base(pathCheck))
		pathCheck = filepath.Dir(pathCheck)
	}
	// Try maximum 2 directories above the given
	// target to find the root package or module.
	for i := 0; i < 2; i++ {
		pkgs, err := packages.Load(config, pathCheck)
		if err != nil {
			return "", fmt.Errorf("load package info: %w", err)
		}
		if len(pkgs) == 0 || len(pkgs[0].Errors) != 0 {
			parts = append(parts, filepath.Base(pathCheck))
			pathCheck = filepath.Dir(pathCheck)
			continue
		}
		pkgPath := pkgs[0].PkgPath
		for j := len(parts) - 1; j >= 0; j-- {
			pkgPath = path.Join(pkgPath, parts[j])
		}
		return pkgPath, nil
	}
	return "", fmt.Errorf("root package or module was not found for: %s", target)
}

// IDType is a custom ID implementation for pflag.
type IDType field.Type

// Set implements the Set method of the flag.Value interface.
func (t *IDType) Set(s string) error {
	switch s {
	case field.TypeInt.String():
		*t = IDType(field.TypeInt)
	case field.TypeInt64.String():
		*t = IDType(field.TypeInt64)
	case field.TypeUint.String():
		*t = IDType(field.TypeUint)
	case field.TypeUint64.String():
		*t = IDType(field.TypeUint64)
	case field.TypeString.String():
		*t = IDType(field.TypeString)
	default:
		return fmt.Errorf("invalid type %q", s)
	}
	return nil
}
