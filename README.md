# Platform Metrics

System for handling metrics in TAP

## Architecture

![](docs/arch.png)

## TODOs

### Metrics DB

* InfluxDB with direct access - simpler that way

### Metrics sink

Accepts connections and pushes retrieved logs into Metrics DB

* Telegraf with TCP reading should be good enough but it will burden clients with writting and no auth

### Metrics read/query API

* should it be direct access to InfluxDB queries? (with auth checks on org/platform) or be simplictly wrapped?

### Collectors (Metrics pushers)

* k8s API
* Catalog
* Data Catalog
* TAP Components status endpoint


## Some further TODOs

* enable auth in InfluxDB https://docs.influxdata.com/influxdb/v0.13/administration/authentication_and_authorization/
* InfluxDB deployment in a cluster


