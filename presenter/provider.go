package main

import (
	"fmt"
	"os"

	//"github.com/prometheus/client_golang/api/prometheus"
	pmodel "github.com/prometheus/common/model"
	"github.com/trustedanalytics/metrics/presenter/model"
	"log"
	"math/rand"
	"net/http"
)

type MetricsProvider interface {
	PlatformMetrics() (*model.PlatformMetrics, error)
	OrganizationMetrics(organization string) (*model.OrganizationMetrics, error)
	SingleMetric(measurement string, fields []string, from, to string) (*model.RawMetrics, error)
	RawQuery(query string) (*model.RawMetrics, error)
	Health() (string, error)
}

type PrometheusMetricsProvider struct {
	//queryAPI prometheus.QueryAPI
	client *http.Client
	url string
}

func (pmp *PrometheusMetricsProvider) execQuery(query string) (pmodel.Value, error) {
	req, err := http.NewRequest("GET", pmp.url + "/api/v1/query", nil)
	if err != nil {
		log.Panic("Error creating request to Prometheus:", err)
		return nil, err
	}
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()
	resp, err := pmp.client.Do(req)
	if err != nil {
		log.Println("Error while making request to Prometheus:", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Println(resp.Body)
	return nil, nil
}

func prometheusValueToRawMetrics(value pmodel.Value) (*model.RawMetrics, error) {
	switch value.Type() {
	case pmodel.ValScalar:
		return nil, nil
	case pmodel.ValVector:
		return nil, nil
	case pmodel.ValMatrix:
		return nil, nil
	case pmodel.ValString:
		return nil, nil
	}
	return nil, nil
}

func (pmp *PrometheusMetricsProvider) RawQuery(query string) (*model.RawMetrics, error) {
	results, err := pmp.execQuery(query)
	if err != nil {
		return nil, err
	}
	return prometheusValueToRawMetrics(results)
}

func (pmp *PrometheusMetricsProvider) PlatformMetrics() (*model.PlatformMetrics, error) {
	//selects := strings.Join([]string{
	//	tap_catalog.OrganizationsCount,
	//	tap_catalog.UsersCount,
	//	tap_catalog.ApplicationsCount,
	//	tap_catalog.ServicesInstancesCount,
	//}, ", ")
	//platformMetricsQuery := fmt.Sprintf("SELECT %s FROM %s LIMIT 1", selects, tap_catalog.PlatformMeasurement)
	//resp, err := imp.execQuery(platformMetricsQuery) // TODO query object could be reused
	//if err != nil {
	//	log.Println("Error retrieving platform metrics: ", resp.Err)
	//	return nil, err
	//}
	return resultsToPlatformMetrics(nil)
}

func resultsToPlatformMetrics(results map[string]interface{}) (*model.PlatformMetrics, error) {
	// TODO
	return &model.PlatformMetrics{
		// mocks
		OrganizationsCount: 1,
		// TODO I think it is missing users
		ApplicationsCount:     2,
		ServiceInstancesCount: 3,
		MemoryUsage:           4 * 1024 * 1024,
		LatestEvents:          5,
		Nodes: []model.NodeMetrics{
			model.NodeMetrics{},
		},
		Components: []model.ComponentMetrics{
			model.ComponentMetrics{},
		},
	}, nil
}

func (pmp *PrometheusMetricsProvider) OrganizationMetrics(organization string) (*model.OrganizationMetrics, error) {
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

func (pmp *PrometheusMetricsProvider) SingleMetric(measurement string, fields []string, from, to string) (*model.RawMetrics, error) {
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
func (pmp *PrometheusMetricsProvider) Health() (string, error) {
	resp, err := pmp.execQuery("up{kubernetes_name=\"metrics-prometheus\"}")
	if err != nil {
		return "Error running simple query to Prometheus", err
	}
	return resp.String(), nil
}

func getPrometheusMetricsProvider() (*PrometheusMetricsProvider, error) {
	addr := fmt.Sprintf("http://%s:%s",
		os.Getenv("METRICS_PROMETHEUS_SERVICE_HOST"),
		os.Getenv("METRICS_PROMETHEUS_SERVICE_PORT"),
	)
	//config := prometheus.Config{Address: addr}
	//client, err := prometheus.New(config)
	//if err != nil {
	//	return nil, err
	//}
	//queryApi := prometheus.NewQueryAPI(client)
	//return &PrometheusMetricsProvider{queryApi}, nil
	return &PrometheusMetricsProvider{http.DefaultClient, addr}, nil
}

func GetMetricsProvider() (MetricsProvider, error) {
	return getPrometheusMetricsProvider()
}
