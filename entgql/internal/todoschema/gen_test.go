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

package todoschema_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCollection_M2MTotalCountSchemaConfig asserts that the generated M2M
// totalCount loader schema-qualifies its join table. This project enables the
// sql/schemaconfig feature, so every "joinT := sql.Table(...)" in the totalCount
// loaders must be followed by a "joinT.Schema(...)" call; otherwise the join
// would be emitted unqualified and break against a schema-aware database.
func TestCollection_M2MTotalCountSchemaConfig(t *testing.T) {
	out, err := os.ReadFile("./ent/gql_collection.go")
	require.NoError(t, err)
	src := string(out)

	require.Contains(t, src, "joinT := sql.Table(group.UsersTable)")
	// The join table is qualified with the schema configured for the M2M
	// relation so it matches the rest of the schema-qualified query.
	require.Contains(t, src, "joinT.Schema(gq.schemaConfig.GroupUsers)")
	// Every M2M totalCount loader must qualify its join table; none should be
	// left as a bare sql.Table without a following Schema call.
	for _, block := range strings.SplitAfter(src, "joinT := sql.Table(") {
		if i := strings.Index(block, "s.Join(joinT)"); i >= 0 {
			require.Contains(t, block[:i], "joinT.Schema(",
				"unqualified M2M join table in totalCount loader:\n%s", block[:i])
		}
	}
}
