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

package uuidgql

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func New() UUID { return UUID(uuid.New()) }

func (u *UUID) UnmarshalGQL(v interface{}) (err error) {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid type %T, expect string", v)
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	*u = UUID(id)
	return nil
}

func (u UUID) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(uuid.UUID(u).String()))
}

func (u *UUID) Scan(src interface{}) error {
	if err := (*uuid.UUID)(u).Scan(src); err != nil {
		return err
	}
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).Value()
}
