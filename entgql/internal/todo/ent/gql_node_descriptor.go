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
//
// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"

	"entgo.io/contrib/entgql/internal/todo/ent/category"
	"entgo.io/contrib/entgql/internal/todo/ent/friendship"
	"entgo.io/contrib/entgql/internal/todo/ent/group"
	"entgo.io/contrib/entgql/internal/todo/ent/todo"
	"entgo.io/contrib/entgql/internal/todo/ent/user"
)

// Node in the graph.
type Node struct {
	ID     int      `json:"id,omitempty"`     // node id.
	Type   string   `json:"type,omitempty"`   // node type.
	Fields []*Field `json:"fields,omitempty"` // node fields.
	Edges  []*Edge  `json:"edges,omitempty"`  // node edges.
}

// Field of a node.
type Field struct {
	Type  string `json:"type,omitempty"`  // field type.
	Name  string `json:"name,omitempty"`  // field name (as in struct).
	Value string `json:"value,omitempty"` // stringified value.
}

// Edges between two nodes.
type Edge struct {
	Type string `json:"type,omitempty"` // edge type.
	Name string `json:"name,omitempty"` // edge name.
	IDs  []int  `json:"ids,omitempty"`  // node ids (where this edge point to).
}

// Node implements Noder interface
func (bp *BillProduct) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     bp.ID,
		Type:   "BillProduct",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(bp.Name); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(bp.Sku); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "sku",
		Value: string(buf),
	}
	if buf, err = json.Marshal(bp.Quantity); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "uint64",
		Name:  "quantity",
		Value: string(buf),
	}
	return node, nil
}

// Node implements Noder interface
func (c *Category) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     c.ID,
		Type:   "Category",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(c.Text); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "text",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Status); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "category.Status",
		Name:  "status",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Config); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "*schematype.CategoryConfig",
		Name:  "config",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Duration); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "time.Duration",
		Name:  "duration",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Count); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "uint64",
		Name:  "count",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Strings); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "[]string",
		Name:  "strings",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "Todo",
		Name: "todos",
	}
	err = c.QueryTodos().
		Select(todo.FieldID).
		Scan(ctx, &node.Edges[0].IDs)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Node implements Noder interface
func (f *Friendship) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     f.ID,
		Type:   "Friendship",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(f.CreatedAt); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "created_at",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.UserID); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "int",
		Name:  "user_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.FriendID); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "friend_id",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "User",
		Name: "user",
	}
	err = f.QueryUser().
		Select(user.FieldID).
		Scan(ctx, &node.Edges[0].IDs)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "User",
		Name: "friend",
	}
	err = f.QueryFriend().
		Select(user.FieldID).
		Scan(ctx, &node.Edges[1].IDs)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Node implements Noder interface
func (gr *Group) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     gr.ID,
		Type:   "Group",
		Fields: make([]*Field, 1),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(gr.Name); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "User",
		Name: "users",
	}
	err = gr.QueryUsers().
		Select(user.FieldID).
		Scan(ctx, &node.Edges[0].IDs)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Node implements Noder interface
func (t *Todo) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     t.ID,
		Type:   "Todo",
		Fields: make([]*Field, 8),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(t.CreatedAt); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "created_at",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.Status); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "todo.Status",
		Name:  "status",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.Priority); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "priority",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.Text); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "text",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.CategoryID); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "category_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.Init); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "map[string]interface {}",
		Name:  "init",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.Custom); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "[]customstruct.Custom",
		Name:  "custom",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.Customp); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "[]*customstruct.Custom",
		Name:  "customp",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "Todo",
		Name: "parent",
	}
	err = t.QueryParent().
		Select(todo.FieldID).
		Scan(ctx, &node.Edges[0].IDs)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "Todo",
		Name: "children",
	}
	err = t.QueryChildren().
		Select(todo.FieldID).
		Scan(ctx, &node.Edges[1].IDs)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		Type: "Category",
		Name: "category",
	}
	err = t.QueryCategory().
		Select(category.FieldID).
		Scan(ctx, &node.Edges[2].IDs)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Node implements Noder interface
func (u *User) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     u.ID,
		Type:   "User",
		Fields: make([]*Field, 2),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(u.Name); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.Password); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "password",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "Group",
		Name: "groups",
	}
	err = u.QueryGroups().
		Select(group.FieldID).
		Scan(ctx, &node.Edges[0].IDs)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "User",
		Name: "friends",
	}
	err = u.QueryFriends().
		Select(user.FieldID).
		Scan(ctx, &node.Edges[1].IDs)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		Type: "Friendship",
		Name: "friendships",
	}
	err = u.QueryFriendships().
		Select(friendship.FieldID).
		Scan(ctx, &node.Edges[2].IDs)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Node returns the node with given global ID.
//
// This API helpful in case you want to build
// an administrator tool to browser all types in system.
func (c *Client) Node(ctx context.Context, id int) (*Node, error) {
	n, err := c.Noder(ctx, id)
	if err != nil {
		return nil, err
	}
	return n.Node(ctx)
}
