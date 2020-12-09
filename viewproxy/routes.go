package viewproxy

import (
	"github.com/prometheus/common/model"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
)

// Handler params
type templateHandler struct {
	queries  []Queries
	template string
	client   prometheus.API
}

// Response passed as array to the template engine
type queryResponse struct {
	Name     string
	Response model.Vector
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var responses []queryResponse = make([]queryResponse, len(th.queries))
	// Run queries and save results
	for index, q := range th.queries {
		data, warn, err := th.client.Query(r.Context(), q.Query, time.Now())
		if warn != nil {
			log.Printf("Warnings emitted from query %s: %v", q.Query, warn)
		}
		if err != nil {
			log.Fatalf("Error emitted from query %s: %v", q.Query, err)
		}
		valType := data.Type()
		if valType == model.ValVector {
			vector := data.(model.Vector)
			responses[index] = queryResponse{Name: q.Name, Response: vector}
		} else {
			log.Printf("Response to query %s was of unexpected type %s. Will not pass response to template", q.Query, valType)
		}
	}

	t, err := template.New("test").Parse(th.template)
	if err != nil {
		log.Fatalf("Failed to parse template: %s: %v", th.template, err)
	}
	// Write the template result directly to the response
	err = t.Execute(w, responses)
	if err != nil {
		log.Fatalf("Failed to execute template for data '%s': %v", responses, err)
	}
}

// NewRoutes generates the HTTP handler for the proxy
func NewRoutes(upstream *url.URL, config Config) http.Handler {
	mux := http.NewServeMux()
	client, err := api.NewClient(api.Config{Address: upstream.String()})
	if err != nil {
		log.Fatalf("Failed to initialize Prometheus API client %v", err)
	}
	apiClient := prometheus.NewAPI(client)
	for path, cfg := range config.Routes {
		log.Printf("Serving Queries '%s' at route '/%s'", cfg.Queries, path)
		mux.Handle("/"+path, &templateHandler{queries: cfg.Queries, template: cfg.Template, client: apiClient})
	}
	return mux
}
