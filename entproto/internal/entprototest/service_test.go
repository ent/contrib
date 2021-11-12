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
	"fmt"
)

type method struct {
	name         string
	input        string
	output       string
	formatInput  bool
	formatOutput bool
}

func (m *method) getFormatInput(name string) string {
	return fmt.Sprintf(m.input, name)
}
func (m *method) getFormatOutput(name string) string {
	return fmt.Sprintf(m.output, name)
}

var (
	methodCreate = method{
		name:         "Create",
		input:        "Create%sRequest",
		output:       "%s",
		formatInput:  true,
		formatOutput: true,
	}

	methodGet = method{
		name:         "Get",
		input:        "Get%sRequest",
		output:       "%s",
		formatInput:  true,
		formatOutput: true,
	}

	methodUpdate = method{
		name:         "Update",
		input:        "Update%sRequest",
		output:       "%s",
		formatInput:  true,
		formatOutput: true,
	}

	methodDelete = method{
		name:         "Delete",
		input:        "Delete%sRequest",
		output:       "Empty",
		formatInput:  true,
		formatOutput: false,
	}
)

func (suite *AdapterTestSuite) TestServiceGeneration() {
	testCases := []struct {
		testName        string
		schemaName      string
		includedMethods []method
		excludedMethods []method
	}{
		{
			testName:   "Default Method Generation",
			schemaName: "BlogPost",
			includedMethods: []method{
				methodCreate,
				methodGet,
				methodUpdate,
				methodDelete,
			},
		},
		{
			testName:   "All Methods Generation",
			schemaName: "AllMethodsService",
			includedMethods: []method{
				methodCreate,
				methodGet,
				methodUpdate,
				methodDelete,
			},
		},
		{
			testName:   "One Method Generation",
			schemaName: "OneMethodService",
			includedMethods: []method{
				methodGet,
			},
			excludedMethods: []method{
				methodCreate,
				methodUpdate,
				methodDelete,
			},
		},
		{
			testName:   "Two Method Generation",
			schemaName: "TwoMethodService",
			includedMethods: []method{
				methodCreate,
				methodGet,
			},
			excludedMethods: []method{
				methodUpdate,
				methodDelete,
			},
		},
	}

	for _, tc := range testCases {
		println(fmt.Sprintf("Test %s", tc.testName))

		fd, err := suite.adapter.GetFileDescriptor(tc.schemaName)
		suite.Require().NoError(err)

		svc := fd.FindService(fmt.Sprintf("entpb.%sService", tc.schemaName))
		suite.NotNil(svc)

		for _, m := range tc.includedMethods {
			getMeth := svc.FindMethodByName(m.name)
			suite.Require().NotNil(getMeth)

			if m.formatInput {
				m.input = m.getFormatInput(tc.schemaName)
			}
			suite.EqualValues(m.input, getMeth.GetInputType().GetName())
			if m.formatOutput {
				m.output = m.getFormatOutput(tc.schemaName)
			}
			suite.EqualValues(m.output, getMeth.GetOutputType().GetName())
		}

		for _, m := range tc.excludedMethods {
			getMeth := svc.FindMethodByName(m.name)
			suite.Nil(getMeth)
		}

		println("PASS")
	}
}
