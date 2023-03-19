// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entproto

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"go.uber.org/multierr"
)

// Hook returns a gen.Hook that invokes Generate.
// To use it programatically:
//
//	entc.Generate("./ent/schema", &gen.Config{
//	  Hooks: []gen.Hook{
//	    entproto.Hook(),
//	  },
//	})
func Hook() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			// Because Generate has side effects (it is writing to the filesystem under gen.Config.Target),
			// we first run all generators, and only then invoke our code. This isn't great, and there's an
			// [open issue](https://github.com/ent/ent/issues/1311) to support this use-case better.
			err := next.Generate(g)
			if err != nil {
				return err
			}
			return Generate(g)
		})
	}
}

// Generate takes a *gen.Graph and creates .proto files. Next to each .proto file, Generate creates a generate.go
// file containing a //go:generate directive to invoke protoc and compile Go code from the protobuf definitions.
// If generate.go already exists next to the .proto file, this step is skipped.
func Generate(g *gen.Graph) error {
	entProtoDir := path.Join(g.Config.Target, "proto")
	adapter, err := LoadAdapter(g)
	if err != nil {
		return fmt.Errorf("entproto: failed parsing ent graph: %w", err)
	}
	var errs error
	for _, schema := range g.Schemas {
		name := schema.Name
		_, err := adapter.GetFileDescriptor(name)
		if err != nil && !errors.Is(err, ErrSchemaSkipped) {
			errs = multierr.Append(errs, err)
		}
	}
	if errs != nil {
		return fmt.Errorf("entproto: failed parsing some schemas: %w", errs)
	}
	allDescriptors := make([]*desc.FileDescriptor, 0, len(adapter.AllFileDescriptors()))
	for _, filedesc := range adapter.AllFileDescriptors() {
		allDescriptors = append(allDescriptors, filedesc)
	}

	// Print the .proto files.
	var printer protoprint.Printer
	if err = printer.PrintProtosToFileSystem(allDescriptors, entProtoDir); err != nil {
		return fmt.Errorf("entproto: failed writing .proto files: %w", err)
	}

	// Print a generate.go file with protoc command for go file generation
	for _, fd := range allDescriptors {
		protoFilePath := filepath.Join(entProtoDir, fd.GetName())
		dir := filepath.Dir(protoFilePath)
		genGoPath := filepath.Join(dir, "generate.go")
		if !fileExists(genGoPath) {
			contents := protocGenerateGo(fd)
			if err := os.WriteFile(genGoPath, []byte(contents), 0600); err != nil {
				return fmt.Errorf("entproto: failed generating generate.go file for %q: %w", protoFilePath, err)
			}
		}
	}
	return nil
}

func fileExists(fpath string) bool {
	if _, err := os.Stat(fpath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func protocGenerateGo(fd *desc.FileDescriptor) string {
	levelsUp := len(strings.Split(fd.GetPackage(), "."))
	toProtoBase := ""
	for i := 0; i < levelsUp; i++ {
		toProtoBase = filepath.Join("..", toProtoBase)
	}
	schemaDir := filepath.Join("..", toProtoBase, "schema")
	protocCmd := []string{
		"protoc",
		"-I=" + toProtoBase,
		"--go_out=" + toProtoBase,
		"--go-grpc_out=" + toProtoBase,
		"--go_opt=paths=source_relative",
		"--go-grpc_opt=paths=source_relative",
		"--entgrpc_out=" + toProtoBase,
		"--entgrpc_opt=paths=source_relative,schema_path=" + schemaDir,
		fd.GetName(),
	}
	goGen := fmt.Sprintf("//go:generate %s", strings.Join(protocCmd, " "))
	goPkgName := extractLastFqnPart(fd.GetPackage())
	return fmt.Sprintf("package %s\n%s\n", goPkgName, goGen)
}
