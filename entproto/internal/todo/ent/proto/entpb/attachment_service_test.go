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
