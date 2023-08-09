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
	FieldID = "id"
	// FieldUserName holds the string denoting the user_name field in the database.
	FieldUserName = "user_name"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldUnnecessary holds the string denoting the unnecessary field in the database.
	FieldUnnecessary = "unnecessary"
	// EdgeBlogPosts holds the string denoting the blog_posts edge name in mutations.
	EdgeBlogPosts = "blog_posts"
	// EdgeProfilePic holds the string denoting the profile_pic edge name in mutations.
	EdgeProfilePic = "profile_pic"
	// EdgeSkipEdge holds the string denoting the skip_edge edge name in mutations.
	EdgeSkipEdge = "skip_edge"
	// Table holds the table name of the user in the database.
	Table = "users"
	// BlogPostsTable is the table that holds the blog_posts relation/edge.
	BlogPostsTable = "blog_posts"
	// BlogPostsInverseTable is the table name for the BlogPost entity.
	// It exists in this package in order to avoid circular dependency with the "blogpost" package.
	BlogPostsInverseTable = "blog_posts"
	// BlogPostsColumn is the table column denoting the blog_posts relation/edge.
	BlogPostsColumn = "blog_post_author"
	// ProfilePicTable is the table that holds the profile_pic relation/edge.
	ProfilePicTable = "users"
	// ProfilePicInverseTable is the table name for the Image entity.
	// It exists in this package in order to avoid circular dependency with the "image" package.
	ProfilePicInverseTable = "images"
	// ProfilePicColumn is the table column denoting the profile_pic relation/edge.
	ProfilePicColumn = "user_profile_pic"
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
	FieldStatus,
	FieldUnnecessary,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "users"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_profile_pic",
}

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

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByUnnecessary orders the results by the unnecessary field.
func ByUnnecessary(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUnnecessary, opts...).ToFunc()
}

// ByBlogPostsCount orders the results by blog_posts count.
func ByBlogPostsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBlogPostsStep(), opts...)
	}
}

// ByBlogPosts orders the results by blog_posts terms.
func ByBlogPosts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBlogPostsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByProfilePicField orders the results by profile_pic field.
func ByProfilePicField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProfilePicStep(), sql.OrderByField(field, opts...))
	}
}

// BySkipEdgeField orders the results by skip_edge field.
func BySkipEdgeField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSkipEdgeStep(), sql.OrderByField(field, opts...))
	}
}
func newBlogPostsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BlogPostsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, BlogPostsTable, BlogPostsColumn),
	)
}
func newProfilePicStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProfilePicInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ProfilePicTable, ProfilePicColumn),
	)
}
func newSkipEdgeStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SkipEdgeInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, SkipEdgeTable, SkipEdgeColumn),
	)
}
