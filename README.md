# Prometheus view proxy

[![Docker Hub](https://img.shields.io/docker/pulls/sbaier1/prometheus-view-proxy.svg)](https://hub.docker.com/r/sbaier1/prometheus-view-proxy)

A simple proxy that queries an upstream Prometheus with the given instant queries and passing the result vectors to a template.

The queries and templates are exposed under various configurable HTTP routes simply as a no-parameters GET request.

The goal of this application is to provide custom, limited insights into Prometheus metrics for building for example customer-facing metric views.

The simple routing using the HTTP GET path allows simply hooking this up using an ingress controller on K8s. [Gorilla MUX](https://github.com/gorilla/mux) is used as the HTTP router, thus allowing for request parameters (and corresponding templating) as well.

See [sample.yaml](sample.yaml) for an example configuration file.

## Templating

Currently, only the gotemplate engine is supported, but support for jsonnet is planned.

## gotemplate functions

* The [sprig functions](http://masterminds.github.io/sprig/) are available in templates
* There is also a custom `getLabel(string, Metric) string` function for getting the value of a label from a Prometheus metric as a string.

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

* Proper leveled logging framework. Fatal logs will lead to exit/restarts at the moment.