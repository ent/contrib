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

package entpb

import (
	"context"
	"github.com/bionicstork/contrib/entproto/internal/todo/ent/enttest"
	"github.com/bionicstork/contrib/entproto/internal/todo/ent/multiwordschema"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMultiWordSchemaService_Get(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewMultiWordSchemaService(client)
	ctx := context.Background()
	entry := client.MultiWordSchema.Create().
		SetUnit(multiwordschema.UnitFt).
		SaveX(ctx)
	get, err := svc.Get(ctx, &GetMultiWordSchemaRequest{Id: int32(entry.ID)})
	require.NoError(t, err)
	require.EqualValues(t, MultiWordSchema_FT, get.GetUnit())
}
