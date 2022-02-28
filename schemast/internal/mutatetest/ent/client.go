// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/migrate"

	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/user"
	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/withfields"
	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/withmodifiedfield"
	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/withnilfields"
	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/withoutfields"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// User is the client for interacting with the User builders.
	User *UserClient
	// WithFields is the client for interacting with the WithFields builders.
	WithFields *WithFieldsClient
	// WithModifiedField is the client for interacting with the WithModifiedField builders.
	WithModifiedField *WithModifiedFieldClient
	// WithNilFields is the client for interacting with the WithNilFields builders.
	WithNilFields *WithNilFieldsClient
	// WithoutFields is the client for interacting with the WithoutFields builders.
	WithoutFields *WithoutFieldsClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.User = NewUserClient(c.config)
	c.WithFields = NewWithFieldsClient(c.config)
	c.WithModifiedField = NewWithModifiedFieldClient(c.config)
	c.WithNilFields = NewWithNilFieldsClient(c.config)
	c.WithoutFields = NewWithoutFieldsClient(c.config)
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:               ctx,
		config:            cfg,
		User:              NewUserClient(cfg),
		WithFields:        NewWithFieldsClient(cfg),
		WithModifiedField: NewWithModifiedFieldClient(cfg),
		WithNilFields:     NewWithNilFieldsClient(cfg),
		WithoutFields:     NewWithoutFieldsClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		config:            cfg,
		User:              NewUserClient(cfg),
		WithFields:        NewWithFieldsClient(cfg),
		WithModifiedField: NewWithModifiedFieldClient(cfg),
		WithNilFields:     NewWithNilFieldsClient(cfg),
		WithoutFields:     NewWithoutFieldsClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		User.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.User.Use(hooks...)
	c.WithFields.Use(hooks...)
	c.WithModifiedField.Use(hooks...)
	c.WithNilFields.Use(hooks...)
	c.WithoutFields.Use(hooks...)
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Create returns a create builder for User.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of User entities.
func (c *UserClient) CreateBulk(builders ...*UserCreate) *UserCreateBulk {
	return &UserCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUser(u))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id int) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUserID(id))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *UserClient) DeleteOneID(id int) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Query returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{
		config: c.config,
	}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id int) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id int) *User {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	return c.hooks.User
}

// WithFieldsClient is a client for the WithFields schema.
type WithFieldsClient struct {
	config
}

// NewWithFieldsClient returns a client for the WithFields from the given config.
func NewWithFieldsClient(c config) *WithFieldsClient {
	return &WithFieldsClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `withfields.Hooks(f(g(h())))`.
func (c *WithFieldsClient) Use(hooks ...Hook) {
	c.hooks.WithFields = append(c.hooks.WithFields, hooks...)
}

// Create returns a create builder for WithFields.
func (c *WithFieldsClient) Create() *WithFieldsCreate {
	mutation := newWithFieldsMutation(c.config, OpCreate)
	return &WithFieldsCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of WithFields entities.
func (c *WithFieldsClient) CreateBulk(builders ...*WithFieldsCreate) *WithFieldsCreateBulk {
	return &WithFieldsCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for WithFields.
func (c *WithFieldsClient) Update() *WithFieldsUpdate {
	mutation := newWithFieldsMutation(c.config, OpUpdate)
	return &WithFieldsUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WithFieldsClient) UpdateOne(wf *WithFields) *WithFieldsUpdateOne {
	mutation := newWithFieldsMutation(c.config, OpUpdateOne, withWithFields(wf))
	return &WithFieldsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *WithFieldsClient) UpdateOneID(id int) *WithFieldsUpdateOne {
	mutation := newWithFieldsMutation(c.config, OpUpdateOne, withWithFieldsID(id))
	return &WithFieldsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WithFields.
func (c *WithFieldsClient) Delete() *WithFieldsDelete {
	mutation := newWithFieldsMutation(c.config, OpDelete)
	return &WithFieldsDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WithFieldsClient) DeleteOne(wf *WithFields) *WithFieldsDeleteOne {
	return c.DeleteOneID(wf.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WithFieldsClient) DeleteOneID(id int) *WithFieldsDeleteOne {
	builder := c.Delete().Where(withfields.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WithFieldsDeleteOne{builder}
}

// Query returns a query builder for WithFields.
func (c *WithFieldsClient) Query() *WithFieldsQuery {
	return &WithFieldsQuery{
		config: c.config,
	}
}

// Get returns a WithFields entity by its id.
func (c *WithFieldsClient) Get(ctx context.Context, id int) (*WithFields, error) {
	return c.Query().Where(withfields.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WithFieldsClient) GetX(ctx context.Context, id int) *WithFields {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *WithFieldsClient) Hooks() []Hook {
	return c.hooks.WithFields
}

// WithModifiedFieldClient is a client for the WithModifiedField schema.
type WithModifiedFieldClient struct {
	config
}

// NewWithModifiedFieldClient returns a client for the WithModifiedField from the given config.
func NewWithModifiedFieldClient(c config) *WithModifiedFieldClient {
	return &WithModifiedFieldClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `withmodifiedfield.Hooks(f(g(h())))`.
func (c *WithModifiedFieldClient) Use(hooks ...Hook) {
	c.hooks.WithModifiedField = append(c.hooks.WithModifiedField, hooks...)
}

// Create returns a create builder for WithModifiedField.
func (c *WithModifiedFieldClient) Create() *WithModifiedFieldCreate {
	mutation := newWithModifiedFieldMutation(c.config, OpCreate)
	return &WithModifiedFieldCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of WithModifiedField entities.
func (c *WithModifiedFieldClient) CreateBulk(builders ...*WithModifiedFieldCreate) *WithModifiedFieldCreateBulk {
	return &WithModifiedFieldCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for WithModifiedField.
func (c *WithModifiedFieldClient) Update() *WithModifiedFieldUpdate {
	mutation := newWithModifiedFieldMutation(c.config, OpUpdate)
	return &WithModifiedFieldUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WithModifiedFieldClient) UpdateOne(wmf *WithModifiedField) *WithModifiedFieldUpdateOne {
	mutation := newWithModifiedFieldMutation(c.config, OpUpdateOne, withWithModifiedField(wmf))
	return &WithModifiedFieldUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *WithModifiedFieldClient) UpdateOneID(id int) *WithModifiedFieldUpdateOne {
	mutation := newWithModifiedFieldMutation(c.config, OpUpdateOne, withWithModifiedFieldID(id))
	return &WithModifiedFieldUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WithModifiedField.
func (c *WithModifiedFieldClient) Delete() *WithModifiedFieldDelete {
	mutation := newWithModifiedFieldMutation(c.config, OpDelete)
	return &WithModifiedFieldDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WithModifiedFieldClient) DeleteOne(wmf *WithModifiedField) *WithModifiedFieldDeleteOne {
	return c.DeleteOneID(wmf.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WithModifiedFieldClient) DeleteOneID(id int) *WithModifiedFieldDeleteOne {
	builder := c.Delete().Where(withmodifiedfield.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WithModifiedFieldDeleteOne{builder}
}

// Query returns a query builder for WithModifiedField.
func (c *WithModifiedFieldClient) Query() *WithModifiedFieldQuery {
	return &WithModifiedFieldQuery{
		config: c.config,
	}
}

// Get returns a WithModifiedField entity by its id.
func (c *WithModifiedFieldClient) Get(ctx context.Context, id int) (*WithModifiedField, error) {
	return c.Query().Where(withmodifiedfield.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WithModifiedFieldClient) GetX(ctx context.Context, id int) *WithModifiedField {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryOwner queries the owner edge of a WithModifiedField.
func (c *WithModifiedFieldClient) QueryOwner(wmf *WithModifiedField) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := wmf.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(withmodifiedfield.Table, withmodifiedfield.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, withmodifiedfield.OwnerTable, withmodifiedfield.OwnerColumn),
		)
		fromV = sqlgraph.Neighbors(wmf.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *WithModifiedFieldClient) Hooks() []Hook {
	return c.hooks.WithModifiedField
}

// WithNilFieldsClient is a client for the WithNilFields schema.
type WithNilFieldsClient struct {
	config
}

// NewWithNilFieldsClient returns a client for the WithNilFields from the given config.
func NewWithNilFieldsClient(c config) *WithNilFieldsClient {
	return &WithNilFieldsClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `withnilfields.Hooks(f(g(h())))`.
func (c *WithNilFieldsClient) Use(hooks ...Hook) {
	c.hooks.WithNilFields = append(c.hooks.WithNilFields, hooks...)
}

// Create returns a create builder for WithNilFields.
func (c *WithNilFieldsClient) Create() *WithNilFieldsCreate {
	mutation := newWithNilFieldsMutation(c.config, OpCreate)
	return &WithNilFieldsCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of WithNilFields entities.
func (c *WithNilFieldsClient) CreateBulk(builders ...*WithNilFieldsCreate) *WithNilFieldsCreateBulk {
	return &WithNilFieldsCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for WithNilFields.
func (c *WithNilFieldsClient) Update() *WithNilFieldsUpdate {
	mutation := newWithNilFieldsMutation(c.config, OpUpdate)
	return &WithNilFieldsUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WithNilFieldsClient) UpdateOne(wnf *WithNilFields) *WithNilFieldsUpdateOne {
	mutation := newWithNilFieldsMutation(c.config, OpUpdateOne, withWithNilFields(wnf))
	return &WithNilFieldsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *WithNilFieldsClient) UpdateOneID(id int) *WithNilFieldsUpdateOne {
	mutation := newWithNilFieldsMutation(c.config, OpUpdateOne, withWithNilFieldsID(id))
	return &WithNilFieldsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WithNilFields.
func (c *WithNilFieldsClient) Delete() *WithNilFieldsDelete {
	mutation := newWithNilFieldsMutation(c.config, OpDelete)
	return &WithNilFieldsDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WithNilFieldsClient) DeleteOne(wnf *WithNilFields) *WithNilFieldsDeleteOne {
	return c.DeleteOneID(wnf.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WithNilFieldsClient) DeleteOneID(id int) *WithNilFieldsDeleteOne {
	builder := c.Delete().Where(withnilfields.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WithNilFieldsDeleteOne{builder}
}

// Query returns a query builder for WithNilFields.
func (c *WithNilFieldsClient) Query() *WithNilFieldsQuery {
	return &WithNilFieldsQuery{
		config: c.config,
	}
}

// Get returns a WithNilFields entity by its id.
func (c *WithNilFieldsClient) Get(ctx context.Context, id int) (*WithNilFields, error) {
	return c.Query().Where(withnilfields.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WithNilFieldsClient) GetX(ctx context.Context, id int) *WithNilFields {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *WithNilFieldsClient) Hooks() []Hook {
	return c.hooks.WithNilFields
}

// WithoutFieldsClient is a client for the WithoutFields schema.
type WithoutFieldsClient struct {
	config
}

// NewWithoutFieldsClient returns a client for the WithoutFields from the given config.
func NewWithoutFieldsClient(c config) *WithoutFieldsClient {
	return &WithoutFieldsClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `withoutfields.Hooks(f(g(h())))`.
func (c *WithoutFieldsClient) Use(hooks ...Hook) {
	c.hooks.WithoutFields = append(c.hooks.WithoutFields, hooks...)
}

// Create returns a create builder for WithoutFields.
func (c *WithoutFieldsClient) Create() *WithoutFieldsCreate {
	mutation := newWithoutFieldsMutation(c.config, OpCreate)
	return &WithoutFieldsCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of WithoutFields entities.
func (c *WithoutFieldsClient) CreateBulk(builders ...*WithoutFieldsCreate) *WithoutFieldsCreateBulk {
	return &WithoutFieldsCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for WithoutFields.
func (c *WithoutFieldsClient) Update() *WithoutFieldsUpdate {
	mutation := newWithoutFieldsMutation(c.config, OpUpdate)
	return &WithoutFieldsUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WithoutFieldsClient) UpdateOne(wf *WithoutFields) *WithoutFieldsUpdateOne {
	mutation := newWithoutFieldsMutation(c.config, OpUpdateOne, withWithoutFields(wf))
	return &WithoutFieldsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *WithoutFieldsClient) UpdateOneID(id int) *WithoutFieldsUpdateOne {
	mutation := newWithoutFieldsMutation(c.config, OpUpdateOne, withWithoutFieldsID(id))
	return &WithoutFieldsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WithoutFields.
func (c *WithoutFieldsClient) Delete() *WithoutFieldsDelete {
	mutation := newWithoutFieldsMutation(c.config, OpDelete)
	return &WithoutFieldsDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *WithoutFieldsClient) DeleteOne(wf *WithoutFields) *WithoutFieldsDeleteOne {
	return c.DeleteOneID(wf.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *WithoutFieldsClient) DeleteOneID(id int) *WithoutFieldsDeleteOne {
	builder := c.Delete().Where(withoutfields.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WithoutFieldsDeleteOne{builder}
}

// Query returns a query builder for WithoutFields.
func (c *WithoutFieldsClient) Query() *WithoutFieldsQuery {
	return &WithoutFieldsQuery{
		config: c.config,
	}
}

// Get returns a WithoutFields entity by its id.
func (c *WithoutFieldsClient) Get(ctx context.Context, id int) (*WithoutFields, error) {
	return c.Query().Where(withoutfields.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WithoutFieldsClient) GetX(ctx context.Context, id int) *WithoutFields {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *WithoutFieldsClient) Hooks() []Hook {
	return c.hooks.WithoutFields
}
