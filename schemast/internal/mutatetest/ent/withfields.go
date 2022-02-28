// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/withfields"
	"entgo.io/ent/dialect/sql"
)

// WithFields is the model entity for the WithFields schema.
type WithFields struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Existing holds the value of the "existing" field.
	Existing string `json:"existing,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*WithFields) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case withfields.FieldID:
			values[i] = new(sql.NullInt64)
		case withfields.FieldExisting:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type WithFields", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the WithFields fields.
func (wf *WithFields) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case withfields.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			wf.ID = int(value.Int64)
		case withfields.FieldExisting:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field existing", values[i])
			} else if value.Valid {
				wf.Existing = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this WithFields.
// Note that you need to call WithFields.Unwrap() before calling this method if this WithFields
// was returned from a transaction, and the transaction was committed or rolled back.
func (wf *WithFields) Update() *WithFieldsUpdateOne {
	return (&WithFieldsClient{config: wf.config}).UpdateOne(wf)
}

// Unwrap unwraps the WithFields entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (wf *WithFields) Unwrap() *WithFields {
	tx, ok := wf.config.driver.(*txDriver)
	if !ok {
		panic("ent: WithFields is not a transactional entity")
	}
	wf.config.driver = tx.drv
	return wf
}

// String implements the fmt.Stringer.
func (wf *WithFields) String() string {
	var builder strings.Builder
	builder.WriteString("WithFields(")
	builder.WriteString(fmt.Sprintf("id=%v", wf.ID))
	builder.WriteString(", existing=")
	builder.WriteString(wf.Existing)
	builder.WriteByte(')')
	return builder.String()
}

// WithFieldsSlice is a parsable slice of WithFields.
type WithFieldsSlice []*WithFields

func (wf WithFieldsSlice) config(cfg config) {
	for _i := range wf {
		wf[_i].config = cfg
	}
}
