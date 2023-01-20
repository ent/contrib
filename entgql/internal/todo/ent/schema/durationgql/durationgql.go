// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package durationgql

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalDuration(t time.Duration) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatInt(int64(t), 10))
	})
}

func UnmarshalDuration(v interface{}) (time.Duration, error) {
	switch v := v.(type) {
	case int64:
		return time.Duration(v), nil
	case string:
		return time.ParseDuration(v)
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return 0, err
		}
		return time.Duration(i), nil
	default:
		return 0, fmt.Errorf("invalid type %T, expect string", v)
	}
}
