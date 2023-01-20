// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/contrib/entproto/internal/entprototest/ent/image"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// Image is the model entity for the Image schema.
type Image struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// URLPath holds the value of the "url_path" field.
	URLPath string `json:"url_path,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ImageQuery when eager-loading is set.
	Edges             ImageEdges `json:"edges"`
	no_backref_images *int
}

// ImageEdges holds the relations/edges for other nodes in the graph.
type ImageEdges struct {
	// UserProfilePic holds the value of the user_profile_pic edge.
	UserProfilePic []*User `json:"user_profile_pic,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// UserProfilePicOrErr returns the UserProfilePic value or an error if the edge
// was not loaded in eager-loading.
func (e ImageEdges) UserProfilePicOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.UserProfilePic, nil
	}
	return nil, &NotLoadedError{edge: "user_profile_pic"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Image) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case image.FieldURLPath:
			values[i] = new(sql.NullString)
		case image.FieldID:
			values[i] = new(uuid.UUID)
		case image.ForeignKeys[0]: // no_backref_images
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Image", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Image fields.
func (i *Image) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case image.FieldID:
			if value, ok := values[j].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[j])
			} else if value != nil {
				i.ID = *value
			}
		case image.FieldURLPath:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field url_path", values[j])
			} else if value.Valid {
				i.URLPath = value.String
			}
		case image.ForeignKeys[0]:
			if value, ok := values[j].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field no_backref_images", value)
			} else if value.Valid {
				i.no_backref_images = new(int)
				*i.no_backref_images = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryUserProfilePic queries the "user_profile_pic" edge of the Image entity.
func (i *Image) QueryUserProfilePic() *UserQuery {
	return NewImageClient(i.config).QueryUserProfilePic(i)
}

// Update returns a builder for updating this Image.
// Note that you need to call Image.Unwrap() before calling this method if this Image
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Image) Update() *ImageUpdateOne {
	return NewImageClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Image entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Image) Unwrap() *Image {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Image is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Image) String() string {
	var builder strings.Builder
	builder.WriteString("Image(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("url_path=")
	builder.WriteString(i.URLPath)
	builder.WriteByte(')')
	return builder.String()
}

// Images is a parsable slice of Image.
type Images []*Image

func (i Images) config(cfg config) {
	for _i := range i {
		i[_i].config = cfg
	}
}
