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

package slice

import (
	"google.golang.org/protobuf/types/known/structpb"
)

// InsertStrings returns a proto list from a slice of strings
func InsertStrings(s []string) *structpb.ListValue {
	insert := make([]interface{}, 0)
	for _, str := range s {
		insert = append(insert, str)
	}
	wrapper, err := structpb.NewList(insert)
	if err != nil {
		return nil
	}
	return wrapper
}

// ExtractStrings returns a slice of strings from a proto list
func ExtractStrings(s *structpb.ListValue) []string {
	extract := []string{}
	for _, str := range s.AsSlice() {
		extract = append(extract, str.(string))
	}
	if len(extract) == 0 {
		return nil
	}
	return extract
}
