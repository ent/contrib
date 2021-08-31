package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	sqc "entgo.io/contrib/sqlcommenter"
	"entgo.io/contrib/sqlcommenter/examples/oc/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-chi/chi"

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
	// create and configure ent client
	db, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	commentedDriver := sqc.NewDriver(dialect.Debug(db), sqc.WithTagger(
		sqc.NewOCTrace(),
		sqc.NewStaticTagger(sqc.SQLCommentTags{
			sqc.ApplicationTagKey: "bootcamp",
			sqc.FrameworkTagKey:   "go-chi",
		}),
		sqc.NewDriverVersionCommenter(),
	))
	client := ent.NewClient(ent.Driver(commentedDriver))
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	client.User.Create().SetName("hedwigz").SaveX(context.Background())
	r := chi.NewRouter()
	r.Get("/users", func(rw http.ResponseWriter, r *http.Request) {
		users := client.User.Query().AllX(r.Context())
		b, _ := json.Marshal(users)
		rw.WriteHeader(http.StatusOK)
		rw.Write(b)
	})

	backend := &ochttp.Handler{
		Handler: r,
	}
	testRequest(backend)
}

func testRequest(handler http.Handler) {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	// debug printer should print sql statement with comments
	handler.ServeHTTP(w, req)
}