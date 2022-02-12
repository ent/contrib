# entproto

A ent plugin to generate protoc from ent schemas files.

### Getting Started

Install the plugin:

```shell
go get entgo.io/contrib/entproto/cmd/entproto
```
Generate protoc files with ent schemas

```shell
entproto -path=ent/schemas -idtype=int64
```
