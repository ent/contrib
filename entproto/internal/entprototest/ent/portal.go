// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/contrib/entproto/internal/entprototest/ent/category"
	"entgo.io/contrib/entproto/internal/entprototest/ent/portal"
	"entgo.io/ent/dialect/sql"
)

// Portal is the model entity for the Portal schema.
type Portal struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PortalQuery when eager-loading is set.
	Edges           PortalEdges `json:"edges"`
	portal_category *int
}

// PortalEdges holds the relations/edges for other nodes in the graph.
type PortalEdges struct {
	// Category holds the value of the category edge.
	Category *Category `json:"category,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// CategoryOrErr returns the Category value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PortalEdges) CategoryOrErr() (*Category, error) {
	if e.loadedTypes[0] {
		if e.Category == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: category.Label}
		}
		return e.Category, nil
	}
	return nil, &NotLoadedError{edge: "category"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Portal) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case portal.FieldID:
			values[i] = new(sql.NullInt64)
		case portal.FieldName, portal.FieldDescription:
			values[i] = new(sql.NullString)
		case portal.ForeignKeys[0]: // portal_category
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Portal", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Portal fields.
func (po *Portal) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case portal.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			po.ID = int(value.Int64)
		case portal.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				po.Name = value.String
			}
		case portal.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				po.Description = value.String
			}
		case portal.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field portal_category", value)
			} else if value.Valid {
				po.portal_category = new(int)
				*po.portal_category = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryCategory queries the "category" edge of the Portal entity.
func (po *Portal) QueryCategory() *CategoryQuery {
	return (&PortalClient{config: po.config}).QueryCategory(po)
}

// Update returns a builder for updating this Portal.
// Note that you need to call Portal.Unwrap() before calling this method if this Portal
// was returned from a transaction, and the transaction was committed or rolled back.
func (po *Portal) Update() *PortalUpdateOne {
	return (&PortalClient{config: po.config}).UpdateOne(po)
}

// Unwrap unwraps the Portal entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (po *Portal) Unwrap() *Portal {
	_tx, ok := po.config.driver.(*txDriver)
	if !ok {
		panic("ent: Portal is not a transactional entity")
	}
	po.config.driver = _tx.drv
	return po
}

// String implements the fmt.Stringer.
func (po *Portal) String() string {
	var builder strings.Builder
	builder.WriteString("Portal(")
	builder.WriteString(fmt.Sprintf("id=%v, ", po.ID))
	builder.WriteString("name=")
	builder.WriteString(po.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(po.Description)
	builder.WriteByte(')')
	return builder.String()
}

// Portals is a parsable slice of Portal.
type Portals []*Portal

func (po Portals) config(cfg config) {
	for _i := range po {
		po[_i].config = cfg
	}
}
