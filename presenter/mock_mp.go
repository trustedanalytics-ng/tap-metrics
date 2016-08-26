package main

import (
	"github.com/trustedanalytics/metrics/presenter/model"
	"math/rand"
)

type MockMetricsProvider struct{}

func (*MockMetricsProvider) PlatformMetrics() (*model.PlatformMetrics, error) {
	return &model.PlatformMetrics{
		OrganizationsCount: 1,
		ApplicationsCount: 2,
		ServiceInstancesCount: 3,
		MemoryUsage: 4 * 1024 * 1024,
		LatestEvents: 5,
		Nodes: []model.NodeMetrics{
			model.NodeMetrics{},
		},
		Components: []model.ComponentMetrics{
			model.ComponentMetrics{},
		},
	}, nil
}

func (*MockMetricsProvider) OrganizationMetrics(organization string) (*model.OrganizationMetrics, error) {
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

func (*MockMetricsProvider) SingleMetric(measurement string, fields []string, from, to string) (*model.RawMetrics, error) {
	metrics := []model.RawMetric{}
	for _, field := range fields {
		values := make([]model.RawMetricValue, 100)
		for i := range values {
			values[i] = model.RawMetricValue{
				Timestamp: int64(1472126069 + i * 12345) * 1000,
				Value: 3.14 * rand.Float32(),
			}
		}
		metrics = append(metrics, model.RawMetric{
			Name: field,
			Values: values,
		})
	}
	return &model.RawMetrics{
		Metrics: metrics,
	}, nil
}

func (m *MockMetricsProvider) RawQuery(query string) (*model.RawMetrics, error) {
	return m.SingleMetric("default", []string{"aa", "bb"}, "123", "456")
}

