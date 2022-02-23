// Code generated by protoc-gen-entgrpc. DO NOT EDIT.
package entpb

import (
	context "context"
	base64 "encoding/base64"
	entproto "github.com/bionicstork/bionicstork/pkg/entproto"
	ent "github.com/bionicstork/bionicstork/pkg/entproto/internal/todo/ent"
	multiwordschema "github.com/bionicstork/bionicstork/pkg/entproto/internal/todo/ent/multiwordschema"
	sqlgraph "entgo.io/ent/dialect/sql/sqlgraph"
	fmt "fmt"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	strconv "strconv"
	strings "strings"
)

// MultiWordSchemaService implements MultiWordSchemaServiceServer
type MultiWordSchemaService struct {
	client *ent.Client
	UnimplementedMultiWordSchemaServiceServer
}

// NewMultiWordSchemaService returns a new MultiWordSchemaService
func NewMultiWordSchemaService(client *ent.Client) *MultiWordSchemaService {
	return &MultiWordSchemaService{
		client: client,
	}
}

func toProtoMultiWordSchema_Unit(e multiwordschema.Unit) MultiWordSchema_Unit {
	if v, ok := MultiWordSchema_Unit_value[strings.ToUpper(string(e))]; ok {
		return MultiWordSchema_Unit(v)
	}
	return MultiWordSchema_Unit(0)
}

func toEntMultiWordSchema_Unit(e MultiWordSchema_Unit) multiwordschema.Unit {
	if v, ok := MultiWordSchema_Unit_name[int32(e)]; ok {
		return multiwordschema.Unit(strings.ToLower(v))
	}
	return ""
}

// toProtoMultiWordSchema transforms the ent type to the pb type
func toProtoMultiWordSchema(e *ent.MultiWordSchema) (*MultiWordSchema, error) {
	v := &MultiWordSchema{}
	id := int32(e.ID)
	v.Id = id
	unit := toProtoMultiWordSchema_Unit(e.Unit)
	v.Unit = unit
	return v, nil
}

// Create implements MultiWordSchemaServiceServer.Create
func (svc *MultiWordSchemaService) Create(ctx context.Context, req *CreateMultiWordSchemaRequest) (*MultiWordSchema, error) {
	multiwordschema := req.GetMultiWordSchema()
	m := svc.client.MultiWordSchema.Create()
	multiwordschemaUnit := toEntMultiWordSchema_Unit(multiwordschema.GetUnit())
	m.SetUnit(multiwordschemaUnit)
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoMultiWordSchema(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Get implements MultiWordSchemaServiceServer.Get
func (svc *MultiWordSchemaService) Get(ctx context.Context, req *GetMultiWordSchemaRequest) (*MultiWordSchema, error) {
	var (
		err error
		get *ent.MultiWordSchema
	)
	id := int(req.GetId())
	switch req.GetView() {
	case GetMultiWordSchemaRequest_VIEW_UNSPECIFIED, GetMultiWordSchemaRequest_BASIC:
		get, err = svc.client.MultiWordSchema.Get(ctx, id)
	case GetMultiWordSchemaRequest_WITH_EDGE_IDS:
		get, err = svc.client.MultiWordSchema.Query().
			Where(multiwordschema.ID(id)).
			Only(ctx)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProtoMultiWordSchema(get)
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}
	return nil, nil

}

// Update implements MultiWordSchemaServiceServer.Update
func (svc *MultiWordSchemaService) Update(ctx context.Context, req *UpdateMultiWordSchemaRequest) (*MultiWordSchema, error) {
	multiwordschema := req.GetMultiWordSchema()
	multiwordschemaID := int(multiwordschema.GetId())
	m := svc.client.MultiWordSchema.UpdateOneID(multiwordschemaID)
	multiwordschemaUnit := toEntMultiWordSchema_Unit(multiwordschema.GetUnit())
	m.SetUnit(multiwordschemaUnit)
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoMultiWordSchema(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Delete implements MultiWordSchemaServiceServer.Delete
func (svc *MultiWordSchemaService) Delete(ctx context.Context, req *DeleteMultiWordSchemaRequest) (*emptypb.Empty, error) {
	var err error
	id := int(req.GetId())
	err = svc.client.MultiWordSchema.DeleteOneID(id).Exec(ctx)
	switch {
	case err == nil:
		return &emptypb.Empty{}, nil
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// List implements MultiWordSchemaServiceServer.List
func (svc *MultiWordSchemaService) List(ctx context.Context, req *ListMultiWordSchemaRequest) (*ListMultiWordSchemaResponse, error) {
	var (
		err      error
		entList  []*ent.MultiWordSchema
		pageSize int
	)
	pageSize = int(req.GetPageSize())
	switch {
	case pageSize < 0:
		return nil, status.Errorf(codes.InvalidArgument, "page size cannot be less than zero")
	case pageSize == 0 || pageSize > entproto.MaxPageSize:
		pageSize = entproto.MaxPageSize
	}
	listQuery := svc.client.MultiWordSchema.Query().
		Order(ent.Desc(multiwordschema.FieldID)).
		Limit(pageSize + 1)
	if req.GetPageToken() != "" {
		bytes, err := base64.StdEncoding.DecodeString(req.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "page token is invalid")
		}
		token, err := strconv.ParseInt(string(bytes), 10, 32)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "page token is invalid")
		}
		pageToken := int(token)
		listQuery = listQuery.
			Where(multiwordschema.IDLTE(pageToken))
	}
	switch req.GetView() {
	case ListMultiWordSchemaRequest_VIEW_UNSPECIFIED, ListMultiWordSchemaRequest_BASIC:
		entList, err = listQuery.All(ctx)
	case ListMultiWordSchemaRequest_WITH_EDGE_IDS:
		entList, err = listQuery.
			All(ctx)
	}
	switch {
	case err == nil:
		var nextPageToken string
		if len(entList) == pageSize+1 {
			nextPageToken = base64.StdEncoding.EncodeToString(
				[]byte(fmt.Sprintf("%v", entList[len(entList)-1].ID)))
			entList = entList[:len(entList)-1]
		}
		var pbList []*MultiWordSchema
		for _, entEntity := range entList {
			pbEntity, err := toProtoMultiWordSchema(entEntity)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "internal error: %s", err)
			}
			pbList = append(pbList, pbEntity)
		}
		return &ListMultiWordSchemaResponse{
			MultiWordSchemaList: pbList,
			NextPageToken:       nextPageToken,
		}, nil
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}
