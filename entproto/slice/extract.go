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
	"strings"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// InsertStrings returns a slice of proto strings from a slice of strings
func InsertStrings(s []string) []*wrapperspb.StringValue {
	wrappers := []*wrapperspb.StringValue{}
	for _, str := range s {
		wrappers = append(wrappers, wrapperspb.String(str))
	}
	return wrappers
}

// ExtractStrings returns a slice of strings from a slice of proto strings
func ExtractStrings(s []*wrapperspb.StringValue) []string {
	extract := []string{}
	for _, str := range s {
		cleanVal := strings.TrimPrefix(str.String(), "value:\"")
		cleanVal = strings.TrimSuffix(cleanVal, "\"")
		extract = append(extract, cleanVal)
	}
	return extract
}
