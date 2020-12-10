package viewproxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	prom_metrics "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/model"

	"github.com/goburrow/cache"
)

// Handler params
type templateHandler struct {
	queries  []Queries
	template *template.Template
	cache    cache.LoadingCache
}

// Response passed as array to the template engine
type queryResponse struct {
	Name     string
	Response model.Vector
}

type queryWithContext struct {
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var responses []queryResponse = make([]queryResponse, len(th.queries))
	// Run queries and save results
	for index, q := range th.queries {
		val, err := th.cache.Get(q.Query)
		if err != nil {
			log.Printf("Failed to load response for query %s: %v", q.Query, err)
		}
		responses[index] = queryResponse{Name: q.Name, Response: val.(model.Vector)}
	}

	// Write the template result directly to the response
	err := th.template.Execute(w, responses)
	if err != nil {
		log.Fatalf("Failed to execute template for data '%s': %v", responses, err)
	}
}

// NewRoutes generates the root HTTP handler for the proxy
func NewRoutes(upstream *url.URL, config Config) http.Handler {
	mux := http.NewServeMux()
	backendQueriesCounter := promauto.NewCounter(prom_metrics.CounterOpts{
		Name: "viewproxy_backend_queries_count",
		Help: "The total number of queries made to the backend",
	})
	backendWarningsCounter := promauto.NewCounter(prom_metrics.CounterOpts{
		Name: "viewproxy_backend_warnings_count",
		Help: "The total number of warnings received from the backend",
	})
	backendErrorsCounter := promauto.NewCounter(prom_metrics.CounterOpts{
		Name: "viewproxy_backend_errors_count",
		Help: "The total number of errors received from the backend",
	})
	invalidTypeCounter := promauto.NewCounter(prom_metrics.CounterOpts{
		Name: "viewproxy_backend_invalid_type_count",
		Help: "The total number of responses with invalid type received from the backend",
	})
	client, err := api.NewClient(api.Config{Address: upstream.String()})
	if err != nil {
		log.Fatalf("Failed to initialize Prometheus API client %v", err)
	}
	apiClient := prometheus.NewAPI(client)

	// Initialize query result cache
	load := func(k cache.Key) (cache.Value, error) {
		query := k.(string)
		backendQueriesCounter.Inc()
		data, warn, err := apiClient.Query(context.Background(), query, time.Now())
		if warn != nil {
			log.Printf("Warnings emitted from query %s: %v", query, warn)
			backendWarningsCounter.Inc()
		}
		if err != nil {
			log.Fatalf("Error emitted from query %s: %v", query, err)
			backendErrorsCounter.Inc()
		}
		valType := data.Type()
		if valType == model.ValVector {
			vector := data.(model.Vector)
			return vector, nil
		}
		log.Printf("Response to query %s was of unexpected type %s. Will not pass response to template", query, valType)
		invalidTypeCounter.Inc()
		return nil, fmt.Errorf("query %s did not return a vector result", query)
	}
	log.Printf("Setting response cache expiry duration to %s", config.ResponseExpiryTime)
	c := cache.NewLoadingCache(load,
		cache.WithMaximumSize(1024),                           // Limit number of entries in the cache.
		cache.WithExpireAfterWrite(config.ResponseExpiryTime), // Expire entries after 2 minutes since last created.
	)

	for path, cfg := range config.Routes {
		t, err := template.New(cfg.Template).Parse(cfg.Template)
		if err != nil {
			log.Fatalf("Failed to parse template: %s: %v", cfg.Template, err)
		}
		if path == "metrics" {
			log.Fatalf("Route /metrics is reserved for this application's metrics")
		}
		log.Printf("Adding route '/%s'", path)
		mux.Handle("/"+path, &templateHandler{queries: cfg.Queries, template: t, cache: c})
	}
	log.Printf("View proxy routes initialized")
	return mux
}
