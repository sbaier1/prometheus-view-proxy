# Prometheus view proxy

[![Docker Hub](https://img.shields.io/docker/pulls/sbaier1/prometheus-view-proxy.svg)](https://hub.docker.com/r/sbaier1/prometheus-view-proxy)

A simple proxy that queries an upstream Prometheus with the given instant queries and passing the result vectors to a template.

The queries and templates are exposed under various configurable HTTP routes simply as a no-parameters GET request.

The goal of this application is to provide custom, limited insights into Prometheus metrics for building for example customer-facing metric views.

The simple routing using the HTTP GET path allows simply hooking this up using an ingress controller on K8s.

See [sample.yaml](sample.yaml) for an example configuration file.

## Templating

Currently, only the gotemplate engine is supported, but support for jsonnet is planned.

## Deploying

### Helm

```sh
helm repo add prom-view-proxy https://sbaier1.github.io/prometheus-view-proxy
# Set at the very least the upstream prometheus URL. You should create the config for your requirements and supply them directly instead if possible though.
helm upgrade --install prometheus-view-proxy prom-view-proxy/prometheus-view-proxy --set config.prometheus.url=http://prometheus.namespace.svc.cluster.local:9090
```

## Metrics

| Name                                 | Description                                                               |
| ------------------------------------ | ------------------------------------------------------------------------- |
| viewproxy_backend_queries_count      | The total number of queries made to the backend                           |
| viewproxy_backend_warnings_count     | The total number of warnings received from the backend                    |
| viewproxy_backend_errors_count       | The total number of errors received from the backend                      |
| viewproxy_backend_invalid_type_count | The total number of responses with invalid type received from the backend |

## TODO

* Proper logging framework
* Consider switching to gorilla mux so regex can be used in path specs (and maybe can be captured and used in the routes/queries?)