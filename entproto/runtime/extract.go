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

package runtime

import (
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ExtractTime returns the time.Time from a proto WKT Timestamp
func ExtractTime(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}

// NewStruct creates a new *structpb.Struct from a map[string]interface{}.
// Returns nil if the input is nil or conversion fails.
func NewStruct(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil
	}
	return s
}

// ExtractStruct converts a *structpb.Struct to map[string]interface{}.
// Returns nil if the input is nil.
func ExtractStruct(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

// NewList creates a new *structpb.ListValue from a []interface{}.
// Returns nil if the input is nil or conversion fails.
func NewList(l []interface{}) *structpb.ListValue {
	if l == nil {
		return nil
	}
	lv, err := structpb.NewList(l)
	if err != nil {
		return nil
	}
	return lv
}

// ExtractList converts a *structpb.ListValue to []interface{}.
// Returns nil if the input is nil.
func ExtractList(l *structpb.ListValue) []interface{} {
	if l == nil {
		return nil
	}
	return l.AsSlice()
}
