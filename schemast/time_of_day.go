package schemast

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log/slog"
)

type TimeOfDay struct {
	Hour   int
	Minute int
	Second int
	driver.Valuer
	sql.Scanner
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *TimeOfDay) UnmarshalGQL(v interface{}) error {
	time, ok := v.(string)
	if !ok {
		return fmt.Errorf("value provided as TimeOfDay is not a string")
	}

	_, err := fmt.Sscanf(time, "%d:%d:%d", &t.Hour, &t.Minute, &t.Second)
	if err != nil {
		return fmt.Errorf("invalid format, could not parse TimeOfDay: %w", err)
	}

	if t.Hour < 0 || t.Hour > 23 {
		return fmt.Errorf("invalid format, hour must be between 0 and 23")
	}
	if t.Minute < 0 || t.Minute > 59 {
		return fmt.Errorf("invalid format, minute must be between 0 and 59")
	}
	if t.Second < 0 || t.Second > 59 {
		return fmt.Errorf("invalid format, second must be between 0 and 59")
	}

	return err
}

// MarshalGQL implements the graphql.Marshaler interface
func (t TimeOfDay) MarshalGQL(w io.Writer) {
	_, err := w.Write([]byte(`"` + t.String() + `"`))
	if err != nil {
		slog.Error("Could not Marshal TimeOfDay")
	}
}

func (t *TimeOfDay) String() string {
	return fmt.Sprintf("%d:%d:%d", t.Hour, t.Minute, t.Second)
}

func (t TimeOfDay) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *TimeOfDay) Scan(src interface{}) error {
	bytes, ok := src.([]uint8)
	if !ok {
		return fmt.Errorf("could not cast database TimeOfDay value to []uint8")
	}

	_, err := fmt.Sscanf(string(bytes), "%d:%d:%d", &t.Hour, &t.Minute, &t.Second)
	if err != nil {
		return fmt.Errorf("invalid format, could not parse TimeOfDay: %w", err)
	}

	return nil
}
