// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/bionicstork/bionicstork/pkg/entproto/internal/todo/ent/nilexample"
	"entgo.io/ent/dialect/sql"
)

// NilExample is the model entity for the NilExample schema.
type NilExample struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// StrNil holds the value of the "str_nil" field.
	StrNil *string `json:"str_nil,omitempty"`
	// TimeNil holds the value of the "time_nil" field.
	TimeNil *time.Time `json:"time_nil,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*NilExample) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case nilexample.FieldID:
			values[i] = new(sql.NullInt64)
		case nilexample.FieldStrNil:
			values[i] = new(sql.NullString)
		case nilexample.FieldTimeNil:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type NilExample", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the NilExample fields.
func (ne *NilExample) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case nilexample.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ne.ID = int(value.Int64)
		case nilexample.FieldStrNil:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field str_nil", values[i])
			} else if value.Valid {
				ne.StrNil = new(string)
				*ne.StrNil = value.String
			}
		case nilexample.FieldTimeNil:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field time_nil", values[i])
			} else if value.Valid {
				ne.TimeNil = new(time.Time)
				*ne.TimeNil = value.Time
			}
		}
	}
	return nil
}

// Update returns a builder for updating this NilExample.
// Note that you need to call NilExample.Unwrap() before calling this method if this NilExample
// was returned from a transaction, and the transaction was committed or rolled back.
func (ne *NilExample) Update() *NilExampleUpdateOne {
	return (&NilExampleClient{config: ne.config}).UpdateOne(ne)
}

// Unwrap unwraps the NilExample entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ne *NilExample) Unwrap() *NilExample {
	tx, ok := ne.config.driver.(*txDriver)
	if !ok {
		panic("ent: NilExample is not a transactional entity")
	}
	ne.config.driver = tx.drv
	return ne
}

// String implements the fmt.Stringer.
func (ne *NilExample) String() string {
	var builder strings.Builder
	builder.WriteString("NilExample(")
	builder.WriteString(fmt.Sprintf("id=%v", ne.ID))
	if v := ne.StrNil; v != nil {
		builder.WriteString(", str_nil=")
		builder.WriteString(*v)
	}
	if v := ne.TimeNil; v != nil {
		builder.WriteString(", time_nil=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// NilExamples is a parsable slice of NilExample.
type NilExamples []*NilExample

func (ne NilExamples) config(cfg config) {
	for _i := range ne {
		ne[_i].config = cfg
	}
}
