# protoc-gen-ent

A protoc plugin to generate ent schemas from .proto files.

### Getting Started

Install the plugin:

```shell
go get github.com/bionicstork/contrib/entproto/cmd/protoc-gen-ent
```

Get `entproto/cmd/protoc-gen-ent/options/opts.proto` from this repository and place it in your project. For example,

```shell
mkdir -p ent/proto/options
wget -O ent/proto/options/opts.proto https://raw.githubusercontent.com/ent/contrib/master/entproto/cmd/protoc-gen-ent/options/ent/opts.proto
```

Let's assume you have a .proto file in `ent/proto/entpb/user.proto`:

```protobuf
syntax = "proto3";

package entpb;

option go_package = "github.com/yourorg/project/ent/proto/entpb";

message User {
  string name = 1;
  string email_address = 2;
}
```

Import the opts.proto file and use the options to annotate your message:

```diff
syntax = "proto3";

package entpb;

++import "options/opts.proto";
 
option go_package = "github.com/yourorg/project/ent/proto/entpb";  

message User {
++  option (ent.schema).gen = true; // <-- tell protoc-gen-ent you want to generate a schema from this message
  string name = 1;
  string email_address = 2;
}
```

Invoke `protoc` with the ent plugin:

```shell
cd ent/
protoc -I=proto/ --ent_out=. --ent_opt=schemadir=./schema proto/entpb/user.proto
```

Observe `schema/user.go` now contains:

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{field.String("name"), field.String("email_address")}
}
func (User) Edges() []ent.Edge {
	return nil
}
```

### Options

[opts.proto](options/ent/opts.proto) contains the message proto configuration extension messages.
The [testdata/](testdata)
directory contains many usage examples. Here are the main ones you should consider:

#### Message Options

Message options configure message/schema wide behavior and are backed by the [Schema](options/ent/opts.proto#L7)
message:

* `gen` (bool) - whether to generate a schema from this message definition

```protobuf
message Basic {
  option (ent.schema).gen = true;
  string name = 1;
}
```

* `name` (string) - specify the generated ent type name.

```protobuf
message CustomName {
  option (ent.schema) = {gen: true, name: "Rotemtam"};
  string name = 1;
}
```

Will generate: 

```go
type Rotemtam struct {
	ent.Schema
}
```

#### Field Options

Field options configure field level behavior and are backed by the [Field](options/ent/opts.proto#L16) message:

For example:

```protobuf
message Pet {
  option (ent.schema).gen = true;
  string name = 1 [(ent.field) = {optional: true}];
}
```

Will generate:
```go
field.String("name").Optional()
```

#### Edge Options

To define an edge between two types we use the [Edge](options/ent/opts.proto#L28) message.

For example:

```protobuf

message Cat {
  option (ent.schema).gen = true;
  string name = 1 [(ent.field) = {optional: true, storage_key: "shem"}];
  Human owner = 2 [(ent.edge) = {}];
}
```
Will define an edge from the Cat type to the Human type named "owner":
```go
edge.To("owner", Human.Type).Unique()
```
To create the inverse (virtual) relation we add `ref: "owner"` to the Human's `cats` field definition:

```go
message Human {
  option (ent.schema).gen = true;
  string name = 1;
  repeated Cat cats = 2 [(ent.edge) = {ref: "owner"}];
}
```

Because, the field is `repeated`, a non-unique edge is created:
```go
edge.From("cats", Cat.Type).Ref("owner")
```