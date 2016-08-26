## API Endpoints

By default works on 8081

### /api/v1/metrics/single?metric=NAME_OF_METRIC&from=now()-5m&to=now()

[response RawMetrics](https://github.com/intel-data/tapng-metrics/blob/presenter/presenter/model/model.go#L58)

### /api/v1/metrics/organization?org=NAME_OF_ORG

[response OrganizationMetrics](https://github.com/intel-data/tapng-metrics/blob/presenter/presenter/model/model.go#L3)

### /api/v1/metrics/platform

[response PlatformMetric](https://github.com/intel-data/tapng-metrics/blob/presenter/presenter/model/model.go#L20)



