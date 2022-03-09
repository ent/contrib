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

package bigintgql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

type BigInt struct {
	*big.Int
}

func NewBigInt(i int64) BigInt {
	return BigInt{Int: big.NewInt(i)}
}

func (b *BigInt) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return err
	}
	if !i.Valid {
		return nil
	}
	if b.Int == nil {
		b.Int = big.NewInt(0)
	}
	// Value came in a floating point format.
	if strings.ContainsAny(i.String, ".+e") {
		f := big.NewFloat(0)
		if _, err := fmt.Sscan(i.String, f); err != nil {
			return err
		}
		b.Int, _ = f.Int(b.Int)
	} else if _, err := fmt.Sscan(i.String, b.Int); err != nil {
		return err
	}
	return nil
}

func (b BigInt) Value() (driver.Value, error) {
	return b.String(), nil
}

func (b BigInt) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(b.String()))
}

func (b *BigInt) UnmarshalGQL(v interface{}) error {
	if bi, ok := v.(string); ok {
		b.Int = new(big.Int)
		b.Int, ok = b.Int.SetString(bi, 10)
		if !ok {
			return fmt.Errorf("invalid big number: %s", bi)
		}

		return nil
	}

	return fmt.Errorf("invalid big number")
}
