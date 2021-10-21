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

package spec

// OpenAPI version 3.0.x is used.
const version = "3.0.3"

const JSON MediaType = "application/json"

const (
	InQuery Location = iota
	InHeader
	InPath
	InCookie
)

type (
	// Spec represents an OpenAPI Specification document https://swagger.io/specification/.
	//
	// Spec is the root struct representing an OAS document.
	Spec struct {
		Info         *Info            `json:"info"`
		Tags         []Tag            `json:"tags,omitempty"`
		Paths        map[string]*Path `json:"paths"`
		Components   *Components      `json:"components"`
		Security     Security         `json:"security,omitempty"`
		ExternalDocs *ExternalDocs    `json:"externalDocs,omitempty"`
	}
	// A Tag is used to group several Operation's.
	Tag struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}
	// The Info struct provides metadata about the API.
	Info struct {
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		TermsOfService string  `json:"termsOfService"`
		Contact        Contact `json:"contact"`
		License        License `json:"license"`
		Version        string  `json:"version"`
	}
	// Contact information for the exposed API.
	Contact struct {
		Name  string `json:"name,omitempty"`
		URL   string `json:"url,omitempty"`
		Email string `json:"email,omitempty"`
	}
	// License information for the exposed API.
	License struct {
		Name string `json:"name"`
		URL  string `json:"url,omitempty"`
	}
	// A Path describes the operations available on a single path.
	Path struct {
		Get        *Operation  `json:"get,omitempty"`
		Post       *Operation  `json:"post,omitempty"`
		Delete     *Operation  `json:"delete,omitempty"`
		Patch      *Operation  `json:"patch,omitempty"`
		Parameters []Parameter `json:"parameters,omitempty"`
	}
	// Parameter describes a single operation parameter.
	//
	// A unique parameter is defined by a combination of a Name and Location.
	Parameter struct {
		Name            string   `json:"name"`
		In              Location `json:"in"`
		Description     string   `json:"description,omitempty"`
		Required        bool     `json:"required,omitempty"`
		Deprecated      bool     `json:"deprecated,omitempty"`
		AllowEmptyValue bool     `json:"allowEmptyValue,omitempty"`
		Schema          *Type    `json:"schema"`
	}
	// A Location describes the location of a Parameter in a request object.
	//
	// Possible values are: InQuery, InHeader, InPath, InCookie.
	Location uint
	// Operation describes a single API operation on a Path.
	Operation struct {
		Summary      string                        `json:"summary,omitempty"`
		Description  string                        `json:"description,omitempty"`
		Tags         []string                      `json:"tags,omitempty"`
		ExternalDocs *ExternalDocs                 `json:"externalDocs,omitempty"`
		OperationID  string                        `json:"operationId"`
		Parameters   []*Parameter                  `json:"parameters,omitempty"`
		RequestBody  *RequestBody                  `json:"requestBody,omitempty"`
		Responses    map[string]*OperationResponse `json:"responses"`
		Deprecated   bool                          `json:"deprecated,omitempty"`
		Security     Security                      `json:"security,omitempty"`
	}
	// Security defines a security scheme that can be used by the operations.
	Security []map[string][]string
	// OperationResponse Describes a single response from an API Operation. Either Ref or Response must be given.
	OperationResponse struct {
		Ref      *Response
		Response *Response
	}
	// ExternalDocs allows referencing an external resource for extended documentation.
	ExternalDocs struct {
		Description string `json:"description"`
		URL         string `json:"url"`
	}
	// RequestBody describes a single request body.
	RequestBody struct {
		Description string  `json:"description"`
		Content     Content `json:"content"`
	}
	// Content of a RequestBody or a Response.
	Content map[MediaType]*MediaTypeObject
	// The MediaType of a MediaTypeObject.
	MediaType string
	// A MediaTypeObject provides schema and examples for the media type identified by its key.
	//
	// Currently, only JSON is supported.
	MediaTypeObject struct {
		Unique bool    `json:"-"`
		Ref    *Schema `json:"-"`
		Schema Schema  `json:"schema"`
	}
	// The Schema Object allows the definition of input and output data types
	Schema struct {
		Name   string
		Fields Fields
		Edges  Edges
	}
	// Fields of a Schema.
	Fields map[string]*Field
	// A Field defines a Field on a Schema.
	Field struct {
		*Type
		Unique   bool        `json:"-"`
		Required bool        `json:"-"`
		Example  interface{} `json:"example,omitempty"`
	}
	// The Type of a field. Primitive only.
	Type struct {
		Type   string `json:"type,omitempty"`
		Format string `json:"format,omitempty"`
		Items  *Type  `json:"items,omitempty"`
	}
	// Edges of a Schema.
	Edges map[string]*Edge
	// An Edge defines an Edge on a Schema. Must contain either Ref or Schema.
	Edge struct {
		Schema *Schema `json:"schema"`
		Ref    *Schema `json:"-"`
		Unique bool    `json:"-"`
	}
	// Response describes a single response from an API Operation.
	Response struct {
		Name        string               `json:"-"`
		Description string               `json:"description"`
		Headers     map[string]Parameter `json:"headers,omitempty"`
		Content     Content              `json:"content,omitempty"`
	}
	// Components holds a set of reusable objects for different aspects of the OAS.
	Components struct {
		Schemas         map[string]*Schema        `json:"schemas"`
		Responses       map[string]*Response      `json:"responses"`
		Parameters      map[string]*Parameter     `json:"parameters"`
		SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
	}
	// SecurityScheme defines a security scheme that can be used by the operations.
	SecurityScheme struct {
		Type             string      `json:"type"`
		Description      string      `json:"description,omitempty"`
		Name             string      `json:"name,omitempty"`
		In               string      `json:"in,omitempty"`
		Scheme           string      `json:"scheme,omitempty"`
		BearerFormat     string      `json:"bearerFormat,omitempty"`
		Flows            *OAuthFlows `json:"flows,omitempty"`
		OpenIDConnectURL string      `json:"openIDConnectURL,omitempty"`
	}
	// OAuthFlows allows configuration of the supported OAuth Flows.
	OAuthFlows struct {
		Implicit          *OAuthFlow `json:"implicit,omitempty"`
		Password          *OAuthFlow `json:"password,omitempty"`
		ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
		AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
	}
	// OAuthFlow configuration details for a supported OAuth Flow.
	OAuthFlow struct {
		AuthorizationURL string            `json:"authorizationURL,omitempty"`
		TokenURL         string            `json:"tokenURL,omitempty"`
		RefreshURL       string            `json:"refreshURL,omitempty"`
		Scopes           map[string]string `json:"scopes,omitempty"`
	}
)
