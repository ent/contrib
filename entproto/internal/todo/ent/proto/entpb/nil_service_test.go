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
	"testing"
	"time"

	"github.com/bionicstork/contrib/entproto/internal/todo/ent/enttest"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNilExampleService_Get(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewNilExampleService(client)
	ctx := context.Background()
	nex := client.NilExample.Create().
		SetStrNil("str").
		SetTimeNil(time.Now()).
		SaveX(ctx)
	get, err := svc.Get(ctx, &GetNilExampleRequest{Id: int32(nex.ID)})
	require.NoError(t, err)
	require.EqualValues(t, nex.TimeNil.Unix(), get.GetTimeNil().AsTime().Unix())
	require.EqualValues(t, *nex.StrNil, get.GetStrNil().GetValue())
}

func TestNilExampleService_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewNilExampleService(client)
	ctx := context.Background()
	ts := time.Now()
	c, err := svc.Create(ctx, &CreateNilExampleRequest{
		NilExample: &NilExample{
			StrNil:  wrapperspb.String("str"),
			TimeNil: timestamppb.New(ts),
		},
	})
	require.NoError(t, err)
	get, err := client.NilExample.Get(ctx, int(c.Id))
	require.NoError(t, err)
	require.EqualValues(t, "str", *get.StrNil)
	require.EqualValues(t, ts.Unix(), get.TimeNil.Unix())
}
