// Code generated by ent, DO NOT EDIT.

package simple

import (
	"fmt"
	"strings"

	"entgo.io/contrib/entoas/internal/simple/category"
	"entgo.io/ent/dialect/sql"
)

// Category is the model entity for the Category schema.
type Category struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Readonly holds the value of the "readonly" field.
	Readonly string `json:"readonly,omitempty"`
	// SkipInSpec holds the value of the "skip_in_spec" field.
	SkipInSpec string `json:"skip_in_spec,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CategoryQuery when eager-loading is set.
	Edges CategoryEdges `json:"edges"`
}

// CategoryEdges holds the relations/edges for other nodes in the graph.
type CategoryEdges struct {
	// Pets holds the value of the pets edge.
	Pets []*Pet `json:"pets,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// PetsOrErr returns the Pets value or an error if the edge
// was not loaded in eager-loading.
func (e CategoryEdges) PetsOrErr() ([]*Pet, error) {
	if e.loadedTypes[0] {
		return e.Pets, nil
	}
	return nil, &NotLoadedError{edge: "pets"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Category) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case category.FieldID:
			values[i] = new(sql.NullInt64)
		case category.FieldName, category.FieldReadonly, category.FieldSkipInSpec:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Category", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Category fields.
func (c *Category) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case category.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case category.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case category.FieldReadonly:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field readonly", values[i])
			} else if value.Valid {
				c.Readonly = value.String
			}
		case category.FieldSkipInSpec:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field skip_in_spec", values[i])
			} else if value.Valid {
				c.SkipInSpec = value.String
			}
		}
	}
	return nil
}

// QueryPets queries the "pets" edge of the Category entity.
func (c *Category) QueryPets() *PetQuery {
	return (&CategoryClient{config: c.config}).QueryPets(c)
}

// Update returns a builder for updating this Category.
// Note that you need to call Category.Unwrap() before calling this method if this Category
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Category) Update() *CategoryUpdateOne {
	return (&CategoryClient{config: c.config}).UpdateOne(c)
}

// Unwrap unwraps the Category entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Category) Unwrap() *Category {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("simple: Category is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Category) String() string {
	var builder strings.Builder
	builder.WriteString("Category(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("readonly=")
	builder.WriteString(c.Readonly)
	builder.WriteString(", ")
	builder.WriteString("skip_in_spec=")
	builder.WriteString(c.SkipInSpec)
	builder.WriteByte(')')
	return builder.String()
}

// Categories is a parsable slice of Category.
type Categories []*Category

func (c Categories) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
