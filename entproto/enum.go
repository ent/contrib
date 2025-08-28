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

package entproto

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-viper/mapstructure/v2"

	"entgo.io/ent/entc/gen"
)

const (
	EnumAnnotation = "ProtoEnum"
)

var (
	ErrEnumFieldsNotAnnotated = errors.New("entproto: all Enum options must be covered with an entproto.Enum annotation")
	normalizeEnumIdent        = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

type EnumOption func(*enum)

// Enum configures the mapping between the ent Enum field and a protobuf Enum.
func Enum(vals map[string]int32, opts ...EnumOption) *enum {
	// apply options
	e := &enum{Options: vals}
	for _, op := range opts {
		op(e)
	}
	return e
}

// OmitFieldPrefix configures the Enum to omit the field name prefix from
// the enum labels on the generated protobuf message. Used for backwards
// compatibility with earlier versions of entproto where the field name
// wasn't prepended to the enum labels.
func OmitFieldPrefix() EnumOption {
	return func(e *enum) {
		e.OmitFieldPrefix = true
	}
}

type enum struct {
	Options         map[string]int32
	OmitFieldPrefix bool
}

func (*enum) Name() string {
	return EnumAnnotation
}

func (e *enum) findByNumber(n int32) string {
	for k, v := range e.Options {
		if v == n {
			return k
		}
	}
	return ""
}

func (e *enum) Verify(fld *gen.Field) error {
	// Verify that all fields on the Enum are in the annotation.
	if len(e.Options) != len(fld.Enums) {
		return ErrEnumFieldsNotAnnotated
	}
	pbIdentifiers := make(map[string]struct{}, len(fld.Enums))
	for _, opt := range fld.Enums {
		if _, ok := e.Options[opt.Value]; !ok {
			return fmt.Errorf("entproto: Enum option %s is not annotated with"+
				" a pbfield number using entproto.Enum", opt.Name)
		}
		pbIdent := NormalizeEnumIdentifier(opt.Value)
		if _, ok := pbIdentifiers[pbIdent]; ok {
			return fmt.Errorf("entproto: Enum option %q produces conflicting pbfield"+
				" name %q after normalization", opt.Name, pbIdent)
		}
		pbIdentifiers[pbIdent] = struct{}{}
	}

	// If default value is set on the pbfield, make sure it's option number is zero.
	if fld.Default {
		dv, ok := fld.DefaultValue().(string)
		if !ok {
			return fmt.Errorf("entproto: default value on Enum pbfield %s should be a string", fld.Name)
		}
		zeroField := e.findByNumber(0)
		if zeroField == "" {
			return fmt.Errorf("entproto: Enum pbfield %q has a default value but"+
				" entproto.Enum annotation doesn't contain an option with number 0", fld.Name)
		}
		if zeroField != dv {
			return fmt.Errorf(
				"entproto: default value for Enum pbfield %q is %q, but the proto annotation pbfield 0 is %q",
				fld.Name, dv, zeroField)
		}
	} else {
		// Make sure no one is using the zero option number.
		zeroField := e.findByNumber(0)
		if zeroField != "" {
			return fmt.Errorf("entproto: Enum pbfield %q has no default value but"+
				" entproto.Enum annotation contains an option with number 0", fld.Name)
		}
	}

	return nil
}

func extractEnumAnnotation(fld *gen.Field) (*enum, error) {
	annot, ok := fld.Annotations[EnumAnnotation]
	if !ok {
		return nil, fmt.Errorf("entproto: field %q does not have an entproto.Enum annotation", fld.Name)
	}

	var out enum
	err := mapstructure.Decode(annot, &out)
	if err != nil {
		return nil, fmt.Errorf("entproto: unable to decode entproto.Enum annotation for field %q: %w",
			fld.Name, err)
	}

	return &out, nil
}

// NormalizeEnumIdentifier normalizes the identifier of an enum pbfield
// to match the Proto Style Guide.
func NormalizeEnumIdentifier(s string) string {
	return strings.ToUpper(normalizeEnumIdent.ReplaceAllString(s, "_"))
}
