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
//
// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/contrib/entgql/internal/todouuid/ent/user"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Username holds the value of the "username" field.
	Username uuid.UUID `json:"username,omitempty"`
	// Password holds the value of the "password" field.
	Password string `json:"-"`
	// Metadata holds the value of the "metadata" field.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserQuery when eager-loading is set.
	Edges UserEdges `json:"edges"`
}

// UserEdges holds the relations/edges for other nodes in the graph.
type UserEdges struct {
	// The groups of the user
	Groups []*Group `json:"groups,omitempty"`
	// Friends holds the value of the friends edge.
	Friends []*User `json:"friends,omitempty"`
	// Friendships holds the value of the friendships edge.
	Friendships []*Friendship `json:"friendships,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
	// totalCount holds the count of the edges above.
	totalCount [3]map[string]int

	namedGroups      map[string][]*Group
	namedFriends     map[string][]*User
	namedFriendships map[string][]*Friendship
}

// GroupsOrErr returns the Groups value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) GroupsOrErr() ([]*Group, error) {
	if e.loadedTypes[0] {
		return e.Groups, nil
	}
	return nil, &NotLoadedError{edge: "groups"}
}

// FriendsOrErr returns the Friends value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) FriendsOrErr() ([]*User, error) {
	if e.loadedTypes[1] {
		return e.Friends, nil
	}
	return nil, &NotLoadedError{edge: "friends"}
}

// FriendshipsOrErr returns the Friendships value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) FriendshipsOrErr() ([]*Friendship, error) {
	if e.loadedTypes[2] {
		return e.Friendships, nil
	}
	return nil, &NotLoadedError{edge: "friendships"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case user.FieldMetadata:
			values[i] = new([]byte)
		case user.FieldName, user.FieldPassword:
			values[i] = new(sql.NullString)
		case user.FieldID, user.FieldUsername:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type User", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case user.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				u.ID = *value
			}
		case user.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				u.Name = value.String
			}
		case user.FieldUsername:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field username", values[i])
			} else if value != nil {
				u.Username = *value
			}
		case user.FieldPassword:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field password", values[i])
			} else if value.Valid {
				u.Password = value.String
			}
		case user.FieldMetadata:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field metadata", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &u.Metadata); err != nil {
					return fmt.Errorf("unmarshal field metadata: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryGroups queries the "groups" edge of the User entity.
func (u *User) QueryGroups() *GroupQuery {
	return NewUserClient(u.config).QueryGroups(u)
}

// QueryFriends queries the "friends" edge of the User entity.
func (u *User) QueryFriends() *UserQuery {
	return NewUserClient(u.config).QueryFriends(u)
}

// QueryFriendships queries the "friendships" edge of the User entity.
func (u *User) QueryFriendships() *FriendshipQuery {
	return NewUserClient(u.config).QueryFriendships(u)
}

// Update returns a builder for updating this User.
// Note that you need to call User.Unwrap() before calling this method if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return NewUserClient(u.config).UpdateOne(u)
}

// Unwrap unwraps the User entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	_tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = _tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v, ", u.ID))
	builder.WriteString("name=")
	builder.WriteString(u.Name)
	builder.WriteString(", ")
	builder.WriteString("username=")
	builder.WriteString(fmt.Sprintf("%v", u.Username))
	builder.WriteString(", ")
	builder.WriteString("password=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("metadata=")
	builder.WriteString(fmt.Sprintf("%v", u.Metadata))
	builder.WriteByte(')')
	return builder.String()
}

// NamedGroups returns the Groups named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedGroups(name string) ([]*Group, error) {
	if u.Edges.namedGroups == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedGroups[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedGroups(name string, edges ...*Group) {
	if u.Edges.namedGroups == nil {
		u.Edges.namedGroups = make(map[string][]*Group)
	}
	if len(edges) == 0 {
		u.Edges.namedGroups[name] = []*Group{}
	} else {
		u.Edges.namedGroups[name] = append(u.Edges.namedGroups[name], edges...)
	}
}

// NamedFriends returns the Friends named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedFriends(name string) ([]*User, error) {
	if u.Edges.namedFriends == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedFriends[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedFriends(name string, edges ...*User) {
	if u.Edges.namedFriends == nil {
		u.Edges.namedFriends = make(map[string][]*User)
	}
	if len(edges) == 0 {
		u.Edges.namedFriends[name] = []*User{}
	} else {
		u.Edges.namedFriends[name] = append(u.Edges.namedFriends[name], edges...)
	}
}

// NamedFriendships returns the Friendships named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedFriendships(name string) ([]*Friendship, error) {
	if u.Edges.namedFriendships == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedFriendships[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedFriendships(name string, edges ...*Friendship) {
	if u.Edges.namedFriendships == nil {
		u.Edges.namedFriendships = make(map[string][]*Friendship)
	}
	if len(edges) == 0 {
		u.Edges.namedFriendships[name] = []*Friendship{}
	} else {
		u.Edges.namedFriendships[name] = append(u.Edges.namedFriendships[name], edges...)
	}
}

// Users is a parsable slice of User.
type Users []*User

func (u Users) config(cfg config) {
	for _i := range u {
		u[_i].config = cfg
	}
}
