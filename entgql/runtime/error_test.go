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

package runtime_test

import (
	"testing"

	"entgo.io/contrib/entgql/runtime"
	"github.com/stretchr/testify/require"
)

func TestErrNodeNotFound(t *testing.T) {
	t.Parallel()
	err := runtime.ErrNodeNotFound(42)
	require.EqualError(t, err, "input: Could not resolve to a node with the global id of '42'")
	require.Equal(t, "NOT_FOUND", err.Extensions["code"])
}
