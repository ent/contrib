// Code generated by ent, DO NOT EDIT.

package user

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "user_id"
	// FieldUserName holds the string denoting the user_name field in the database.
	FieldUserName = "user_name"
	// FieldJoined holds the string denoting the joined field in the database.
	FieldJoined = "joined"
	// FieldPoints holds the string denoting the points field in the database.
	FieldPoints = "points"
	// FieldExp holds the string denoting the exp field in the database.
	FieldExp = "exp"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldExternalID holds the string denoting the external_id field in the database.
	FieldExternalID = "external_id"
	// FieldCrmID holds the string denoting the crm_id field in the database.
	FieldCrmID = "crm_id"
	// FieldBanned holds the string denoting the banned field in the database.
	FieldBanned = "banned"
	// FieldCustomPb holds the string denoting the custom_pb field in the database.
	FieldCustomPb = "custom_pb"
	// FieldOptNum holds the string denoting the opt_num field in the database.
	FieldOptNum = "opt_num"
	// FieldOptStr holds the string denoting the opt_str field in the database.
	FieldOptStr = "opt_str"
	// FieldOptBool holds the string denoting the opt_bool field in the database.
	FieldOptBool = "opt_bool"
	// FieldBigInt holds the string denoting the big_int field in the database.
	FieldBigInt = "big_int"
	// FieldBUser1 holds the string denoting the b_user_1 field in the database.
	FieldBUser1 = "b_user_1"
	// FieldHeightInCm holds the string denoting the height_in_cm field in the database.
	FieldHeightInCm = "height_in_cm"
	// FieldAccountBalance holds the string denoting the account_balance field in the database.
	FieldAccountBalance = "account_balance"
	// FieldUnnecessary holds the string denoting the unnecessary field in the database.
	FieldUnnecessary = "unnecessary"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldLabels holds the string denoting the labels field in the database.
	FieldLabels = "labels"
	// FieldDeviceType holds the string denoting the device_type field in the database.
	FieldDeviceType = "device_type"
	// FieldOmitPrefix holds the string denoting the omit_prefix field in the database.
	FieldOmitPrefix = "omit_prefix"
	// FieldMimeType holds the string denoting the mime_type field in the database.
	FieldMimeType = "mime_type"
	// EdgeGroup holds the string denoting the group edge name in mutations.
	EdgeGroup = "group"
	// EdgeAttachment holds the string denoting the attachment edge name in mutations.
	EdgeAttachment = "attachment"
	// EdgeReceived1 holds the string denoting the received_1 edge name in mutations.
	EdgeReceived1 = "received_1"
	// EdgePet holds the string denoting the pet edge name in mutations.
	EdgePet = "pet"
	// EdgeSkipEdge holds the string denoting the skip_edge edge name in mutations.
	EdgeSkipEdge = "skip_edge"
	// GroupFieldID holds the string denoting the ID field of the Group.
	GroupFieldID = "id"
	// AttachmentFieldID holds the string denoting the ID field of the Attachment.
	AttachmentFieldID = "id"
	// PetFieldID holds the string denoting the ID field of the Pet.
	PetFieldID = "id"
	// SkipEdgeExampleFieldID holds the string denoting the ID field of the SkipEdgeExample.
	SkipEdgeExampleFieldID = "id"
	// Table holds the table name of the user in the database.
	Table = "users"
	// GroupTable is the table that holds the group relation/edge.
	GroupTable = "users"
	// GroupInverseTable is the table name for the Group entity.
	// It exists in this package in order to avoid circular dependency with the "group" package.
	GroupInverseTable = "groups"
	// GroupColumn is the table column denoting the group relation/edge.
	GroupColumn = "user_group"
	// AttachmentTable is the table that holds the attachment relation/edge.
	AttachmentTable = "attachments"
	// AttachmentInverseTable is the table name for the Attachment entity.
	// It exists in this package in order to avoid circular dependency with the "attachment" package.
	AttachmentInverseTable = "attachments"
	// AttachmentColumn is the table column denoting the attachment relation/edge.
	AttachmentColumn = "user_attachment"
	// Received1Table is the table that holds the received_1 relation/edge. The primary key declared below.
	Received1Table = "attachment_recipients"
	// Received1InverseTable is the table name for the Attachment entity.
	// It exists in this package in order to avoid circular dependency with the "attachment" package.
	Received1InverseTable = "attachments"
	// PetTable is the table that holds the pet relation/edge.
	PetTable = "pets"
	// PetInverseTable is the table name for the Pet entity.
	// It exists in this package in order to avoid circular dependency with the "pet" package.
	PetInverseTable = "pets"
	// PetColumn is the table column denoting the pet relation/edge.
	PetColumn = "user_pet"
	// SkipEdgeTable is the table that holds the skip_edge relation/edge.
	SkipEdgeTable = "skip_edge_examples"
	// SkipEdgeInverseTable is the table name for the SkipEdgeExample entity.
	// It exists in this package in order to avoid circular dependency with the "skipedgeexample" package.
	SkipEdgeInverseTable = "skip_edge_examples"
	// SkipEdgeColumn is the table column denoting the skip_edge relation/edge.
	SkipEdgeColumn = "user_skip_edge"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldUserName,
	FieldJoined,
	FieldPoints,
	FieldExp,
	FieldStatus,
	FieldExternalID,
	FieldCrmID,
	FieldBanned,
	FieldCustomPb,
	FieldOptNum,
	FieldOptStr,
	FieldOptBool,
	FieldBigInt,
	FieldBUser1,
	FieldHeightInCm,
	FieldAccountBalance,
	FieldUnnecessary,
	FieldType,
	FieldLabels,
	FieldDeviceType,
	FieldOmitPrefix,
	FieldMimeType,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "users"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_group",
}

var (
	// Received1PrimaryKey and Received1Column2 are the table columns denoting the
	// primary key for the received_1 relation (M2M).
	Received1PrimaryKey = []string{"attachment_id", "user_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultBanned holds the default value on creation for the "banned" field.
	DefaultBanned bool
	// DefaultHeightInCm holds the default value on creation for the "height_in_cm" field.
	DefaultHeightInCm float32
	// DefaultAccountBalance holds the default value on creation for the "account_balance" field.
	DefaultAccountBalance float64
)

// Status defines the type for the "status" enum field.
type Status string

// Status values.
const (
	StatusPending Status = "pending"
	StatusActive  Status = "active"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusPending, StatusActive:
		return nil
	default:
		return fmt.Errorf("user: invalid enum value for status field: %q", s)
	}
}

// DeviceType defines the type for the "device_type" enum field.
type DeviceType string

// DeviceTypeGLOWY9000 is the default value of the DeviceType enum.
const DefaultDeviceType = DeviceTypeGLOWY9000

// DeviceType values.
const (
	DeviceTypeGLOWY9000 DeviceType = "GLOWY9000"
	DeviceTypeSPEEDY300 DeviceType = "SPEEDY300"
)

func (dt DeviceType) String() string {
	return string(dt)
}

// DeviceTypeValidator is a validator for the "device_type" field enum values. It is called by the builders before save.
func DeviceTypeValidator(dt DeviceType) error {
	switch dt {
	case DeviceTypeGLOWY9000, DeviceTypeSPEEDY300:
		return nil
	default:
		return fmt.Errorf("user: invalid enum value for device_type field: %q", dt)
	}
}

// OmitPrefix defines the type for the "omit_prefix" enum field.
type OmitPrefix string

// OmitPrefix values.
const (
	OmitPrefixFoo OmitPrefix = "foo"
	OmitPrefixBar OmitPrefix = "bar"
)

func (op OmitPrefix) String() string {
	return string(op)
}

// OmitPrefixValidator is a validator for the "omit_prefix" field enum values. It is called by the builders before save.
func OmitPrefixValidator(op OmitPrefix) error {
	switch op {
	case OmitPrefixFoo, OmitPrefixBar:
		return nil
	default:
		return fmt.Errorf("user: invalid enum value for omit_prefix field: %q", op)
	}
}

// MimeType defines the type for the "mime_type" enum field.
type MimeType string

// MimeType values.
const (
	MimeTypePng MimeType = "image/png"
	MimeTypeSvg MimeType = "image/xml+svg"
)

func (mt MimeType) String() string {
	return string(mt)
}

// MimeTypeValidator is a validator for the "mime_type" field enum values. It is called by the builders before save.
func MimeTypeValidator(mt MimeType) error {
	switch mt {
	case MimeTypePng, MimeTypeSvg:
		return nil
	default:
		return fmt.Errorf("user: invalid enum value for mime_type field: %q", mt)
	}
}

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByUserName orders the results by the user_name field.
func ByUserName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUserName, opts...).ToFunc()
}

// ByJoined orders the results by the joined field.
func ByJoined(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldJoined, opts...).ToFunc()
}

// ByPoints orders the results by the points field.
func ByPoints(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPoints, opts...).ToFunc()
}

// ByExp orders the results by the exp field.
func ByExp(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExp, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByExternalID orders the results by the external_id field.
func ByExternalID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExternalID, opts...).ToFunc()
}

// ByCrmID orders the results by the crm_id field.
func ByCrmID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCrmID, opts...).ToFunc()
}

// ByBanned orders the results by the banned field.
func ByBanned(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBanned, opts...).ToFunc()
}

// ByCustomPb orders the results by the custom_pb field.
func ByCustomPb(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomPb, opts...).ToFunc()
}

// ByOptNum orders the results by the opt_num field.
func ByOptNum(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOptNum, opts...).ToFunc()
}

// ByOptStr orders the results by the opt_str field.
func ByOptStr(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOptStr, opts...).ToFunc()
}

// ByOptBool orders the results by the opt_bool field.
func ByOptBool(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOptBool, opts...).ToFunc()
}

// ByBigInt orders the results by the big_int field.
func ByBigInt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBigInt, opts...).ToFunc()
}

// ByBUser1 orders the results by the b_user_1 field.
func ByBUser1(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBUser1, opts...).ToFunc()
}

// ByHeightInCm orders the results by the height_in_cm field.
func ByHeightInCm(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHeightInCm, opts...).ToFunc()
}

// ByAccountBalance orders the results by the account_balance field.
func ByAccountBalance(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAccountBalance, opts...).ToFunc()
}

// ByUnnecessary orders the results by the unnecessary field.
func ByUnnecessary(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUnnecessary, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByDeviceType orders the results by the device_type field.
func ByDeviceType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeviceType, opts...).ToFunc()
}

// ByOmitPrefix orders the results by the omit_prefix field.
func ByOmitPrefix(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOmitPrefix, opts...).ToFunc()
}

// ByMimeType orders the results by the mime_type field.
func ByMimeType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMimeType, opts...).ToFunc()
}

// ByGroupField orders the results by group field.
func ByGroupField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newGroupStep(), sql.OrderByField(field, opts...))
	}
}

// ByAttachmentField orders the results by attachment field.
func ByAttachmentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAttachmentStep(), sql.OrderByField(field, opts...))
	}
}

// ByReceived1Count orders the results by received_1 count.
func ByReceived1Count(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newReceived1Step(), opts...)
	}
}

// ByReceived1 orders the results by received_1 terms.
func ByReceived1(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newReceived1Step(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByPetField orders the results by pet field.
func ByPetField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPetStep(), sql.OrderByField(field, opts...))
	}
}

// BySkipEdgeField orders the results by skip_edge field.
func BySkipEdgeField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSkipEdgeStep(), sql.OrderByField(field, opts...))
	}
}
func newGroupStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(GroupInverseTable, GroupFieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, GroupTable, GroupColumn),
	)
}
func newAttachmentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AttachmentInverseTable, AttachmentFieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, AttachmentTable, AttachmentColumn),
	)
}
func newReceived1Step() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Received1InverseTable, AttachmentFieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, Received1Table, Received1PrimaryKey...),
	)
}
func newPetStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PetInverseTable, PetFieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, PetTable, PetColumn),
	)
}
func newSkipEdgeStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SkipEdgeInverseTable, SkipEdgeExampleFieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, SkipEdgeTable, SkipEdgeColumn),
	)
}
