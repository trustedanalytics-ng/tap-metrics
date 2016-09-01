package main

import (
	"fmt"
	"os"

	"encoding/json"
	influx "github.com/influxdata/influxdb/client/v2"
	tap_catalog "github.com/trustedanalytics/metrics/collectors/tap_catalog/metrics"
	"github.com/trustedanalytics/metrics/presenter/model"
	"log"
	"math/rand"
	"strings"
	"time"
)

type MetricsProvider interface {
	PlatformMetrics() (*model.PlatformMetrics, error)
	OrganizationMetrics(organization string) (*model.OrganizationMetrics, error)
	SingleMetric(measurement string, fields []string, from, to string) (*model.RawMetrics, error)
	RawQuery(query string) (*model.RawMetrics, error)
	Health() (string, error)
}

type InfluxMetricsProvider struct {
	client influx.Client
	dbName string
}

func (imp *InfluxMetricsProvider) execQuery(query string) (*influx.Response, error) {
	q := influx.Query{
		Command:  query,
		Database: imp.dbName,
	}
	return imp.client.Query(q)
}

func (imp *InfluxMetricsProvider) RawQuery(query string) (*model.RawMetrics, error) {
	// TODO
	results, err := imp.execQuery(query)
	fmt.Println(results, err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (imp *InfluxMetricsProvider) PlatformMetrics() (*model.PlatformMetrics, error) {
	selects := strings.Join([]string{
		tap_catalog.OrganizationsCount,
		tap_catalog.UsersCount,
		tap_catalog.ApplicationsCount,
		tap_catalog.ServicesInstancesCount,
	}, ", ")
	platformMetricsQuery := fmt.Sprintf("SELECT %s FROM %s LIMIT 1", selects, tap_catalog.PlatformMeasurement)
	resp, err := imp.execQuery(platformMetricsQuery) // TODO query object could be reused
	if err != nil {
		log.Println("Error retrieving platform metrics: ", resp.Err)
		return nil, err
	}
	return resultsToPlatformMetrics(resp.Results)
}

func toInt(number interface{}) int {
	n, err := number.(json.Number).Int64()
	if err != nil {
		log.Panic("Unable to cast to int value: ", number)
	}
	return int(n)
}

func resultsToPlatformMetrics(results []influx.Result) (*model.PlatformMetrics, error) {
	row := results[0].Series[0].Values[0]
	return &model.PlatformMetrics{
		OrganizationsCount: toInt(row[1]),
		// TODO I think it is missing users
		ApplicationsCount:     toInt(row[3]),
		ServiceInstancesCount: toInt(row[4]),
		// mocks
		MemoryUsage:  4 * 1024 * 1024,
		LatestEvents: 5,
		Nodes: []model.NodeMetrics{
			model.NodeMetrics{},
		},
		Components: []model.ComponentMetrics{
			model.ComponentMetrics{},
		},
	}, nil
}

func (imp *InfluxMetricsProvider) OrganizationMetrics(organization string) (*model.OrganizationMetrics, error) {
	// TODO this is mock
	return &model.OrganizationMetrics{
		OrganizationID: "defaultID",
		Name:           "default",
		Status:         "OK",
		ApplicationsRunningCount: 1,
		ApplicationsFailedCount:  2,
		ServicesCount:            3,
		ServicesUsagePercentage:  0.1,
		UsersCount:               4,
		MemoryUsage:              5 * 1024 * 1024,
		MemoryUsagePercentage:    0.2,
		CpuUsage:                 6,
		CpuUsagePercentage:       0.3,
		PublicDatasetsCount:      7,
		PrivateDatasetsCount:     8,
	}, nil
}

func (imp *InfluxMetricsProvider) SingleMetric(measurement string, fields []string, from, to string) (*model.RawMetrics, error) {
	// TODO this is mock
	metrics := []model.RawMetric{}
	for _, field := range fields {
		values := make([]model.RawMetricValue, 100)
		for i := range values {
			values[i] = model.RawMetricValue{
				Timestamp: int64(1472126069+i*12345) * 1000,
				Value:     3.14 * rand.Float32(),
			}
		}
		metrics = append(metrics, model.RawMetric{
			Name:   field,
			Values: values,
		})
	}
	return &model.RawMetrics{
		Metrics: metrics,
	}, nil
}
func (imp *InfluxMetricsProvider) Health() (string, error) {
	duration, msg, err := imp.client.Ping(time.Second)
	retMsg := fmt.Sprintf("Ping: %v\nMessage: %s\nError: %v\n", duration, msg, err)
	return retMsg, err
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
