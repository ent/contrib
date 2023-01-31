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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	load, err := Load("./internal/loadtest/ent/schema")
	require.NoError(t, err)

	fd, ok := load.lookupMethod("Message", "Fields")
	require.True(t, ok)
	require.EqualValues(t, fd.Name.Name, "Fields")
	require.True(t, load.HasType("Message"))
	require.False(t, load.HasType("MessageXX"))

	buf := bytes.Buffer{}
	err = printer.Fprint(&buf, load.SchemaPackage.Fset, fd)
	require.NoError(t, err)
	require.EqualValues(t, `// Fields of the Message.
func (Message) Fields() []ent.Field {
	return nil
}`, buf.String())
}
