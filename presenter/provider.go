package main

import (
	"fmt"
	"os"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/trustedanalytics/metrics/presenter/model"
)

type MetricsProvider interface {
	PlatformMetrics() (*model.PlatformMetrics, error)
	OrganizationMetrics(organization string) (*model.OrganizationMetrics, error)
	SingleMetric(measurement string, fields []string, from, to string) (*model.RawMetrics, error)
	RawQuery(query string) (*model.RawMetrics, error)
}

type InfluxMetricsProvider struct {
	client influx.Client
	dbName string
}

func executeQuery(measurement string, fields []string, from, to, where, groupby string) {
}

func (imp *InfluxMetricsProvider) RawQuery(query string) (*influx.Response, error) {
	q := influx.Query{
		Command:  query,
		Database: imp.dbName,
	}
	return imp.client.Query(q)
}

func (imp *InfluxMetricsProvider) PlatformMetrics() (*model.PlatformMetrics, error) {
	return nil, nil
}

func (imp *InfluxMetricsProvider) OrganizationMetrics(organization string) (*model.OrganizationMetrics, error) {
	return nil, nil
}

func influxClientConfig() *influx.HTTPConfig {
	addr := fmt.Sprintf("http://%s:%s",
		os.Getenv("METRICS_INFLUXDB_SERVICE_HOST"),
		os.Getenv("METRICS_INFLUXDB_SERVICE_PORT_API"),
	)
	return &influx.HTTPConfig{
		Addr: addr,
	}
}

func NewInfluxDBMetricsProvider() (*InfluxMetricsProvider, error) {
	client, err := influx.NewHTTPClient(*influxClientConfig())
	if err != nil {
		return nil, err
	}
	dbName := os.Getenv("METRICS_INFLUXDB_DB_NAME")
	if dbName == "" {
		dbName = "TAPMetrics"
	}
	return &InfluxMetricsProvider{client, dbName}, nil
}
