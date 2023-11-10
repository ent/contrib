package todo

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"entgo.io/contrib/entgql/internal/todo/ent"
)

func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return r.client.Noder(ctx, id)
}

func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return r.client.Noders(ctx, ids)
}

func (r *queryResolver) BillProducts(ctx context.Context) ([]*ent.BillProduct, error) {
	return r.client.BillProduct.Query().All(ctx)
}

func (r *queryResolver) Categories(ctx context.Context, limit *int, offset *int, orderBy []*ent.CategoryOrder, where *ent.CategoryWhereInput) (*ent.CategoryList, error) {
	return r.client.Category.Query().
		PaginateLimitOffset(ctx, limit, offset,
			ent.WithCategoryOrder(orderBy),
			ent.WithCategoryFilter(where.Filter),
		)
}

func (r *queryResolver) Groups(ctx context.Context, limit *int, offset *int, where *ent.GroupWhereInput) (*ent.GroupList, error) {
	return r.client.Group.Query().
		PaginateLimitOffset(ctx, limit, offset,
			ent.WithGroupFilter(where.Filter),
		)
}

func (r *queryResolver) OneToMany(ctx context.Context, limit *int, offset *int, orderBy *ent.OneToManyOrder, where *ent.OneToManyWhereInput) (*ent.OneToManyList, error) {
	return r.client.OneToMany.Query().
		PaginateLimitOffset(ctx, limit, offset,
			ent.WithOneToManyOrder(orderBy),
			ent.WithOneToManyFilter(where.Filter),
		)
}

func (r *queryResolver) Todos(ctx context.Context, limit *int, offset *int, orderBy *ent.TodoOrder, where *ent.TodoWhereInput) (*ent.TodoList, error) {
	return r.client.Todo.Query().
		PaginateLimitOffset(ctx, limit, offset,
			ent.WithTodoOrder(orderBy),
			ent.WithTodoFilter(where.Filter),
		)
}

func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int, orderBy *ent.UserOrder, where *ent.UserWhereInput) (*ent.UserList, error) {
	return r.client.User.Query().
		PaginateLimitOffset(ctx, limit, offset,
			ent.WithUserOrder(orderBy),
			ent.WithUserFilter(where.Filter),
		)
}

// Category returns CategoryResolver implementation.
func (r *Resolver) Category() CategoryResolver { return &categoryResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Todo returns TodoResolver implementation.
func (r *Resolver) Todo() TodoResolver { return &todoResolver{r} }

// CreateCategoryInput returns CreateCategoryInputResolver implementation.
func (r *Resolver) CreateCategoryInput() CreateCategoryInputResolver {
	return &createCategoryInputResolver{r}
}

// TodoWhereInput returns TodoWhereInputResolver implementation.
func (r *Resolver) TodoWhereInput() TodoWhereInputResolver { return &todoWhereInputResolver{r} }

type categoryResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
type createCategoryInputResolver struct{ *Resolver }
type todoWhereInputResolver struct{ *Resolver }
