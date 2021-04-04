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

	"entgo.io/contrib/entproto/internal/todo/ent/enttest"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAttachmentService_Get(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewAttachmentService(client)

	ctx := context.Background()
	attachment := client.Attachment.Create().SaveX(ctx)
	id, err := attachment.ID.MarshalBinary()
	require.NoError(t, err)

	get, err := svc.Get(ctx, &GetAttachmentRequest{Id: id})
	require.EqualValues(t, get.Id, id)
	respStatus, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.OK)

	get, err = svc.Get(ctx, &GetAttachmentRequest{Id: []byte("short")})
	require.Nil(t, get)
	respStatus, ok = status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.InvalidArgument)
}
