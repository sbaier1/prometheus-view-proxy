# Prometheus View Proxy Helm Repository

## Add the Helm repository

```
helm repo add prom-view-proxy https://sbaier1.github.io/prometheus-view-proxy
```

## Install the HiveMQ operator

```
helm upgrade --install prometheus-view-proxy prom-view-proxy/prometheus-view-proxy
```