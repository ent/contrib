// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/contrib/entproto/internal/entprototest/ent/messagewithints"
	"entgo.io/ent/dialect/sql"
)

// MessageWithInts is the model entity for the MessageWithInts schema.
type MessageWithInts struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Ints holds the value of the "ints" field.
	Ints []int `json:"ints,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*MessageWithInts) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case messagewithints.FieldInts:
			values[i] = new([]byte)
		case messagewithints.FieldID:
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type MessageWithInts", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the MessageWithInts fields.
func (mwi *MessageWithInts) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case messagewithints.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			mwi.ID = int(value.Int64)
		case messagewithints.FieldInts:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field ints", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &mwi.Ints); err != nil {
					return fmt.Errorf("unmarshal field ints: %w", err)
				}
			}
		}
	}
	return nil
}

// Update returns a builder for updating this MessageWithInts.
// Note that you need to call MessageWithInts.Unwrap() before calling this method if this MessageWithInts
// was returned from a transaction, and the transaction was committed or rolled back.
func (mwi *MessageWithInts) Update() *MessageWithIntsUpdateOne {
	return (&MessageWithIntsClient{config: mwi.config}).UpdateOne(mwi)
}

// Unwrap unwraps the MessageWithInts entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (mwi *MessageWithInts) Unwrap() *MessageWithInts {
	_tx, ok := mwi.config.driver.(*txDriver)
	if !ok {
		panic("ent: MessageWithInts is not a transactional entity")
	}
	mwi.config.driver = _tx.drv
	return mwi
}

// String implements the fmt.Stringer.
func (mwi *MessageWithInts) String() string {
	var builder strings.Builder
	builder.WriteString("MessageWithInts(")
	builder.WriteString(fmt.Sprintf("id=%v, ", mwi.ID))
	builder.WriteString("ints=")
	builder.WriteString(fmt.Sprintf("%v", mwi.Ints))
	builder.WriteByte(')')
	return builder.String()
}

// MessageWithIntsSlice is a parsable slice of MessageWithInts.
type MessageWithIntsSlice []*MessageWithInts

func (mwi MessageWithIntsSlice) config(cfg config) {
	for _i := range mwi {
		mwi[_i].config = cfg
	}
}
