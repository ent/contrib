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

Next, set up your HTTP listener:
```go
//get from github.com/prometheus/client_golang/prometheus/promhttp
http.Handle("/metrics", promhttp.Handler())
```

For testing you can run the example service
```shell
go run entprom/internal/main.go
"2021/08/09 14:47:46 server starting on port 8080"
```
Run mutations, second curl will fail due to constraint.
```shell
curl localhost:8080
"check the metrics%"
curl localhost:8080
"ent: constraint failed: insert node to table "users": UNIQUE constraint failed: users.name"
```
After your server is running and a few mutations were made, you can get the metrics from the `/metrics` endpoint.
```shell
# HELP ent_operation_duration_seconds Time in seconds per operation
# TYPE ent_operation_duration_seconds histogram
ent_operation_duration_seconds_bucket{environment="dev",mutation_op="OpCreate",mutation_type="File",le="0.005"} 2
ent_operation_duration_seconds_bucket{environment="dev",mutation_op="OpCreate",mutation_type="File",le="0.01"} 2
.
.
.
ent_operation_duration_seconds_bucket{environment="dev",mutation_op="OpCreate",mutation_type="File",le="+Inf"} 2
ent_operation_duration_seconds_sum{environment="dev",mutation_op="OpCreate",mutation_type="File"} 0.000125273
.
.
.
ent_operation_duration_seconds_bucket{environment="dev",mutation_op="OpCreate",mutation_type="User",le="10"} 2
ent_operation_duration_seconds_count{environment="dev",mutation_op="OpCreate",mutation_type="User"} 2
# HELP ent_operation_error Number of failed ent mutation operations
# TYPE ent_operation_error counter
ent_operation_error{environment="dev",mutation_op="OpCreate",mutation_type="User"} 1
# HELP ent_operation_total Number of ent mutation operations
# TYPE ent_operation_total counter
ent_operation_total{environment="dev",mutation_op="OpCreate",mutation_type="File"} 2
ent_operation_total{environment="dev",mutation_op="OpCreate",mutation_type="User"} 2
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
```
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
