// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/contrib/entproto/internal/todo/ent/attachment"
	"entgo.io/contrib/entproto/internal/todo/ent/group"
	"entgo.io/contrib/entproto/internal/todo/ent/pet"
	"entgo.io/contrib/entproto/internal/todo/ent/schema"
	"entgo.io/contrib/entproto/internal/todo/ent/skipedgeexample"
	"entgo.io/contrib/entproto/internal/todo/ent/user"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// UserName holds the value of the "user_name" field.
	UserName string `json:"user_name,omitempty"`
	// Joined holds the value of the "joined" field.
	Joined time.Time `json:"joined,omitempty"`
	// Points holds the value of the "points" field.
	Points uint `json:"points,omitempty"`
	// Exp holds the value of the "exp" field.
	Exp uint64 `json:"exp,omitempty"`
	// Status holds the value of the "status" field.
	Status user.Status `json:"status,omitempty"`
	// ExternalID holds the value of the "external_id" field.
	ExternalID int `json:"external_id,omitempty"`
	// CrmID holds the value of the "crm_id" field.
	CrmID uuid.UUID `json:"crm_id,omitempty"`
	// Banned holds the value of the "banned" field.
	Banned bool `json:"banned,omitempty"`
	// CustomPb holds the value of the "custom_pb" field.
	CustomPb uint8 `json:"custom_pb,omitempty"`
	// OptNum holds the value of the "opt_num" field.
	OptNum int `json:"opt_num,omitempty"`
	// OptStr holds the value of the "opt_str" field.
	OptStr string `json:"opt_str,omitempty"`
	// OptBool holds the value of the "opt_bool" field.
	OptBool bool `json:"opt_bool,omitempty"`
	// BigInt holds the value of the "big_int" field.
	BigInt schema.BigInt `json:"big_int,omitempty"`
	// BUser1 holds the value of the "b_user_1" field.
	BUser1 int `json:"b_user_1,omitempty"`
	// HeightInCm holds the value of the "height_in_cm" field.
	HeightInCm float32 `json:"height_in_cm,omitempty"`
	// AccountBalance holds the value of the "account_balance" field.
	AccountBalance float64 `json:"account_balance,omitempty"`
	// Unnecessary holds the value of the "unnecessary" field.
	Unnecessary string `json:"unnecessary,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels []string `json:"labels,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserQuery when eager-loading is set.
	Edges      UserEdges `json:"edges"`
	user_group *int
}

// UserEdges holds the relations/edges for other nodes in the graph.
type UserEdges struct {
	// Group holds the value of the group edge.
	Group *Group `json:"group,omitempty"`
	// Attachment holds the value of the attachment edge.
	Attachment *Attachment `json:"attachment,omitempty"`
	// Received1 holds the value of the received_1 edge.
	Received1 []*Attachment `json:"received_1,omitempty"`
	// Pet holds the value of the pet edge.
	Pet *Pet `json:"pet,omitempty"`
	// SkipEdge holds the value of the skip_edge edge.
	SkipEdge *SkipEdgeExample `json:"skip_edge,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// GroupOrErr returns the Group value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserEdges) GroupOrErr() (*Group, error) {
	if e.loadedTypes[0] {
		if e.Group == nil {
			// The edge group was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: group.Label}
		}
		return e.Group, nil
	}
	return nil, &NotLoadedError{edge: "group"}
}

// AttachmentOrErr returns the Attachment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserEdges) AttachmentOrErr() (*Attachment, error) {
	if e.loadedTypes[1] {
		if e.Attachment == nil {
			// The edge attachment was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: attachment.Label}
		}
		return e.Attachment, nil
	}
	return nil, &NotLoadedError{edge: "attachment"}
}

// Received1OrErr returns the Received1 value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) Received1OrErr() ([]*Attachment, error) {
	if e.loadedTypes[2] {
		return e.Received1, nil
	}
	return nil, &NotLoadedError{edge: "received_1"}
}

// PetOrErr returns the Pet value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserEdges) PetOrErr() (*Pet, error) {
	if e.loadedTypes[3] {
		if e.Pet == nil {
			// The edge pet was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: pet.Label}
		}
		return e.Pet, nil
	}
	return nil, &NotLoadedError{edge: "pet"}
}

// SkipEdgeOrErr returns the SkipEdge value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserEdges) SkipEdgeOrErr() (*SkipEdgeExample, error) {
	if e.loadedTypes[4] {
		if e.SkipEdge == nil {
			// The edge skip_edge was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: skipedgeexample.Label}
		}
		return e.SkipEdge, nil
	}
	return nil, &NotLoadedError{edge: "skip_edge"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case user.FieldLabels:
			values[i] = new([]byte)
		case user.FieldBigInt:
			values[i] = new(schema.BigInt)
		case user.FieldBanned, user.FieldOptBool:
			values[i] = new(sql.NullBool)
		case user.FieldHeightInCm, user.FieldAccountBalance:
			values[i] = new(sql.NullFloat64)
		case user.FieldID, user.FieldPoints, user.FieldExp, user.FieldExternalID, user.FieldCustomPb, user.FieldOptNum, user.FieldBUser1:
			values[i] = new(sql.NullInt64)
		case user.FieldUserName, user.FieldStatus, user.FieldOptStr, user.FieldUnnecessary, user.FieldType:
			values[i] = new(sql.NullString)
		case user.FieldJoined:
			values[i] = new(sql.NullTime)
		case user.FieldCrmID:
			values[i] = new(uuid.UUID)
		case user.ForeignKeys[0]: // user_group
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type User", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case user.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			u.ID = int(value.Int64)
		case user.FieldUserName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field user_name", values[i])
			} else if value.Valid {
				u.UserName = value.String
			}
		case user.FieldJoined:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field joined", values[i])
			} else if value.Valid {
				u.Joined = value.Time
			}
		case user.FieldPoints:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field points", values[i])
			} else if value.Valid {
				u.Points = uint(value.Int64)
			}
		case user.FieldExp:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field exp", values[i])
			} else if value.Valid {
				u.Exp = uint64(value.Int64)
			}
		case user.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				u.Status = user.Status(value.String)
			}
		case user.FieldExternalID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field external_id", values[i])
			} else if value.Valid {
				u.ExternalID = int(value.Int64)
			}
		case user.FieldCrmID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field crm_id", values[i])
			} else if value != nil {
				u.CrmID = *value
			}
		case user.FieldBanned:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field banned", values[i])
			} else if value.Valid {
				u.Banned = value.Bool
			}
		case user.FieldCustomPb:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field custom_pb", values[i])
			} else if value.Valid {
				u.CustomPb = uint8(value.Int64)
			}
		case user.FieldOptNum:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field opt_num", values[i])
			} else if value.Valid {
				u.OptNum = int(value.Int64)
			}
		case user.FieldOptStr:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field opt_str", values[i])
			} else if value.Valid {
				u.OptStr = value.String
			}
		case user.FieldOptBool:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field opt_bool", values[i])
			} else if value.Valid {
				u.OptBool = value.Bool
			}
		case user.FieldBigInt:
			if value, ok := values[i].(*schema.BigInt); !ok {
				return fmt.Errorf("unexpected type %T for field big_int", values[i])
			} else if value != nil {
				u.BigInt = *value
			}
		case user.FieldBUser1:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field b_user_1", values[i])
			} else if value.Valid {
				u.BUser1 = int(value.Int64)
			}
		case user.FieldHeightInCm:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field height_in_cm", values[i])
			} else if value.Valid {
				u.HeightInCm = float32(value.Float64)
			}
		case user.FieldAccountBalance:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field account_balance", values[i])
			} else if value.Valid {
				u.AccountBalance = value.Float64
			}
		case user.FieldUnnecessary:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field unnecessary", values[i])
			} else if value.Valid {
				u.Unnecessary = value.String
			}
		case user.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				u.Type = value.String
			}
		case user.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &u.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case user.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field user_group", value)
			} else if value.Valid {
				u.user_group = new(int)
				*u.user_group = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryGroup queries the "group" edge of the User entity.
func (u *User) QueryGroup() *GroupQuery {
	return (&UserClient{config: u.config}).QueryGroup(u)
}

// QueryAttachment queries the "attachment" edge of the User entity.
func (u *User) QueryAttachment() *AttachmentQuery {
	return (&UserClient{config: u.config}).QueryAttachment(u)
}

// QueryReceived1 queries the "received_1" edge of the User entity.
func (u *User) QueryReceived1() *AttachmentQuery {
	return (&UserClient{config: u.config}).QueryReceived1(u)
}

// QueryPet queries the "pet" edge of the User entity.
func (u *User) QueryPet() *PetQuery {
	return (&UserClient{config: u.config}).QueryPet(u)
}

// QuerySkipEdge queries the "skip_edge" edge of the User entity.
func (u *User) QuerySkipEdge() *SkipEdgeExampleQuery {
	return (&UserClient{config: u.config}).QuerySkipEdge(u)
}

// Update returns a builder for updating this User.
// Note that you need to call User.Unwrap() before calling this method if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return (&UserClient{config: u.config}).UpdateOne(u)
}

// Unwrap unwraps the User entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v, ", u.ID))
	builder.WriteString("user_name=")
	builder.WriteString(u.UserName)
	builder.WriteString(", ")
	builder.WriteString("joined=")
	builder.WriteString(u.Joined.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("points=")
	builder.WriteString(fmt.Sprintf("%v", u.Points))
	builder.WriteString(", ")
	builder.WriteString("exp=")
	builder.WriteString(fmt.Sprintf("%v", u.Exp))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", u.Status))
	builder.WriteString(", ")
	builder.WriteString("external_id=")
	builder.WriteString(fmt.Sprintf("%v", u.ExternalID))
	builder.WriteString(", ")
	builder.WriteString("crm_id=")
	builder.WriteString(fmt.Sprintf("%v", u.CrmID))
	builder.WriteString(", ")
	builder.WriteString("banned=")
	builder.WriteString(fmt.Sprintf("%v", u.Banned))
	builder.WriteString(", ")
	builder.WriteString("custom_pb=")
	builder.WriteString(fmt.Sprintf("%v", u.CustomPb))
	builder.WriteString(", ")
	builder.WriteString("opt_num=")
	builder.WriteString(fmt.Sprintf("%v", u.OptNum))
	builder.WriteString(", ")
	builder.WriteString("opt_str=")
	builder.WriteString(u.OptStr)
	builder.WriteString(", ")
	builder.WriteString("opt_bool=")
	builder.WriteString(fmt.Sprintf("%v", u.OptBool))
	builder.WriteString(", ")
	builder.WriteString("big_int=")
	builder.WriteString(fmt.Sprintf("%v", u.BigInt))
	builder.WriteString(", ")
	builder.WriteString("b_user_1=")
	builder.WriteString(fmt.Sprintf("%v", u.BUser1))
	builder.WriteString(", ")
	builder.WriteString("height_in_cm=")
	builder.WriteString(fmt.Sprintf("%v", u.HeightInCm))
	builder.WriteString(", ")
	builder.WriteString("account_balance=")
	builder.WriteString(fmt.Sprintf("%v", u.AccountBalance))
	builder.WriteString(", ")
	builder.WriteString("unnecessary=")
	builder.WriteString(u.Unnecessary)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(u.Type)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", u.Labels))
	builder.WriteByte(')')
	return builder.String()
}

// Users is a parsable slice of User.
type Users []*User

func (u Users) config(cfg config) {
	for _i := range u {
		u[_i].config = cfg
	}
}
