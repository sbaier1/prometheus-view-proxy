package main

import (
	"flag"
	"github.com/sbaier1/prometheus-view-proxy/viewproxy"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v2"
)

func main() {
	var (
		insecureListenAddress string
		upstream              string
		config                string
	)

	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&insecureListenAddress, "insecure-listen-address", "", "The address the prometheus-view-proxy HTTP server should listen on.")
	flagset.StringVar(&upstream, "upstream", "", "The upstream URL to proxy to.")
	flagset.StringVar(&config, "config", "", "The config to load for queries to perform.")
	//nolint: errcheck // Parse() will exit on error.
	flagset.Parse(os.Args[1:])
	if config == "" {
		log.Fatalf("-config flag cannot be empty")
	}

	upstreamURL, err := url.Parse(upstream)
	if err != nil {
		log.Fatalf("Failed to build parse upstream URL: %v", err)
	}

	if upstreamURL.Scheme != "http" && upstreamURL.Scheme != "https" {
		log.Fatalf("Invalid scheme for upstream URL %q, only 'http' and 'https' are supported", upstream)
	}

	f, err := os.Open(config)
	if err != nil {
		log.Fatalf("Could not read config at %s: %v", config, err)
	}
	defer f.Close()

	var cfg viewproxy.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Could not decode config at %s: %v", config, err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", viewproxy.NewRoutes(upstreamURL, cfg))

	srv := &http.Server{Handler: mux}

	l, err := net.Listen("tcp", insecureListenAddress)
	if err != nil {
		log.Fatalf("Failed to listen on insecure address: %v", err)
	}

	errCh := make(chan error)
	go func() {
		log.Printf("Listening insecurely on %v", l.Addr())
		errCh <- srv.Serve(l)
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		log.Print("Received SIGTERM, exiting gracefully...")
		srv.Close()
	case err := <-errCh:
		if err != http.ErrServerClosed {
			log.Printf("Server stopped with %v", err)
		}
		os.Exit(1)
	}
}
