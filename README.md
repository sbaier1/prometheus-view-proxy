# Prometheus view proxy

A simple proxy that queries an upstream Prometheus with the given instant queries and passing the result vectors to a template.

The queries and templates are exposed under various configurable HTTP routes simply as a no-parameters GET request.

The goal of this application is to provide custom, limited insights into Prometheus metrics for building for example customer-facing metric views.

The simple routing using the HTTP GET path allows simply hooking this up using an ingress controller on K8s.

See [sample.yaml](sample.yaml) for an example configuration file.

## Templating

Currently, only the gotemplate engine is supported, but support for jsonnet is planned.

## TODO

* Proper logging framework