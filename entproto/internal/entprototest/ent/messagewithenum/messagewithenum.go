// Code generated by ent, DO NOT EDIT.

package messagewithenum

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the messagewithenum type in the database.
	Label = "message_with_enum"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldEnumType holds the string denoting the enum_type field in the database.
	FieldEnumType = "enum_type"
	// FieldEnumWithoutDefault holds the string denoting the enum_without_default field in the database.
	FieldEnumWithoutDefault = "enum_without_default"
	// Table holds the table name of the messagewithenum in the database.
	Table = "message_with_enums"
)

// Columns holds all SQL columns for messagewithenum fields.
var Columns = []string{
	FieldID,
	FieldEnumType,
	FieldEnumWithoutDefault,
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

// EnumType defines the type for the "enum_type" enum field.
type EnumType string

// EnumTypePending is the default value of the EnumType enum.
const DefaultEnumType = EnumTypePending

// EnumType values.
const (
	EnumTypePending   EnumType = "pending"
	EnumTypeActive    EnumType = "active"
	EnumTypeSuspended EnumType = "suspended"
	EnumTypeDeleted   EnumType = "deleted"
)

func (et EnumType) String() string {
	return string(et)
}

// EnumTypeValidator is a validator for the "enum_type" field enum values. It is called by the builders before save.
func EnumTypeValidator(et EnumType) error {
	switch et {
	case EnumTypePending, EnumTypeActive, EnumTypeSuspended, EnumTypeDeleted:
		return nil
	default:
		return fmt.Errorf("messagewithenum: invalid enum value for enum_type field: %q", et)
	}
}

// EnumWithoutDefault defines the type for the "enum_without_default" enum field.
type EnumWithoutDefault string

// EnumWithoutDefault values.
const (
	EnumWithoutDefaultFirst  EnumWithoutDefault = "first"
	EnumWithoutDefaultSecond EnumWithoutDefault = "second"
)

func (ewd EnumWithoutDefault) String() string {
	return string(ewd)
}

// EnumWithoutDefaultValidator is a validator for the "enum_without_default" field enum values. It is called by the builders before save.
func EnumWithoutDefaultValidator(ewd EnumWithoutDefault) error {
	switch ewd {
	case EnumWithoutDefaultFirst, EnumWithoutDefaultSecond:
		return nil
	default:
		return fmt.Errorf("messagewithenum: invalid enum value for enum_without_default field: %q", ewd)
	}
}

// OrderOption defines the ordering options for the MessageWithEnum queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByEnumType orders the results by the enum_type field.
func ByEnumType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnumType, opts...).ToFunc()
}

// ByEnumWithoutDefault orders the results by the enum_without_default field.
func ByEnumWithoutDefault(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnumWithoutDefault, opts...).ToFunc()
}
