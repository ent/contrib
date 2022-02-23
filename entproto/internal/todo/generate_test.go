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

package todo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/bionicstork/bionicstork/pkg/entproto"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	tgt, err := ioutil.TempDir(os.TempDir(), "entproto-test-*")
	defer os.RemoveAll(tgt)
	require.NoError(t, err)
	graph, err := entc.LoadGraph("./ent/schema", &gen.Config{
		Target: tgt,
	})
	require.NoError(t, err)

	err = entproto.Generate(graph, "")
	require.NoError(t, err)

	bytes, err := ioutil.ReadFile(filepath.Join(tgt, "proto", "entpb", "entpb.proto"))
	require.NoError(t, err)
	require.True(t, strings.Contains(string(bytes), "// Code generated by entproto. DO NOT EDIT."))
}
