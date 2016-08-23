# Platform Metrics

System for handling metrics in TAP

## Architecture

![](docs/arch.png)

* Each specyfic collector runs as separate pod to reduce friction between different kinds of services and reduce impact of any failure
  * Task of adding new collector was simplified due to provided library
* Using Telegraf on each specyfic collector to ensure (in future) secured connection with InfluxDB

## TODOs

NEEDS HEAPSTER TO COLLECT DATA RELATED TO CLUSTER PERFORMANCE AND STATE

### Metrics DB

* InfluxDB with direct access - simpler that way

### Metrics read/query API

* should it be direct access to InfluxDB queries? (with auth checks on org/platform) or be simplictly wrapped?

### Collectors (Metrics pushers)

* k8s API
* Catalog
* Data Catalog
* TAP Components status endpoint


## Some further TODOs

* enable auth in InfluxDB https://docs.influxdata.com/influxdb/v0.13/administration/authentication_and_authorization/
* InfluxDB deployment in a HA mode
* Retention Policy in Influx to store data newer than X (3 days?)


