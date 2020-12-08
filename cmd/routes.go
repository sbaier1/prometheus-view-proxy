package main

import (
	"log"
	"net/http"
	"net/url"
)

type templateHandler struct {
	query    string
	template string
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + th.query))
}

// NewRoutes generates the HTTP handler for the proxy
func NewRoutes(upstream *url.URL, config Config) http.Handler {
	mux := http.NewServeMux()
	for path, cfg := range config.Routes {
		log.Printf("Serving query '%s' at route '/%s'", cfg.Query, path)
		mux.Handle("/"+path, &templateHandler{query: cfg.Query, template: cfg.Template})
	}
	return mux
}
