package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	sqc "entgo.io/contrib/sqlcommenter"
	"entgo.io/contrib/sqlcommenter/examples/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"

	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	_ "github.com/mattn/go-sqlite3"
)

const (
	metricsLogFile = "/tmp/metrics.log"
	tracesLogFile  = "/tmp/trace.log"
)

func initTracer() func() {
	// Using log exporter to export metrics but you can choose any supported exporter.
	exporter, err := exporter.NewLogExporter(exporter.Options{
		ReportingInterval: 10 * time.Second,
		MetricsLogFile:    metricsLogFile,
		TracesLogFile:     tracesLogFile,
	})
	if err != nil {
		log.Fatalf("Error creating log exporter: %v", err)
	}
	exporter.Start()
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Report stats at every second.
	view.SetReportingPeriod(1 * time.Second)
	return func() {
		exporter.Stop()
		exporter.Close()
	}
}

func main() {
	closeTracer := initTracer()
	defer closeTracer()
	// create db driver
	db, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// create sqlcommenter driver which wraps debug driver which wraps sqlite driver
	// we should have sqlcommenter and debug logs on every query to our sqlite DB
	commentedDriver := sqc.NewDriver(dialect.Debug(db),
		// add OpenCensus tracing tags
		sqc.WithTagger(sqc.NewOCTagger()),
		sqc.WithDriverVersion(),
		sqc.WithTags(sqc.Tags{
			sqc.KeyAppliaction: "users",
			sqc.KeyFramework:   "net/http",
		}),
	)
	// create and configure ent client
	client := ent.NewClient(ent.Driver(commentedDriver))
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	client.User.Create().SetName("hedwigz").SaveX(context.Background())
	getUsersHandler := func(rw http.ResponseWriter, r *http.Request) {
		users := client.User.Query().AllX(r.Context())
		b, _ := json.Marshal(users)
		rw.WriteHeader(http.StatusOK)
		rw.Write(b)
	}

	backend := &ochttp.Handler{
		Handler: http.HandlerFunc(getUsersHandler),
	}
	testRequest(backend)
}

func testRequest(handler http.Handler) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// debug printer should print sql statement with comment
	handler.ServeHTTP(w, req)
}
