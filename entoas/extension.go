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

package entoas

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"entgo.io/contrib/entoas/spec"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	// Config provides global metadata for the generator. It is injected into the gen.Graph.
	Config struct {
		// DefaultPolicy defines the default policy for endpoint generation.
		// It is used if no policy is set on a (sub-)resource.
		// Defaults to PolicyExpose.
		DefaultPolicy Policy
		// Whether or whether not to generate simple models instead of one model per endpoint.
		//
		// The OAS generator by default creates one view per endpoint. The naming strategy is the following:
		// - <S><Op> for 1st level operation Op on schema S
		// - <S><Op>_<E> for an eager loaded edge E on 1st level operation Op on schema S
		// - <S>_<E><Op> for a 2nd level operation Op on edge E on schema S
		//
		// By enabling she SimpleModels configuration the generator simply adds the defined schemas with all fields and edges.
		// Serialization groups have no effects in this mode.
		SimpleModels bool
	}
	// Extension implements entc.Extension interface for providing OpenAPI Specification generation.
	Extension struct {
		entc.DefaultExtension
		config    *Config
		mutations []MutateFunc
		out       io.Writer
	}
	// ExtensionOption allows managing Extension configuration using functional arguments.
	ExtensionOption func(*Extension) error
	// MutateFunc defines a closure to be called on a generated spec.
	MutateFunc func(*gen.Graph, *spec.Spec) error
)

// NewExtension returns a new entoas extension with default values.
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{config: &Config{DefaultPolicy: PolicyExpose}}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	return ex, nil
}

// Hooks of the Extension.
func (ex *Extension) Hooks() []gen.Hook {
	return []gen.Hook{ex.generate}
}

// Annotations of the extensions.
func (ex *Extension) Annotations() []entc.Annotation {
	return []entc.Annotation{ex.config}
}

// DefaultPolicy sets the default ExclusionPolicy to use of none is given on a (sub-)schema.
func DefaultPolicy(p Policy) ExtensionOption {
	return func(ex *Extension) error {
		ex.config.DefaultPolicy = p
		return nil
	}
}

// Mutations adds the given mutations to the spec generator.
//
// A MutateFunc is a simple closure that can be used to edit the generated spec.
// It can be used to add custom endpoints or alter the spec in any other way.
func Mutations(ms ...MutateFunc) ExtensionOption {
	return func(ex *Extension) error {
		ex.mutations = append(ex.mutations, ms...)
		return nil
	}
}

// SimpleModels enables the simple model generation feature.
//
// Further information can be found at Config.SimpleModels.
func SimpleModels() ExtensionOption {
	return func(ex *Extension) error {
		ex.config.SimpleModels = true
		return nil
	}
}

// SpecTitle sets the title of the Info block.
func SpecTitle(v string) ExtensionOption {
	return Mutations(func(_ *gen.Graph, spec *spec.Spec) error {
		spec.Info.Title = v
		return nil
	})
}

// SpecDescription sets the title of the Info block.
func SpecDescription(v string) ExtensionOption {
	return Mutations(func(_ *gen.Graph, spec *spec.Spec) error {
		spec.Info.Description = v
		return nil
	})
}

// SpecVersion sets the version of the Info block.
func SpecVersion(v string) ExtensionOption {
	return Mutations(func(_ *gen.Graph, spec *spec.Spec) error {
		spec.Info.Version = v
		return nil
	})
}

// WriteTo writes the current specs content to the given io.Writer.
func WriteTo(out io.Writer) ExtensionOption {
	return func(ex *Extension) error {
		ex.out = out
		return nil
	}
}

// generator returns a gen.Hook that creates a Spec for the given gen.Graph.
func (ex *Extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		// Let ent create all the files.
		if err := next.Generate(g); err != nil {
			return err
		}
		// Spec stub to fill.
		s := &spec.Spec{
			Info: &spec.Info{
				Title:       "Ent Schema API",
				Description: "This is an auto generated API description made out of an Ent schema definition",
				Version:     "0.0.0",
			},
			Components: &spec.Components{
				Schemas:    make(map[string]*spec.Schema),
				Responses:  make(map[string]*spec.Response),
				Parameters: make(map[string]*spec.Parameter),
			},
		}
		// Run the generator.
		if err := generate(g, s); err != nil {
			return err
		}
		// Run the user provided mutations.
		for _, m := range ex.mutations {
			if err := m(g, s); err != nil {
				return err
			}
		}
		// Dump the spec.
		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		if ex.out != nil {
			_, err = ex.out.Write(b)
			return err
		}
		return os.WriteFile(filepath.Join(g.Target, "openapi.json"), b, 0644)
	})
}

// Name implements entc.Annotation interface.
func (c Config) Name() string {
	return "EntOASConfig"
}

// Decode from ent.
func (c *Config) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, c)
}

// GetConfig loads the entoas.Config from the given *gen.Config object.
func GetConfig(cfg *gen.Config) (*Config, error) {
	c := &Config{}
	if cfg == nil && cfg.Annotations == nil && cfg.Annotations[c.Name()] == nil {
		return nil, errors.New("entoas extension configuration not found")
	}
	return c, c.Decode(cfg.Annotations[c.Name()])
}

var (
	_ entc.Extension  = (*Extension)(nil)
	_ entc.Annotation = (*Config)(nil)
)
