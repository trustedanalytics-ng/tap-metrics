# Platform Metrics

System for handling metrics in TAP

## Architecture

Yellow components are the ones strictly for metrics sub-system

![](docs/arch.png)

### Core components

[Prometheus](http://prometheus.io/)

[Grafana](http://grafana.org/)


## Dashboards

WIP: now only definitions

### TAP dashboard (aka Platform Dashboard)

### Organization dashboard

### User applications dashboard

TODO - we need a way to easily identify users apps



## You want your app to have metrics collected?

In order to make your metrics available in Grafana you need to:

* collect them into Prometheus. You can:
    * instrument your code end expose collected data through endpoint. [Here you can related libraries](https://prometheus.io/docs/instrumenting/clientlibs/)
    * create exposer service that can collect metrics from services which code we can't instrument. [More info here](https://prometheus.io/docs/instrumenting/exporters/)
    * (not recommended in most cases) push your metrics to Prometheus directly. [See here for more info](https://prometheus.io/docs/instrumenting/pushing/)
* inform Prometheus that your service can be scrapped. [See here](https://github.com/prometheus/prometheus/blob/50e044bb006f74b14cc44fc65a1f3bdad0ed5676/documentation/examples/prometheus-kubernetes.yml#L73)

## TODOs

* secure connection between Grafana and Prometheus
* OAuth in Grafana (in 4.x.x - currently in beta). We need to integrate UAA
* HA deployment for Prometheus (2x instances, petsets + PV) and Grafana (2x instances) -- in deployment related repo.



