// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package todo

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/contrib/entgql/internal/todouuid/ent"
	"github.com/google/uuid"
)

type NamedNode interface {
	IsNamedNode()
}

type CategoryTypes struct {
	Public *bool `json:"public,omitempty"`
}

type CategoryTypesInput struct {
	Public *bool `json:"public,omitempty"`
}

type OneToMany struct {
	ID       uuid.UUID    `json:"id"`
	Name     string       `json:"name"`
	Field2   *string      `json:"field2,omitempty"`
	Parent   *OneToMany   `json:"parent,omitempty"`
	Children []*OneToMany `json:"children,omitempty"`
}

func (OneToMany) IsNode() {}

// A connection to a list of items.
type OneToManyConnection struct {
	// A list of edges.
	Edges []*OneToManyEdge `json:"edges,omitempty"`
	// Information to aid in pagination.
	PageInfo *entgql.PageInfo[uuid.UUID] `json:"pageInfo"`
	// Identifies the total count of items in the connection.
	TotalCount int `json:"totalCount"`
}

// An edge in a connection.
type OneToManyEdge struct {
	// The item at the end of the edge.
	Node *OneToMany `json:"node"`
	// A cursor for use in pagination.
	Cursor entgql.Cursor[uuid.UUID] `json:"cursor"`
}

// Ordering options for OneToMany connections
type OneToManyOrder struct {
	// The ordering direction.
	Direction entgql.OrderDirection `json:"direction"`
	// The field by which to order OneToManies.
	Field OneToManyOrderField `json:"field"`
}

// OneToManyWhereInput is used for filtering OneToMany objects.
// Input was generated by ent.
type OneToManyWhereInput struct {
	Not *OneToManyWhereInput   `json:"not,omitempty"`
	And []*OneToManyWhereInput `json:"and,omitempty"`
	Or  []*OneToManyWhereInput `json:"or,omitempty"`
	// id field predicates
	ID      *uuid.UUID  `json:"id,omitempty"`
	IDNeq   *uuid.UUID  `json:"idNEQ,omitempty"`
	IDIn    []uuid.UUID `json:"idIn,omitempty"`
	IDNotIn []uuid.UUID `json:"idNotIn,omitempty"`
	IDGt    *uuid.UUID  `json:"idGT,omitempty"`
	IDGte   *uuid.UUID  `json:"idGTE,omitempty"`
	IDLt    *uuid.UUID  `json:"idLT,omitempty"`
	IDLte   *uuid.UUID  `json:"idLTE,omitempty"`
	// name field predicates
	Name             *string  `json:"name,omitempty"`
	NameNeq          *string  `json:"nameNEQ,omitempty"`
	NameIn           []string `json:"nameIn,omitempty"`
	NameNotIn        []string `json:"nameNotIn,omitempty"`
	NameGt           *string  `json:"nameGT,omitempty"`
	NameGte          *string  `json:"nameGTE,omitempty"`
	NameLt           *string  `json:"nameLT,omitempty"`
	NameLte          *string  `json:"nameLTE,omitempty"`
	NameContains     *string  `json:"nameContains,omitempty"`
	NameHasPrefix    *string  `json:"nameHasPrefix,omitempty"`
	NameHasSuffix    *string  `json:"nameHasSuffix,omitempty"`
	NameEqualFold    *string  `json:"nameEqualFold,omitempty"`
	NameContainsFold *string  `json:"nameContainsFold,omitempty"`
	// field2 field predicates
	Field2             *string  `json:"field2,omitempty"`
	Field2neq          *string  `json:"field2NEQ,omitempty"`
	Field2In           []string `json:"field2In,omitempty"`
	Field2NotIn        []string `json:"field2NotIn,omitempty"`
	Field2gt           *string  `json:"field2GT,omitempty"`
	Field2gte          *string  `json:"field2GTE,omitempty"`
	Field2lt           *string  `json:"field2LT,omitempty"`
	Field2lte          *string  `json:"field2LTE,omitempty"`
	Field2Contains     *string  `json:"field2Contains,omitempty"`
	Field2HasPrefix    *string  `json:"field2HasPrefix,omitempty"`
	Field2HasSuffix    *string  `json:"field2HasSuffix,omitempty"`
	Field2IsNil        *bool    `json:"field2IsNil,omitempty"`
	Field2NotNil       *bool    `json:"field2NotNil,omitempty"`
	Field2EqualFold    *string  `json:"field2EqualFold,omitempty"`
	Field2ContainsFold *string  `json:"field2ContainsFold,omitempty"`
	// parent edge predicates
	HasParent     *bool                  `json:"hasParent,omitempty"`
	HasParentWith []*OneToManyWhereInput `json:"hasParentWith,omitempty"`
	// children edge predicates
	HasChildren     *bool                  `json:"hasChildren,omitempty"`
	HasChildrenWith []*OneToManyWhereInput `json:"hasChildrenWith,omitempty"`
}

// OrganizationWhereInput is used for filtering Workspace objects.
// Input was generated by ent.
type OrganizationWhereInput struct {
	Not *OrganizationWhereInput   `json:"not,omitempty"`
	And []*OrganizationWhereInput `json:"and,omitempty"`
	Or  []*OrganizationWhereInput `json:"or,omitempty"`
	// id field predicates
	ID      *uuid.UUID  `json:"id,omitempty"`
	IDNeq   *uuid.UUID  `json:"idNEQ,omitempty"`
	IDIn    []uuid.UUID `json:"idIn,omitempty"`
	IDNotIn []uuid.UUID `json:"idNotIn,omitempty"`
	IDGt    *uuid.UUID  `json:"idGT,omitempty"`
	IDGte   *uuid.UUID  `json:"idGTE,omitempty"`
	IDLt    *uuid.UUID  `json:"idLT,omitempty"`
	IDLte   *uuid.UUID  `json:"idLTE,omitempty"`
	// name field predicates
	Name             *string  `json:"name,omitempty"`
	NameNeq          *string  `json:"nameNEQ,omitempty"`
	NameIn           []string `json:"nameIn,omitempty"`
	NameNotIn        []string `json:"nameNotIn,omitempty"`
	NameGt           *string  `json:"nameGT,omitempty"`
	NameGte          *string  `json:"nameGTE,omitempty"`
	NameLt           *string  `json:"nameLT,omitempty"`
	NameLte          *string  `json:"nameLTE,omitempty"`
	NameContains     *string  `json:"nameContains,omitempty"`
	NameHasPrefix    *string  `json:"nameHasPrefix,omitempty"`
	NameHasSuffix    *string  `json:"nameHasSuffix,omitempty"`
	NameEqualFold    *string  `json:"nameEqualFold,omitempty"`
	NameContainsFold *string  `json:"nameContainsFold,omitempty"`
}

type Project struct {
	ID    uuid.UUID           `json:"id"`
	Todos *ent.TodoConnection `json:"todos"`
}

func (Project) IsNode() {}

// ProjectWhereInput is used for filtering Project objects.
// Input was generated by ent.
type ProjectWhereInput struct {
	Not *ProjectWhereInput   `json:"not,omitempty"`
	And []*ProjectWhereInput `json:"and,omitempty"`
	Or  []*ProjectWhereInput `json:"or,omitempty"`
	// id field predicates
	ID      *uuid.UUID  `json:"id,omitempty"`
	IDNeq   *uuid.UUID  `json:"idNEQ,omitempty"`
	IDIn    []uuid.UUID `json:"idIn,omitempty"`
	IDNotIn []uuid.UUID `json:"idNotIn,omitempty"`
	IDGt    *uuid.UUID  `json:"idGT,omitempty"`
	IDGte   *uuid.UUID  `json:"idGTE,omitempty"`
	IDLt    *uuid.UUID  `json:"idLT,omitempty"`
	IDLte   *uuid.UUID  `json:"idLTE,omitempty"`
	// todos edge predicates
	HasTodos     *bool                 `json:"hasTodos,omitempty"`
	HasTodosWith []*ent.TodoWhereInput `json:"hasTodosWith,omitempty"`
}

// UpdateFriendshipInput is used for update Friendship object.
// Input was generated by ent.
type UpdateFriendshipInput struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UserID    *uuid.UUID `json:"userID,omitempty"`
	FriendID  *uuid.UUID `json:"friendID,omitempty"`
}

// Properties by which OneToMany connections can be ordered.
type OneToManyOrderField string

const (
	OneToManyOrderFieldName OneToManyOrderField = "NAME"
)

var AllOneToManyOrderField = []OneToManyOrderField{
	OneToManyOrderFieldName,
}

func (e OneToManyOrderField) IsValid() bool {
	switch e {
	case OneToManyOrderFieldName:
		return true
	}
	return false
}

func (e OneToManyOrderField) String() string {
	return string(e)
}

func (e *OneToManyOrderField) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OneToManyOrderField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OneToManyOrderField", str)
	}
	return nil
}

func (e OneToManyOrderField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
