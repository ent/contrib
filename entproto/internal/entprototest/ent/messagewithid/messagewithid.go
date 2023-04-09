// Code generated by ent, DO NOT EDIT.

package messagewithid

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the messagewithid type in the database.
	Label = "message_with_id"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// Table holds the table name of the messagewithid in the database.
	Table = "message_with_ids"
)

// Columns holds all SQL columns for messagewithid fields.
var Columns = []string{
	FieldID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Order defines the ordering method for the MessageWithID queries.
type Order func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) Order {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}
