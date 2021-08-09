## prometheus-hook

`entprom.Hook` Is a hook for sending `ent` metrics to your prometheus when mutating the graph.

#### How to install and run?

```go

client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
if err != nil {
    log.Fatalf("failed opening connection to sqlite: %v", err)
}
defer client.Close()
ctx := context.Background()

if err := client.Schema.Create(ctx); err != nil {
    log.Fatalf("failed creating schema resources: %v", err)
}
// Basic usage of the hook, this will apply only for mutations.
client.Use(entprom.Hook())
```

If you want to add some custom labels you can use `entprom.Labels`.

```go
client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
if err != nil {
    log.Fatalf("failed opening connection to sqlite: %v", err)
}
ctx := context.Background()

if err := client.Schema.Create(ctx); err != nil {
    log.Fatalf("failed creating schema resources: %v", err)
}
//use the hook and add custom labels of your choice.
client.Use(entprom.Hook(
    entprom.Labels(map[string]string{"environment": "dev"}),
))
return client
```
If you want to register to specific mutations.

```go
client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
if err != nil {
    log.Fatalf("failed opening connection to sqlite: %v", err)
}
ctx := context.Background()

if err := client.Schema.Create(ctx); err != nil {
    log.Fatalf("failed creating schema resources: %v", err)
}
//use the hook only on OpUpdate and OpUpdateOne.
client.Use(hook.On(entprom.Hook(), ent.OpUpdate|ent.OpUpdateOne))
return client
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


Adding prometheus handler to the http router for scraping.

```go
dbClient = initDB()
defer dbClient.Close()
http.HandleFunc("/", handler)

//get from github.com/prometheus/client_golang/prometheus/promhttp
http.Handle("/metrics", promhttp.Handler())

log.Println("server starting on port 8080")
log.Fatal(http.ListenAndServe(":8080", nil))
```