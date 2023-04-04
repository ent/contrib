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

package schemast

import (
	"bytes"
	"go/printer"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"golang.org/x/tools/imports"
)

// PrintOption modifies the behavior of Print.
type PrintOption func(opt *printOpts)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// Print writes the updated .go files from Context into path, the directory for the "schema" package in an
// ent project.  Print receives functional options of type PrintOption that modify its behavior.
func (c *Context) Print(path string, opts ...PrintOption) error {
	options := &printOpts{}
	for _, apply := range opts {
		apply(options)
	}
	for _, file := range c.syntax() {
		base := filepath.Base(c.SchemaPackage.Fset.File(file.Pos()).Name())
		var buf bytes.Buffer
		if err := printer.Fprint(&buf, c.SchemaPackage.Fset, file); err != nil {
			return err
		}
		fn := filepath.Join(path, base)
		t1 := time.Now()
		process, err := imports.Process(base, buf.Bytes(), nil)
		TimeTrack(t1, "imports.Process")
		if err != nil {
			return err
		}
		if options.headerComment != "" {
			if s := string(process); s != "" && options.commentRegexp.FindString(s) == "" {
				process = []byte(options.headerComment + "\n\n" + s)
			}
		}
		t2 := time.Now()
		if err := os.WriteFile(fn, process, 0600); err != nil {
			return err
		}
		TimeTrack(t2, "os.WriteFile")
	}
	return nil
}

// Header modifies Print to include a comment at the top of the printed .go files.
// If the file already contains the comment, even if it is not located at the very top of the file
// the comment will not be appended.
// Example:
//
//	ctx.Print("./schema", schemast.Header("File generated with ent-codegen-plugin.")
func Header(c string) PrintOption {
	return func(opt *printOpts) {
		opt.headerComment = "// " + c
		opt.commentRegexp = regexp.MustCompile("(?m)^" + opt.headerComment + "$")
	}
}

type printOpts struct {
	headerComment string
	commentRegexp *regexp.Regexp
}
