package main

import (
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
)

type templateHandler struct {
	query    string
	template string
	client   prometheus.API
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, warn, err := th.client.Query(r.Context(), th.query, time.Now())
	if warn != nil {
		log.Printf("Warnings emitted from query %s: %v", th.query, warn)
	}
	if err != nil {
		log.Fatalf("Error emitted from query %s: %v", th.query, err)
	}
	t, err := template.New(th.query).Parse(th.template)
	if err != nil {
		log.Fatalf("Failed to parse template: %s: %v", th.template, err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatalf("Failed to execute template for data %s: %v", data, err)
	}
	w.Write([]byte("The time is: " + th.query))
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
		log.Printf("Serving query '%s' at route '/%s'", cfg.Query, path)
		mux.Handle("/"+path, &templateHandler{query: cfg.Query, template: cfg.Template, client: apiClient})
	}
	return mux
}
