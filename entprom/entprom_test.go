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

package entprom

import (
	"context"
	"testing"

	"entgo.io/ent/examples/fs/ent"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestHookDefault(t *testing.T) {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)
	defer client.Close()
	ctx := context.Background()
	require.NoError(t, client.Schema.Create(ctx))

	ph := &hook{}
	client.Use(newHook(ph))
	client.File.Create().SetName("a8m").SaveX(ctx)

	require.Equal(t, 1, testutil.CollectAndCount(ph.opsProcessedTotal))
}

func TestExtraLabelsOption(t *testing.T) {
	ph := &hook{}
	expectedLabels := map[string]string{"dog": "fortuna"}
	Labels(expectedLabels)(ph)
	require.EqualValues(t, expectedLabels, ph.extraLabels)
}
