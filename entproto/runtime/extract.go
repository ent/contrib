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

package runtime

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ExtractTime returns the time.Time from a proto WKT Timestamp
func ExtractTime(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}

// MustExtractUUIDBytes returns the []byte representation of the uuid.UUID, if it fails it panics.
func MustExtractUUIDBytes(u uuid.UUID) []byte {
	b, err := u.MarshalBinary()
	if err != nil {
		panic("entproto: cannot marshal UUID to byte slice")
	}
	return b
}

// MustBytesToUUID returns a uuid.UUID from byte-slice b, if it fails it panics.
func MustBytesToUUID(b []byte) uuid.UUID {
	u, err := uuid.FromBytes(b)
	if err != nil {
		panic("entproto: cannot unmarshal UUID from bytes")
	}
	return u
}

// ValidateUUID returns an error if b is not a valid UUID.
func ValidateUUID(b []byte) error {
	_, err := uuid.FromBytes(b)
	return err
}
