package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	sqc "entgo.io/contrib/sqlcommenter"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"

	"entgo.io/contrib/sqlcommenter/examples/ent"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type (
	routeKey          struct{}
	MyCustomCommenter struct{}
)

func (mcc MyCustomCommenter) Tag(ctx context.Context) sqc.Tags {
	return sqc.Tags{
		"key": "value",
	}
}

func initTracer() *sdktrace.TracerProvider {
	exporter, err := stdout.New(stdout.WithWriter(io.Discard))
	if err != nil {
		log.Fatal(err)
	}
	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("ExampleService"))),
	)
	otel.SetTracerProvider(tp)
	// Add propagation.TaceContext{} which will be used by OtelTagger to inject trace information.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// create db driver
	db, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	commentedDriver := sqc.NewDriver(dialect.Debug(db),
		sqc.WithTagger(
			// add tracing info with Open Telemetry.
			sqc.NewOtelTagger(),
			// use your custom commenter
			MyCustomCommenter{},
			// map routeKey{} from context to tag named "route"
			sqc.NewContextMapper(sqc.KeyRoute, routeKey{}),
		),
		// add `db_driver` version tag
		sqc.WithDriverVersion(),
		// add some global tags to all queries
		sqc.WithTags(sqc.Tags{
			sqc.KeyAppliaction: "bootcamp",
			sqc.KeyFramework:   "go-chi",
		}))
	// create and configure ent client
	client := ent.NewClient(ent.Driver(commentedDriver))
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	client.User.Create().SetName("hedwigz").SaveX(context.Background())

	// this http middleware adds the url path to the request context, to later be used by sqc.ContextMapper.
	middleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c := context.WithValue(r.Context(), routeKey{}, r.URL.Path)
			next.ServeHTTP(w, r.WithContext(c))
		}
		return http.HandlerFunc(fn)
	}
	// some app http handler
	getUsersHandler := func(rw http.ResponseWriter, r *http.Request) {
		users := client.User.Query().AllX(r.Context())
		b, _ := json.Marshal(users)
		rw.WriteHeader(http.StatusOK)
		rw.Write(b)
	}

	backend := otelhttp.NewHandler(middleware(http.HandlerFunc(getUsersHandler)), "app")
	testRequest(backend)
}

func testRequest(handler http.Handler) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// debug printer should print sql statement with comment
	handler.ServeHTTP(w, req)
}
