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

package entgql

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"text/template/parse"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/exp/slices"
)

var (
	// CollectionTemplate adds fields collection support using auto eager-load ent edges.
	// More info can be found here: https://spec.graphql.org/June2018/#sec-Field-Collection.
	CollectionTemplate = parseT("template/collection.tmpl")

	// EnumTemplate adds a template implementing MarshalGQL/UnmarshalGQL methods for enums.
	EnumTemplate = parseT("template/enum.tmpl")

	// NodeTemplate implements the Relay Node interface for all types.
	NodeTemplate = parseT("template/node.tmpl")

	// NodeDescriptorTemplate implements the Node descriptor API for all types.
	NodeDescriptorTemplate = parseT("template/node_descriptor.tmpl")

	// PaginationTemplate adds pagination support according to the GraphQL Cursor Connections Spec.
	// More info can be found in the following link: https://relay.dev/graphql/connections.htm.
	PaginationTemplate = parseT("template/pagination.tmpl")

	// TransactionTemplate adds support for ent.Client for opening transactions for the transaction
	// middleware. See transaction.go for for information.
	TransactionTemplate = parseT("template/transaction.tmpl")

	// EdgeTemplate adds edge resolution using eager-loading with a query fallback.
	EdgeTemplate = parseT("template/edge.tmpl")

	// WhereTemplate adds a template for generating <T>WhereInput filters for each schema type.
	WhereTemplate = parseT("template/where_input.tmpl")

	// MutationInputTemplate adds a template for generating Create<T>Input and Update<T>Input for each schema type.
	MutationInputTemplate = parseT("template/mutation_input.tmpl").SkipIf(skipMutationTemplate)

	// AllTemplates holds all templates for extending ent to support GraphQL.
	AllTemplates = []*gen.Template{
		CollectionTemplate,
		EnumTemplate,
		NodeTemplate,
		PaginationTemplate,
		TransactionTemplate,
		EdgeTemplate,
		MutationInputTemplate,
	}

	// TemplateFuncs contains the extra template functions used by entgql.
	TemplateFuncs = template.FuncMap{
		"fieldCollections":    fieldCollections,
		"fieldMapping":        fieldMapping,
		"filterEdges":         filterEdges,
		"filterFields":        filterFields,
		"filterNodes":         filterNodes,
		"gqlIDType":           gqlIDType,
		"gqlMarshaler":        gqlMarshaler,
		"gqlUnmarshaler":      gqlUnmarshaler,
		"hasWhereInput":       hasWhereInput,
		"isRelayConn":         isRelayConn,
		"isSkipMode":          isSkipMode,
		"mutationInputs":      mutationInputs,
		"nodeImplementors":    nodeImplementors,
		"nodeImplementorsVar": nodeImplementorsVar,
		"nodePaginationNames": nodePaginationNames,
		"orderFields":         orderFields,
		"skipMode":            skipModeFromString,
	}

	//go:embed template/*
	_templates embed.FS

	marshalerType   = reflect.TypeOf((*graphql.Marshaler)(nil)).Elem()
	unmarshalerType = reflect.TypeOf((*graphql.Unmarshaler)(nil)).Elem()
)

func parseT(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(TemplateFuncs).
		ParseFS(_templates, path))
}

// idType is returned by the gqlIDType below to describe the
// Go scalar type of the GraphQL ID. Note that, the type is
// not exported to avoid its usage outside the templates.
type idType struct {
	*field.TypeInfo
	// Mixed indicates if the ID type involves more than
	// single Go type and requires normalization to string.
	Mixed bool
}

// gqlIDType returns the scalar (Go) type of the GraphQL ID.
func gqlIDType(nodes []*gen.Type, defaultType *field.TypeInfo) (*idType, error) {
	if len(nodes) == 0 {
		return &idType{TypeInfo: defaultType}, nil
	}
	var mixed bool
	for i := 1; i < len(nodes); i++ {
		id1, id2 := nodes[i-1].ID, nodes[i].ID
		// Field type does not match.
		if mixed = id1.Type.Type != id2.Type.Type; mixed {
			break
		}
		// Underlying Go type does not match.
		if mixed = id1.HasGoType() != id2.HasGoType() || (id1.HasGoType() && id1.Type.RType.Ident != id2.Type.RType.Ident); mixed {
			break
		}
	}
	if !mixed {
		return &idType{TypeInfo: nodes[0].ID.Type}, nil
	}
	// If there are mixed types, expect all of them
	// to be either string or graphql.Marshaler.
	for _, n := range nodes {
		// Skip basic string types.
		if n.ID.IsString() && !n.ID.HasGoType() {
			continue
		}
		// Expect type to be un/marshaller to GraphQL scalar.
		if !n.ID.HasGoType() || !n.ID.Type.RType.Implements(marshalerType) || !n.ID.Type.RType.Implements(unmarshalerType) {
			return nil, errors.New("entgql: mixed id types must be type string or implement the graphql.Marshaller/graphql.Unmarshaller interfaces")
		}
	}
	return &idType{
		Mixed: true,
		TypeInfo: &field.TypeInfo{
			Type: field.TypeString,
		},
	}, nil
}

func gqlMarshaler(f *gen.Field) bool {
	return f.HasGoType() && f.Type.RType.Implements(marshalerType)
}

func gqlUnmarshaler(f *gen.Field) bool {
	return f.HasGoType() && f.Type.RType.Implements(unmarshalerType)
}

type fieldCollection struct {
	Edge    *gen.Edge
	Mapping []string
}

func fieldCollections(edges []*gen.Edge) ([]*fieldCollection, error) {
	collect := make([]*fieldCollection, 0, len(edges))
	for _, e := range edges {
		ant, err := annotation(e.Annotations)
		if err != nil {
			return nil, err
		}
		switch {
		case len(ant.Mapping) > 0:
			if !ant.Unbind {
				return nil, errors.New("bind and mapping annotations are mutually exclusive")
			}
			collect = append(collect, &fieldCollection{Edge: e, Mapping: ant.Mapping})
		case !ant.Unbind:
			mapping := []string{camel(e.Name)}
			collect = append(collect, &fieldCollection{Edge: e, Mapping: mapping})
		}
	}
	return collect, nil
}

// MutationDescriptor holds information about a GraphQL mutation input.
type MutationDescriptor struct {
	*gen.Type
	IsCreate bool
}

// Input returns the input's name.
func (m *MutationDescriptor) Input() (string, error) {
	gqlType, _, err := gqlTypeFromNode(m.Type)
	if err != nil {
		return "", err
	}
	if m.IsCreate {
		return fmt.Sprintf("Create%sInput", gqlType), nil
	}
	return fmt.Sprintf("Update%sInput", gqlType), nil
}

// Builders return the builder's names to apply the input.
func (m *MutationDescriptor) Builders() []string {
	if m.IsCreate {
		return []string{m.Type.CreateName()}
	}

	return []string{m.Type.UpdateName(), m.Type.UpdateOneName()}
}

// InputFieldDescriptor holds the information
// about a field in the input type.
// It's shared between GQL and Go types.
type InputFieldDescriptor struct {
	*gen.Field
	// AppendOp indicates if the field has the Append operator
	AppendOp bool
	// ClearOp indicates if the field has the Clear operator
	ClearOp bool
	// Nullable indicates if the field is nullable.
	Nullable bool
}

// IsPointer returns true if the Go type should be a pointer
func (f *InputFieldDescriptor) IsPointer() bool {
	if f.Type.Nillable || f.Type.RType.IsPtr() {
		return false
	}
	return f.Nullable
}

// InputFields returns the list of fields in the input type.
func (m *MutationDescriptor) InputFields() ([]*InputFieldDescriptor, error) {
	fields := make([]*InputFieldDescriptor, 0, len(m.Type.Fields))
	for _, f := range m.Type.Fields {
		ant, err := annotation(f.Annotations)
		if err != nil {
			return nil, err
		}
		if f.IsEdgeField() || m.skip(f.Immutable, ant.Skip) {
			continue
		}

		fields = append(fields, &InputFieldDescriptor{
			Field:    f,
			AppendOp: !m.IsCreate && f.SupportsMutationAppend(),
			ClearOp:  !m.IsCreate && f.Optional,
			Nullable: !m.IsCreate || f.Optional || f.Default || f.DefaultFunc(),
		})
	}

	return fields, nil
}

// InputEdges returns the list of fields in the input type.
//
// NOTE(giautm): This method should refactor to
// return a list of InputFieldDescriptor.
func (m *MutationDescriptor) InputEdges() ([]*gen.Edge, error) {
	edges := make([]*gen.Edge, 0, len(m.Type.Edges))
	for _, e := range m.Type.Edges {
		ant, err := annotation(e.Annotations)
		if err != nil {
			return nil, err
		}
		if e.Type.IsEdgeSchema() || m.skip(e.Immutable, ant.Skip) {
			continue
		}
		edges = append(edges, e)
	}
	return edges, nil
}

func (m *MutationDescriptor) skip(immutable bool, skip SkipMode) bool {
	if m.IsCreate {
		return skip.Is(SkipMutationCreateInput)
	}
	return immutable || skip.Is(SkipMutationUpdateInput)
}

// mutationInputs returns the list of input types for the mutation.
func mutationInputs(nodes []*gen.Type) ([]*MutationDescriptor, error) {
	filteredNodes := make([]*MutationDescriptor, 0, len(nodes))
	for _, n := range nodes {
		ant, err := annotation(n.Annotations)
		if err != nil {
			return nil, err
		}

		for _, a := range ant.MutationInputs {
			if (a.IsCreate && ant.Skip.Is(SkipMutationCreateInput)) ||
				(!a.IsCreate && ant.Skip.Is(SkipMutationUpdateInput)) {
				continue
			}

			filteredNodes = append(filteredNodes, &MutationDescriptor{
				Type:     n,
				IsCreate: a.IsCreate,
			})
		}
	}
	return filteredNodes, nil
}

// filterNodes filters out nodes that should not be included in the GraphQL schema.
func filterNodes(nodes []*gen.Type, skip SkipMode) ([]*gen.Type, error) {
	filteredNodes := make([]*gen.Type, 0, len(nodes))
	for _, n := range nodes {
		if n.HasCompositeID() {
			continue
		}
		ant, err := annotation(n.Annotations)
		if err != nil {
			return nil, err
		}
		if !ant.Skip.Is(skip) {
			filteredNodes = append(filteredNodes, n)
		}
	}
	return filteredNodes, nil
}

// filterEdges filters out edges that should not be included in the GraphQL schema.
func filterEdges(edges []*gen.Edge, skip SkipMode) ([]*gen.Edge, error) {
	filteredEdges := make([]*gen.Edge, 0, len(edges))
	for _, e := range edges {
		if e.Type.HasCompositeID() {
			continue
		}
		antE, err := annotation(e.Annotations)
		if err != nil {
			return nil, err
		}
		antT, err := annotation(e.Type.Annotations)
		if err != nil {
			return nil, err
		}
		if !antE.Skip.Is(skip) && !antT.Skip.Is(skip) {
			filteredEdges = append(filteredEdges, e)
		}
	}
	return filteredEdges, nil
}

// filterFields filters out fields that should not be included in the GraphQL schema.
func filterFields(fields []*gen.Field, skip SkipMode) ([]*gen.Field, error) {
	filteredFields := make([]*gen.Field, 0, len(fields))
	for _, f := range fields {
		ant, err := annotation(f.Annotations)
		if err != nil {
			return nil, err
		}
		if !ant.Skip.Is(skip) {
			filteredFields = append(filteredFields, f)
		}
	}
	return filteredFields, nil
}

// OrderTerm is a struct that represents a single GraphQL order term.
type OrderTerm struct {
	// The type that owns the order field.
	Owner *gen.Type
	// The GraphQL name of the field.
	GQL string
	// The type that owns the field. For type fields, it equals to Owner.
	// For edge fields, it equals to the underlying edge's type.
	Type *gen.Type
	// Not nil if it is a type/edge field.
	Field *gen.Field
	// Not nil if it is an edge field or count.
	Edge *gen.Edge
	// True if it is a count field.
	Count bool
}

// IsFieldTerm returns true if the order term is a type field term.
func (o *OrderTerm) IsFieldTerm() bool {
	return o.Field != nil && o.Edge == nil
}

// IsEdgeFieldTerm returns true if the order term is an edge field term.
func (o *OrderTerm) IsEdgeFieldTerm() bool {
	return o.Field != nil && o.Edge != nil
}

// IsEdgeCountTerm returns true if the order term is an edge count term.
func (o *OrderTerm) IsEdgeCountTerm() bool {
	return o.Field == nil && o.Edge != nil && o.Count
}

// VarName returns the name of the variable holding the order term.
func (o *OrderTerm) VarName() (string, error) {
	switch prefix := paginationNames(o.Owner.Name).OrderField; {
	case o.IsFieldTerm():
		return prefix + o.Field.StructField(), nil
	case o.IsEdgeFieldTerm():
		return prefix + o.Edge.StructField() + o.Field.StructField(), nil
	case o.IsEdgeCountTerm():
		return prefix + o.Edge.StructField() + "Count", nil
	default:
		return "", fmt.Errorf("entgql: invalid order term %v", o)
	}
}

// VarField returns the field name inside the variable holding the order term.
func (o *OrderTerm) VarField() (string, error) {
	switch {
	case o.IsFieldTerm():
		return fmt.Sprintf("%s.%s", o.Type.Package(), o.Field.Constant()), nil
	case o.IsEdgeFieldTerm(), o.IsEdgeCountTerm():
		return strconv.Quote(strings.ToLower(o.GQL)), nil
	default:
		return "", fmt.Errorf("entgql: invalid order term %v", o)
	}
}

// orderFields returns the GraphQL fields of the given node with the `OrderField` annotation.
func orderFields(n *gen.Type) ([]*OrderTerm, error) {
	var (
		terms  []*OrderTerm
		fields = n.Fields
	)
	if n.HasOneFieldID() {
		fields = append([]*gen.Field{n.ID}, fields...)
	}
	for _, f := range fields {
		switch ant, err := annotation(f.Annotations); {
		case err != nil:
			return nil, err
		case ant.Skip.Is(SkipOrderField), ant.OrderField == "":
		case !f.Type.Comparable():
			return nil, fmt.Errorf("entgql: ordered field %s.%s must be comparable", n.Name, f.Name)
		default:
			terms = append(terms, &OrderTerm{
				Owner: n,
				GQL:   ant.OrderField,
				Type:  n,
				Field: f,
			})
		}
	}
	for _, e := range n.Edges {
		name := strings.ToUpper(e.Name)
		switch ant, err := annotation(e.Annotations); {
		case err != nil:
			return nil, err
		case ant.Skip.Is(SkipOrderField), ant.OrderField == "":
		case ant.OrderField == fmt.Sprintf("%s_COUNT", name):
			// Validate that the edge has a count ordering.
			if _, err := e.OrderCountName(); err != nil {
				return nil, fmt.Errorf("entgql: invalid order field %s defined on edge %s.%s: %w", ant.OrderField, n.Name, e.Name, err)
			}
			terms = append(terms, &OrderTerm{
				Owner: n,
				GQL:   ant.OrderField,
				Type:  n,
				Edge:  e,
				Count: true,
			})
		case strings.HasPrefix(ant.OrderField, name+"_"):
			// Validate that the edge has a edge field ordering.
			if _, err := e.OrderFieldName(); err != nil {
				return nil, fmt.Errorf("entgql: invalid order field %s defined on edge %s.%s: %w", ant.OrderField, n.Name, e.Name, err)
			}
			ef := strings.TrimPrefix(ant.OrderField, name+"_")
			idx := slices.IndexFunc(e.Type.Fields, func(f *gen.Field) bool {
				ant, err := annotation(f.Annotations)
				return err == nil && ant.OrderField == ef
			})
			if idx == -1 {
				return nil, fmt.Errorf("entgql: order field %s defined on edge %s.%s was not found on its reference", ant.OrderField, n.Name, e.Name)
			}
			terms = append(terms, &OrderTerm{
				Owner: n,
				GQL:   ant.OrderField,
				Edge:  e,
				Type:  e.Type,
				Field: e.Type.Fields[idx],
			})
		default:
			return nil, fmt.Errorf("entgql: invalid order field defined on edge %s.%s", n.Name, e.Name)
		}
	}
	return terms, nil
}

// hasWhereInput returns true if neither the edge nor its
// node type has the SkipWhereInput annotation
func hasWhereInput(n *gen.Edge) (v bool, err error) {
	antEdge, err := annotation(n.Annotations)
	if err != nil || antEdge.Skip.Is(SkipWhereInput) {
		return false, err
	}
	ant, err := annotation(n.Type.Annotations)
	if err != nil || ant.Skip.Is(SkipWhereInput) {
		return false, err
	}
	return true, nil
}

// skipModeFromString returns SkipFlag from a string
func skipModeFromString(modes ...string) (SkipMode, error) {
	var m SkipMode
	for _, s := range modes {
		switch s {
		case "type":
			m |= SkipType
		case "enum_field":
			m |= SkipEnumField
		case "order_field":
			m |= SkipOrderField
		case "where_input":
			m |= SkipWhereInput
		case "mutation_create_input":
			m |= SkipMutationCreateInput
		case "mutation_update_input":
			m |= SkipMutationUpdateInput
		default:
			return 0, fmt.Errorf("invalid skip mode: %s", s)
		}
	}
	return m, nil
}

func isSkipMode(antSkip interface{}, m string) (bool, error) {
	skip, err := skipModeFromString(m)
	if err != nil || antSkip == nil {
		return false, err
	}
	if raw, ok := antSkip.(float64); ok {
		return SkipMode(raw).Is(skip), nil
	}
	return false, fmt.Errorf("invalid annotation skip: %v", antSkip)
}

func isRelayConn(e *gen.Edge) (bool, error) {
	ant, err := annotation(e.Annotations)
	if err != nil {
		return false, err
	}
	return ant.RelayConnection, nil
}

// PaginationNames holds the names of the pagination fields.
type PaginationNames struct {
	Connection string
	Edge       string
	Node       string
	Order      string
	OrderField string
	WhereInput string
}

func (p *PaginationNames) TypeDefs() []*ast.Definition {
	return []*ast.Definition{
		{
			Name:        p.Edge,
			Kind:        ast.Object,
			Description: "An edge in a connection.",
			Fields: []*ast.FieldDefinition{
				{
					Name:        "node",
					Type:        ast.NamedType(p.Node, nil),
					Description: "The item at the end of the edge.",
				},
				{
					Name:        "cursor",
					Type:        ast.NonNullNamedType(RelayCursor, nil),
					Description: "A cursor for use in pagination.",
				},
			},
		},
		{
			Name:        p.Connection,
			Kind:        ast.Object,
			Description: "A connection to a list of items.",
			Fields: []*ast.FieldDefinition{
				{
					Name:        "edges",
					Type:        ast.ListType(ast.NamedType(p.Edge, nil), nil),
					Description: "A list of edges.",
				},
				{
					Name:        "pageInfo",
					Type:        ast.NonNullNamedType(RelayPageInfo, nil),
					Description: "Information to aid in pagination.",
				},
				{
					Name:        "totalCount",
					Type:        ast.NonNullNamedType("Int", nil),
					Description: "Identifies the total count of items in the connection.",
				},
			},
		},
	}
}

func (p *PaginationNames) OrderInputDef() *ast.Definition {
	return &ast.Definition{
		Name:        p.Order,
		Kind:        ast.InputObject,
		Description: fmt.Sprintf("Ordering options for %s connections", p.Node),
		Fields: ast.FieldList{
			{
				Name: "direction",
				Type: ast.NonNullNamedType(OrderDirectionEnum, nil),
				DefaultValue: &ast.Value{
					Raw:  "ASC",
					Kind: ast.EnumValue,
				},
				Description: "The ordering direction.",
			},
			{
				Name:        "field",
				Type:        ast.NonNullNamedType(p.OrderField, nil),
				Description: fmt.Sprintf("The field by which to order %s.", plural(p.Node)),
			},
		},
	}
}

func (p *PaginationNames) ConnectionField(name string, hasOrderBy, multiOrder, hasWhereInput bool) *ast.FieldDefinition {
	def := &ast.FieldDefinition{
		Name: name,
		Type: ast.NonNullNamedType(p.Connection, nil),
		Arguments: ast.ArgumentDefinitionList{
			{
				Name:        "after",
				Type:        ast.NamedType(RelayCursor, nil),
				Description: "Returns the elements in the list that come after the specified cursor.",
			},
			{
				Name:        "first",
				Type:        ast.NamedType("Int", nil),
				Description: "Returns the first _n_ elements from the list.",
			},
			{
				Name:        "before",
				Type:        ast.NamedType(RelayCursor, nil),
				Description: "Returns the elements in the list that come before the specified cursor.",
			},
			{
				Name:        "last",
				Type:        ast.NamedType("Int", nil),
				Description: "Returns the last _n_ elements from the list.",
			},
		},
	}
	if hasOrderBy {
		orderT := ast.NamedType(p.Order, nil)
		if multiOrder {
			orderT = ast.ListType(ast.NonNullNamedType(p.Order, nil), nil)
		}
		def.Arguments = append(def.Arguments, &ast.ArgumentDefinition{
			Name:        "orderBy",
			Type:        orderT,
			Description: fmt.Sprintf("Ordering options for %s returned from the connection.", plural(p.Node)),
		})
	}
	if hasWhereInput {
		def.Arguments = append(def.Arguments, &ast.ArgumentDefinition{
			Name:        "where",
			Type:        ast.NamedType(p.WhereInput, nil),
			Description: fmt.Sprintf("Filtering options for %s returned from the connection.", plural(p.Node)),
		})
	}

	return def
}

func gqlTypeFromNode(t *gen.Type) (gqlType string, ant *Annotation, err error) {
	if ant, err = annotation(t.Annotations); err != nil {
		return
	}
	gqlType = t.Name
	if ant.Type != "" {
		gqlType = ant.Type
	}
	return
}

// nodePaginationNames returns the names of the pagination types for the node.
func nodePaginationNames(t *gen.Type) (*PaginationNames, error) {
	node, _, err := gqlTypeFromNode(t)
	if err != nil {
		return nil, err
	}

	return paginationNames(node), nil
}

func paginationNames(node string) *PaginationNames {
	return &PaginationNames{
		Connection: fmt.Sprintf("%sConnection", node),
		Edge:       fmt.Sprintf("%sEdge", node),
		Node:       node,
		Order:      fmt.Sprintf("%sOrder", node),
		OrderField: fmt.Sprintf("%sOrderField", node),
		WhereInput: fmt.Sprintf("%sWhereInput", node),
	}
}

// removeOldAssets removes files that were generated before v0.1.0.
func removeOldAssets(next gen.Generator) gen.Generator {
	const prefix = "gql_"
	templates := []*gen.Template{WhereTemplate}
	templates = append(templates, AllTemplates...)
	return gen.GenerateFunc(func(g *gen.Graph) error {
		for _, rootT := range templates {
			for _, t := range rootT.Templates() {
				if parse.IsEmptyTree(t.Root) {
					continue
				}
				if !strings.HasPrefix(t.Name(), prefix) {
					continue
				}
				name := strings.TrimPrefix(t.Name(), prefix)
				if err := removeOldTemplate(g, name); err != nil {
					return err
				}
			}
		}
		return next.Generate(g)
	})
}

func removeOldTemplate(g *gen.Graph, name string) error {
	// Check if name already taken by existing schema field.
	for _, n := range g.Nodes {
		if n.Package() == name {
			return nil
		}
	}
	err := os.Remove(filepath.Join(g.Target, name+".go"))
	if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func skipMutationTemplate(g *gen.Graph) bool {
	for _, n := range g.Nodes {
		ant, err := annotation(n.Annotations)
		if err != nil {
			continue
		}
		for _, i := range ant.MutationInputs {
			if (i.IsCreate && !ant.Skip.Is(SkipMutationCreateInput)) ||
				(!i.IsCreate && !ant.Skip.Is(SkipMutationUpdateInput)) {
				return false
			}
		}
	}
	return true
}

func nodeImplementors(n *gen.Type) (ifaces []string, err error) {
	ant, err := annotation(n.Annotations)
	if err != nil {
		return nil, err
	}
	if !ant.Skip.Is(SkipType) && !slices.Contains(ant.Implements, "Node") {
		ifaces = append(ifaces, "Node")
	}
	return append(ifaces, ant.Implements...), nil
}

func nodeImplementorsVar(n *gen.Type) string {
	return strings.ToLower(n.Name) + "Implementors"
}
