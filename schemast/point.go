package schemast

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

type Point struct {
	Lat float64 `graphql:"lat"`
	Lon float64 `graphql:"lon"`
}

// Scan implements the Scanner interface.
func (p *Point) Scan(value any) error {
	bin, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid binary value for point")
	}
	var op orb.Point
	if err := wkb.Scanner(&op).Scan(bin[4:]); err != nil {
		return err
	}
	p.Lon, p.Lat = op.X(), op.Y()
	return nil
}

// Value implements the driver Valuer interface.
func (p Point) Value() (driver.Value, error) {
	op := orb.Point{p.Lon, p.Lat}
	return wkb.Value(op).Value()
}

// FormatParam implements the sql.ParamFormatter interface to tell the SQL
// builder that the placeholder for a Point parameter needs to be formatted.
func (p Point) FormatParam(placeholder string, info *sql.StmtInfo) string {
	if info.Dialect == dialect.MySQL {
		return "ST_GeomFromWKB(" + placeholder + ")"
	}
	return placeholder
}

// SchemaType defines the schema-type of the Point object.
func (Point) SchemaType() map[string]string {
	return map[string]string{
		dialect.MySQL: "POINT",
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (p *Point) UnmarshalGQL(v any) error {
	pt, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid format")
	}

	*p = Point{
		Lat: pt["lat"].(float64),
		Lon: pt["lon"].(float64),
	}

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (p Point) MarshalGQL(w io.Writer) {
	b, err := json.Marshal(&p)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
