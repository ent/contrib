// Code generated by protoc-gen-entgrpc. DO NOT EDIT.
package entpb

import (
	context "context"
	base64 "encoding/base64"
	entproto "github.com/bionicstork/contrib/entproto"
	ent "github.com/bionicstork/contrib/entproto/internal/todo/ent"
	attachment "github.com/bionicstork/contrib/entproto/internal/todo/ent/attachment"
	group "github.com/bionicstork/contrib/entproto/internal/todo/ent/group"
	pet "github.com/bionicstork/contrib/entproto/internal/todo/ent/pet"
	schema "github.com/bionicstork/contrib/entproto/internal/todo/ent/schema"
	user "github.com/bionicstork/contrib/entproto/internal/todo/ent/user"
	runtime "github.com/bionicstork/contrib/entproto/runtime"
	sqlgraph "entgo.io/ent/dialect/sql/sqlgraph"
	errors "errors"
	fmt "fmt"
	uuid "github.com/google/uuid"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	strconv "strconv"
	strings "strings"
)

// UserService implements UserServiceServer
type UserService struct {
	client *ent.Client
	UnimplementedUserServiceServer
}

// NewUserService returns a new UserService
func NewUserService(client *ent.Client) *UserService {
	return &UserService{
		client: client,
	}
}

func toProtoUser_Status(e user.Status) User_Status {
	if v, ok := User_Status_value[strings.ToUpper(string(e))]; ok {
		return User_Status(v)
	}
	return User_Status(0)
}

func toEntUser_Status(e User_Status) user.Status {
	if v, ok := User_Status_name[int32(e)]; ok {
		return user.Status(strings.ToLower(v))
	}
	return ""
}

// toProtoUser transforms the ent type to the pb type
func toProtoUser(e *ent.User) (*User, error) {
	v := &User{}
	accountbalance := e.AccountBalance
	v.AccountBalance = accountbalance
	buser1 := wrapperspb.Int32(int32(e.BUser1))
	v.BUser_1 = buser1
	banned := e.Banned
	v.Banned = banned
	bigintValue, err := e.BigInt.Value()
	if err != nil {
		return nil, err
	}
	bigintTyped, ok := bigintValue.(string)
	if !ok {
		return nil, errors.New("casting value to string")
	}
	bigint := wrapperspb.String(bigintTyped)
	v.BigInt = bigint
	crmid, err := e.CrmID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	v.CrmId = crmid
	custompb := uint64(e.CustomPb)
	v.CustomPb = custompb
	exp := e.Exp
	v.Exp = exp
	externalid := int32(e.ExternalID)
	v.ExternalId = externalid
	heightincm := e.HeightInCm
	v.HeightInCm = heightincm
	id := int32(e.ID)
	v.Id = id
	joined := timestamppb.New(e.Joined)
	v.Joined = joined
	optbool := wrapperspb.Bool(e.OptBool)
	v.OptBool = optbool
	optnum := wrapperspb.Int32(int32(e.OptNum))
	v.OptNum = optnum
	optstr := wrapperspb.String(e.OptStr)
	v.OptStr = optstr
	points := uint32(e.Points)
	v.Points = points
	status := toProtoUser_Status(e.Status)
	v.Status = status
	username := e.UserName
	v.UserName = username
	if edg := e.Edges.Attachment; edg != nil {
		id, err := edg.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		v.Attachment = &Attachment{
			Id: id,
		}
	}
	if edg := e.Edges.Group; edg != nil {
		id := int32(edg.ID)
		v.Group = &Group{
			Id: id,
		}
	}
	if edg := e.Edges.Pet; edg != nil {
		id := int32(edg.ID)
		v.Pet = &Pet{
			Id: id,
		}
	}
	for _, edg := range e.Edges.Received1 {
		id, err := edg.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		v.Received_1 = append(v.Received_1, &Attachment{
			Id: id,
		})
	}
	return v, nil
}

// Create implements UserServiceServer.Create
func (svc *UserService) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	user := req.GetUser()
	m := svc.client.User.Create()
	userAccountBalance := float64(user.GetAccountBalance())
	m.SetAccountBalance(userAccountBalance)
	if user.GetBUser_1() != nil {
		userBUser1 := int(user.GetBUser_1().GetValue())
		m.SetBUser1(userBUser1)
	}
	userBanned := user.GetBanned()
	m.SetBanned(userBanned)
	if user.GetBigInt() != nil {
		userBigInt := schema.BigInt{}
		if err := (&userBigInt).Scan(user.GetBigInt().GetValue()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetBigInt(userBigInt)
	}
	var userCrmID uuid.UUID
	if err := (&userCrmID).UnmarshalBinary(user.GetCrmId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}
	m.SetCrmID(userCrmID)
	userCustomPb := uint8(user.GetCustomPb())
	m.SetCustomPb(userCustomPb)
	userExp := uint64(user.GetExp())
	m.SetExp(userExp)
	userExternalID := int(user.GetExternalId())
	m.SetExternalID(userExternalID)
	userHeightInCm := float32(user.GetHeightInCm())
	m.SetHeightInCm(userHeightInCm)
	userJoined := runtime.ExtractTime(user.GetJoined())
	m.SetJoined(userJoined)
	if user.GetOptBool() != nil {
		userOptBool := user.GetOptBool().GetValue()
		m.SetOptBool(userOptBool)
	}
	if user.GetOptNum() != nil {
		userOptNum := int(user.GetOptNum().GetValue())
		m.SetOptNum(userOptNum)
	}
	if user.GetOptStr() != nil {
		userOptStr := user.GetOptStr().GetValue()
		m.SetOptStr(userOptStr)
	}
	userPoints := uint(user.GetPoints())
	m.SetPoints(userPoints)
	userStatus := toEntUser_Status(user.GetStatus())
	m.SetStatus(userStatus)
	userUserName := user.GetUserName()
	m.SetUserName(userUserName)
	if user.GetAttachment() != nil {
		var userAttachment uuid.UUID
		if err := (&userAttachment).UnmarshalBinary(user.GetAttachment().GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetAttachmentID(userAttachment)
	}
	if user.GetGroup() != nil {
		userGroup := int(user.GetGroup().GetId())
		m.SetGroupID(userGroup)
	}
	if user.GetPet() != nil {
		userPet := int(user.GetPet().GetId())
		m.SetPetID(userPet)
	}
	for _, item := range user.GetReceived_1() {
		var received1 uuid.UUID
		if err := (&received1).UnmarshalBinary(item.GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.AddReceived1IDs(received1)
	}
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoUser(res)
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

// Get implements UserServiceServer.Get
func (svc *UserService) Get(ctx context.Context, req *GetUserRequest) (*User, error) {
	var (
		err error
		get *ent.User
	)
	id := int(req.GetId())
	switch req.GetView() {
	case GetUserRequest_VIEW_UNSPECIFIED, GetUserRequest_BASIC:
		get, err = svc.client.User.Get(ctx, id)
	case GetUserRequest_WITH_EDGE_IDS:
		get, err = svc.client.User.Query().
			Where(user.ID(id)).
			WithAttachment(func(query *ent.AttachmentQuery) {
				query.Select(attachment.FieldID)
			}).
			WithGroup(func(query *ent.GroupQuery) {
				query.Select(group.FieldID)
			}).
			WithPet(func(query *ent.PetQuery) {
				query.Select(pet.FieldID)
			}).
			WithReceived1(func(query *ent.AttachmentQuery) {
				query.Select(attachment.FieldID)
			}).
			Only(ctx)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProtoUser(get)
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}
	return nil, nil

}

// Update implements UserServiceServer.Update
func (svc *UserService) Update(ctx context.Context, req *UpdateUserRequest) (*User, error) {
	user := req.GetUser()
	userID := int(user.GetId())
	m := svc.client.User.UpdateOneID(userID)
	userAccountBalance := float64(user.GetAccountBalance())
	m.SetAccountBalance(userAccountBalance)
	if user.GetBUser_1() != nil {
		userBUser1 := int(user.GetBUser_1().GetValue())
		m.SetBUser1(userBUser1)
	}
	userBanned := user.GetBanned()
	m.SetBanned(userBanned)
	if user.GetBigInt() != nil {
		userBigInt := schema.BigInt{}
		if err := (&userBigInt).Scan(user.GetBigInt().GetValue()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetBigInt(userBigInt)
	}
	var userCrmID uuid.UUID
	if err := (&userCrmID).UnmarshalBinary(user.GetCrmId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}
	m.SetCrmID(userCrmID)
	userCustomPb := uint8(user.GetCustomPb())
	m.SetCustomPb(userCustomPb)
	userExp := uint64(user.GetExp())
	m.SetExp(userExp)
	userExternalID := int(user.GetExternalId())
	m.SetExternalID(userExternalID)
	userHeightInCm := float32(user.GetHeightInCm())
	m.SetHeightInCm(userHeightInCm)
	if user.GetOptBool() != nil {
		userOptBool := user.GetOptBool().GetValue()
		m.SetOptBool(userOptBool)
	}
	if user.GetOptNum() != nil {
		userOptNum := int(user.GetOptNum().GetValue())
		m.SetOptNum(userOptNum)
	}
	if user.GetOptStr() != nil {
		userOptStr := user.GetOptStr().GetValue()
		m.SetOptStr(userOptStr)
	}
	userPoints := uint(user.GetPoints())
	m.SetPoints(userPoints)
	userStatus := toEntUser_Status(user.GetStatus())
	m.SetStatus(userStatus)
	userUserName := user.GetUserName()
	m.SetUserName(userUserName)
	if user.GetAttachment() != nil {
		var userAttachment uuid.UUID
		if err := (&userAttachment).UnmarshalBinary(user.GetAttachment().GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetAttachmentID(userAttachment)
	}
	if user.GetGroup() != nil {
		userGroup := int(user.GetGroup().GetId())
		m.SetGroupID(userGroup)
	}
	if user.GetPet() != nil {
		userPet := int(user.GetPet().GetId())
		m.SetPetID(userPet)
	}
	for _, item := range user.GetReceived_1() {
		var received1 uuid.UUID
		if err := (&received1).UnmarshalBinary(item.GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.AddReceived1IDs(received1)
	}
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoUser(res)
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

// Delete implements UserServiceServer.Delete
func (svc *UserService) Delete(ctx context.Context, req *DeleteUserRequest) (*emptypb.Empty, error) {
	var err error
	id := int(req.GetId())
	err = svc.client.User.DeleteOneID(id).Exec(ctx)
	switch {
	case err == nil:
		return &emptypb.Empty{}, nil
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// List implements UserServiceServer.List
func (svc *UserService) List(ctx context.Context, req *ListUserRequest) (*ListUserResponse, error) {
	var (
		err      error
		entList  []*ent.User
		pageSize int
	)
	pageSize = int(req.GetPageSize())
	switch {
	case pageSize < 0:
		return nil, status.Errorf(codes.InvalidArgument, "page size cannot be less than zero")
	case pageSize == 0 || pageSize > entproto.MaxPageSize:
		pageSize = entproto.MaxPageSize
	}
	listQuery := svc.client.User.Query().
		Order(ent.Desc(user.FieldID)).
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
			Where(user.IDLTE(pageToken))
	}
	switch req.GetView() {
	case ListUserRequest_VIEW_UNSPECIFIED, ListUserRequest_BASIC:
		entList, err = listQuery.All(ctx)
	case ListUserRequest_WITH_EDGE_IDS:
		entList, err = listQuery.
			WithAttachment(func(query *ent.AttachmentQuery) {
				query.Select(attachment.FieldID)
			}).
			WithGroup(func(query *ent.GroupQuery) {
				query.Select(group.FieldID)
			}).
			WithPet(func(query *ent.PetQuery) {
				query.Select(pet.FieldID)
			}).
			WithReceived1(func(query *ent.AttachmentQuery) {
				query.Select(attachment.FieldID)
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
		var pbList []*User
		for _, entEntity := range entList {
			pbEntity, err := toProtoUser(entEntity)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "internal error: %s", err)
			}
			pbList = append(pbList, pbEntity)
		}
		return &ListUserResponse{
			UserList:      pbList,
			NextPageToken: nextPageToken,
		}, nil
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}
