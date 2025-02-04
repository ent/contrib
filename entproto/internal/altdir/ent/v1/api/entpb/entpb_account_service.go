// Code generated by protoc-gen-entgrpc. DO NOT EDIT.
package entpb

import (
	context "context"
	base64 "encoding/base64"
	entproto "entgo.io/contrib/entproto"
	ent "entgo.io/contrib/entproto/internal/altdir/ent"
	account "entgo.io/contrib/entproto/internal/altdir/ent/account"
	user "entgo.io/contrib/entproto/internal/altdir/ent/user"
	sqlgraph "entgo.io/ent/dialect/sql/sqlgraph"
	fmt "fmt"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	strconv "strconv"
)

// AccountService implements AccountServiceServer
type AccountService struct {
	client *ent.Client
	UnimplementedAccountServiceServer
}

// NewAccountService returns a new AccountService
func NewAccountService(client *ent.Client) *AccountService {
	return &AccountService{
		client: client,
	}
}

// toProtoAccount transforms the ent type to the pb type
func toProtoAccount(e *ent.Account) (*Account, error) {
	v := &Account{}
	id := int64(e.ID)
	v.Id = id
	if edg := e.Edges.Owner; edg != nil {
		id := int64(edg.ID)
		v.Owner = &User{
			Id: id,
		}
	}
	return v, nil
}

// toProtoAccountList transforms a list of ent type to a list of pb type
func toProtoAccountList(e []*ent.Account) ([]*Account, error) {
	var pbList []*Account
	for _, entEntity := range e {
		pbEntity, err := toProtoAccount(entEntity)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		pbList = append(pbList, pbEntity)
	}
	return pbList, nil
}

// Create implements AccountServiceServer.Create
func (svc *AccountService) Create(ctx context.Context, req *CreateAccountRequest) (*Account, error) {
	account := req.GetAccount()
	m, err := svc.createBuilder(account)
	if err != nil {
		return nil, err
	}
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoAccount(res)
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

// Get implements AccountServiceServer.Get
func (svc *AccountService) Get(ctx context.Context, req *GetAccountRequest) (*Account, error) {
	var (
		err error
		get *ent.Account
	)
	id := int(req.GetId())
	switch req.GetView() {
	case GetAccountRequest_VIEW_UNSPECIFIED, GetAccountRequest_BASIC:
		get, err = svc.client.Account.Get(ctx, id)
	case GetAccountRequest_WITH_EDGE_IDS:
		get, err = svc.client.Account.Query().
			Where(account.ID(id)).
			WithOwner(func(query *ent.UserQuery) {
				query.Select(user.FieldID)
			}).
			Only(ctx)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProtoAccount(get)
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Update implements AccountServiceServer.Update
func (svc *AccountService) Update(ctx context.Context, req *UpdateAccountRequest) (*Account, error) {
	account := req.GetAccount()
	accountID := int(account.GetId())
	m := svc.client.Account.UpdateOneID(accountID)

	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoAccount(res)
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

// Delete implements AccountServiceServer.Delete
func (svc *AccountService) Delete(ctx context.Context, req *DeleteAccountRequest) (*emptypb.Empty, error) {
	var err error
	id := int(req.GetId())
	err = svc.client.Account.DeleteOneID(id).Exec(ctx)
	switch {
	case err == nil:
		return &emptypb.Empty{}, nil
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// List implements AccountServiceServer.List
func (svc *AccountService) List(ctx context.Context, req *ListAccountRequest) (*ListAccountResponse, error) {
	var (
		err      error
		entList  []*ent.Account
		pageSize int
	)
	pageSize = int(req.GetPageSize())
	switch {
	case pageSize < 0:
		return nil, status.Errorf(codes.InvalidArgument, "page size cannot be less than zero")
	case pageSize == 0 || pageSize > entproto.MaxPageSize:
		pageSize = entproto.MaxPageSize
	}
	listQuery := svc.client.Account.Query().
		Order(ent.Desc(account.FieldID)).
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
			Where(account.IDLTE(pageToken))
	}
	switch req.GetView() {
	case ListAccountRequest_VIEW_UNSPECIFIED, ListAccountRequest_BASIC:
		entList, err = listQuery.All(ctx)
	case ListAccountRequest_WITH_EDGE_IDS:
		entList, err = listQuery.
			WithOwner(func(query *ent.UserQuery) {
				query.Select(user.FieldID)
			}).
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
		protoList, err := toProtoAccountList(entList)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return &ListAccountResponse{
			AccountList:   protoList,
			NextPageToken: nextPageToken,
		}, nil
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// BatchCreate implements AccountServiceServer.BatchCreate
func (svc *AccountService) BatchCreate(ctx context.Context, req *BatchCreateAccountsRequest) (*BatchCreateAccountsResponse, error) {
	requests := req.GetRequests()
	if len(requests) > entproto.MaxBatchCreateSize {
		return nil, status.Errorf(codes.InvalidArgument, "batch size cannot be greater than %d", entproto.MaxBatchCreateSize)
	}
	bulk := make([]*ent.AccountCreate, len(requests))
	for i, req := range requests {
		account := req.GetAccount()
		var err error
		bulk[i], err = svc.createBuilder(account)
		if err != nil {
			return nil, err
		}
	}
	res, err := svc.client.Account.CreateBulk(bulk...).Save(ctx)
	switch {
	case err == nil:
		protoList, err := toProtoAccountList(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return &BatchCreateAccountsResponse{
			Accounts: protoList,
		}, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

func (svc *AccountService) createBuilder(account *Account) (*ent.AccountCreate, error) {
	m := svc.client.Account.Create()
	if account.GetOwner() != nil {
		accountOwner := int(account.GetOwner().GetId())
		m.SetOwnerID(accountOwner)
	}
	return m, nil
}
