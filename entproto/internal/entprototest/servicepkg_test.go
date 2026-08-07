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

package entprototest

import (
	"path"
	"strings"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/suite"
)

type ServicePkgTestSuite struct {
	suite.Suite
	schema, pkg, protoFile string
}

func (suite *ServicePkgTestSuite) SetupSuite() {
	suite.DirExists(suite.schema)
	suite.FileExists(suite.protoFile)
}

func TestServicePkgSuite(t *testing.T) {
	suite.Run(t, &ServicePkgTestSuite{
		schema:    "../todo/ent/schema",
		pkg:       "entgo.io/contrib/entproto/internal/entprototest/db",
		protoFile: "./db/proto/entpb/entpb.proto",
	})
}

func (suite *ServicePkgTestSuite) TestDefaultPackage() {
	graph, err := entc.LoadGraph(suite.schema, &gen.Config{})
	suite.Require().NoError(err)

	suite.processImports(false, graph)
}

func (suite *ServicePkgTestSuite) TestCustomPackage() {
	graph, err := entc.LoadGraph(suite.schema, &gen.Config{Package: suite.pkg})
	suite.Require().NoError(err)

	suite.processImports(true, graph)
}

func (suite *ServicePkgTestSuite) processImports(shouldMatch bool, graph *gen.Graph) bool {
	fdesc, err := protoparse.Parser{}.ParseFiles(suite.protoFile)
	suite.Require().NoError(err)

	var (
		match    func(value bool, msgAndArgs ...interface{}) bool
		b, m     bool
		typeName string
		entType  *gen.Type
	)

	if shouldMatch == true {
		match = suite.True
	} else {
		match = suite.False
	}

	for _, fd := range fdesc {
		for _, sd := range fd.GetFile().GetServices() {
			entType = nil
			typeName = strings.TrimSuffix(sd.GetName(), "Service")

			for _, gt := range graph.Nodes {
				if gt.Name == typeName {
					entType = gt
					break
				}
			}
			if !suite.NotNil(entType) {
				continue
			}

			for _, si := range entType.SiblingImports() {
				m, err = path.Match(suite.pkg+"/*", si.Path)
				if !suite.NoError(err) || !match(m, "type %s has unexpected import: %q", typeName, si.Path) {
					b = true
				}
			}
		}
	}

	return !b
}
