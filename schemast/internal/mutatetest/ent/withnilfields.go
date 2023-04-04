// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/contrib/schemast/internal/mutatetest/ent/withnilfields"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// WithNilFields is the model entity for the WithNilFields schema.
type WithNilFields struct {
	config
	// ID of the ent.
	ID           int `json:"id,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*WithNilFields) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case withnilfields.FieldID:
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the WithNilFields fields.
func (wnf *WithNilFields) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case withnilfields.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			wnf.ID = int(value.Int64)
		default:
			wnf.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the WithNilFields.
// This includes values selected through modifiers, order, etc.
func (wnf *WithNilFields) Value(name string) (ent.Value, error) {
	return wnf.selectValues.Get(name)
}

// Update returns a builder for updating this WithNilFields.
// Note that you need to call WithNilFields.Unwrap() before calling this method if this WithNilFields
// was returned from a transaction, and the transaction was committed or rolled back.
func (wnf *WithNilFields) Update() *WithNilFieldsUpdateOne {
	return NewWithNilFieldsClient(wnf.config).UpdateOne(wnf)
}

// Unwrap unwraps the WithNilFields entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (wnf *WithNilFields) Unwrap() *WithNilFields {
	_tx, ok := wnf.config.driver.(*txDriver)
	if !ok {
		panic("ent: WithNilFields is not a transactional entity")
	}
	wnf.config.driver = _tx.drv
	return wnf
}

// String implements the fmt.Stringer.
func (wnf *WithNilFields) String() string {
	var builder strings.Builder
	builder.WriteString("WithNilFields(")
	builder.WriteString(fmt.Sprintf("id=%v", wnf.ID))
	builder.WriteByte(')')
	return builder.String()
}

// WithNilFieldsSlice is a parsable slice of WithNilFields.
type WithNilFieldsSlice []*WithNilFields
