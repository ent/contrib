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

package entoas

import (
	"os"
	"testing"

	"entgo.io/ent/entc/gen"
	"github.com/ogen-go/ogen"
	"github.com/stretchr/testify/require"
)

func TestExtension(t *testing.T) {
	t.Parallel()
	ex, err := NewExtension(
		DefaultPolicy(PolicyExpose),
		MinItemsPerPage(20),
		MaxItemsPerPage(40),
		Mutations(func(_ *gen.Graph, spec *ogen.Spec) error {
			spec.Info.
				SetTitle("Spec Title").
				SetDescription("Spec Description").
				SetVersion("Spec Version")
			return nil
		}),
		WriteTo(os.Stdout),
	)
	require.NoError(t, err)
	require.Equal(t, ex.config.DefaultPolicy, PolicyExpose)
	require.Len(t, ex.mutations, 1)
	require.Equal(t, os.Stdout, ex.out)
	require.Equal(t, int64(20), ex.config.MinItemsPerPage)
	require.Equal(t, int64(40), ex.config.MaxItemsPerPage)
}
