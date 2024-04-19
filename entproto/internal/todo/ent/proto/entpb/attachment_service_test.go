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

package entpb

import (
	"context"
	"strconv"
	"testing"
	"time"

	"entgo.io/contrib/entproto/internal/todo/ent"
	"entgo.io/contrib/entproto/internal/todo/ent/enttest"
	"entgo.io/contrib/entproto/internal/todo/ent/user"
	"github.com/google/uuid"
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
	require.NoError(t, err)
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

func TestAttachmentService_MultiEdge(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewAttachmentService(client)
	ctx := context.Background()
	var users []*ent.User
	for i := 0; i < 5; i++ {
		users = append(users, client.User.Create().
			SetUserName(strconv.Itoa(i)).
			SetJoined(time.Now()).
			SetPoints(10).
			SetExp(1000).
			SetStatus("pending").
			SetExternalID(i+1).
			SetCrmID(uuid.New()).
			SetCustomPb(1).
			SetLabels(nil).
			SetOmitPrefix(user.OmitPrefixFoo).
			SetMimeType(user.MimeTypeSvg).
			SaveX(ctx))
	}
	att, err := svc.Create(ctx, &CreateAttachmentRequest{Attachment: &Attachment{
		User: &User{
			Id: users[0].ID,
		},
		Recipients: []*User{
			{Id: users[1].ID},
			{Id: users[2].ID},
			{Id: users[3].ID},
			{Id: users[4].ID},
		},
	}})
	all := client.Attachment.Query().WithRecipients(func(query *ent.UserQuery) {
		query.Select(user.FieldID)
	}).AllX(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.Len(t, all[0].Edges.Recipients, 4)

	get, err := svc.Get(ctx, &GetAttachmentRequest{
		Id:   att.GetId(),
		View: GetAttachmentRequest_WITH_EDGE_IDS,
	})
	require.NoError(t, err)
	require.NotNil(t, get.User)
	require.Len(t, get.Recipients, 4)
}
