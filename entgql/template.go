package entgql

import (
	"text/template"

	"github.com/facebookincubator/ent-contrib/entgql/internal"

	"github.com/facebook/ent/entc/gen"
	_ "github.com/go-bindata/go-bindata"
)

var (
	// CollectionTemplate adds fields collection support using auto eager-load ent edges.
	// More info can be found here: https://spec.graphql.org/June2018/#sec-Field-Collection.
	CollectionTemplate = parse("template/collection.tmpl")

	// EnumTemplate adds a template implementing MarshalGQL/UnmarshalGQL methods for enums.
	EnumTemplate = parse("template/enum.tmpl")

	// NodeTemplate implements the Relay Node interface for all types.
	NodeTemplate = parse("template/node.tmpl")

	// PaginationTemplate adds pagination support according to the GraphQL Cursor Connections Spec.
	// More info can be found in the following link: https://relay.dev/graphql/connections.htm.
	PaginationTemplate = parse("template/pagination.tmpl")

	// TransactionTemplate adds support for ent.Client for opening transactions for the transaction
	// middleware. See transaction.go for for information.
	TransactionTemplate = parse("template/transaction.tmpl")

	// AllTemplates holds all templates for extending ent to support GraphQL.
	AllTemplates = []*template.Template{
		CollectionTemplate,
		EnumTemplate,
		NodeTemplate,
		PaginationTemplate,
		TransactionTemplate,
	}
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -o=internal/bindata.go -pkg=internal -modtime=1 ./template

func parse(path string) *template.Template {
	text := string(internal.MustAsset(path))
	return template.Must(template.New(path).
		Funcs(gen.Funcs).
		Parse(text))
}
