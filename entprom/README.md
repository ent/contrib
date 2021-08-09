## entprom

`entprom.Hook` is a hook for exporting metrics about `ent` mutations to [Prometheus](https://prometheus.io).

### Getting Started

1. Install the latest `ent/contrib`:
   `go get -u entgo.io/contrib/`
2. Use `entprom.Hook` as a [https://entgo.io/docs/hooks/#runtime-hooks](Runtime Hook):
```go {2-5}
client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
if err != nil {
    log.Fatalf("failed opening connection to sqlite: %v", err)
}
defer client.Close()
ctx := context.Background()

if err := client.Schema.Create(ctx); err != nil {
    log.Fatalf("failed creating schema resources: %v", err)
}

// Register the hook with default options.
client.Use(entprom.Hook()) 
```

Next, wherever you are setting up your HTTP listener:
```go
//get from github.com/prometheus/client_golang/prometheus/promhttp
http.Handle("/metrics", promhttp.Handler())
```

After your server is running and a few mutations were made, you can get the metrics from the `/metrics` endpoint.

TODO(yoni): add curl example and sample output. 

### Configuration

#### Custom Labels

To set up your Prometheus counters/histograms with custom constant labels, use `entprom.Labels`:

```go
client.Use(entprom.Hook(
    entprom.Labels(map[string]string{"environment": "dev"}),
))
```

#### Register to specific Ops

Using `entprom` as shown in the example above will register metric collection for all entity types
and all operations. To collect metrics only on specific operations you can use:

```go
//use the hook only on OpUpdate and OpUpdateOne.
client.Use(hook.On(entprom.Hook(), ent.OpUpdate|ent.OpUpdateOne))
```

Attach hook only for specific schema - for example User

```go
// Hooks of the User.
func (User) Hooks() []ent.Hook {
    return []ent.Hook{
        entprom.Hook(),
    }
}
```
