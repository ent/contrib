// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package pulidgql implements the pulid type for gqlgen use.
// A pulid is an identifier that is a two-byte prefixed ULIDs, with the first two bytes encoding the type of the entity.
package pulidgql

import (
	"crypto/rand"
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/oklog/ulid/v2"
)

// ID implements a PULID - a prefixed ULID.
type ID string

var entropy *ulid.MonotonicEntropy

func init() {
	// TODO: Real applications would likely seed entropy in different ways.
	entropy = ulid.Monotonic(rand.Reader, 0)
}

func newULID() ulid.ULID {
	// TODO: This is unrealistic as it fixes to a particular time. Real applications would have a different scheme here.
	t := time.Unix(1000000, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

func New(prefix string) ID { return ID(prefix + fmt.Sprint(newULID())) }

func (u *ID) UnmarshalGQL(v interface{}) (err error) {
	s, _ := v.(string)
	*u = ID(s)
	return nil
}

func (u ID) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(string(u)))
}

func (u *ID) Scan(src interface{}) error {
	if src != nil {
		s, _ := src.(string)
		*u = ID(s)
	}
	return nil
}

func (u ID) Value() (driver.Value, error) {
	return string(u), nil
}
