package main

import (
	utils "github.com/trustedanalytics/metrics/collectors"
	catalogMetrics "github.com/trustedanalytics/metrics/collectors/tap_catalog/metrics"
	tapCatalogClient "github.com/trustedanalytics/tap-catalog/client"

	tapImageFactory "github.com/trustedanalytics/tap-image-factory/app"
	// Owner of Catalog could provide proper method for initialization
	// that is located there

	"log"
	"os"
	"time"
)

const (
	collectorAddress = "localhost:8094"
	metricsBuffer    = 1000
	defaultSchedule  = "*/10 * * * *"
)

type TAPCatalogMetricsCollector struct {
	client tapCatalogClient.TapCatalogApi
}

func (catalogMC *TAPCatalogMetricsCollector) CollectMetrics() ([]*utils.Metrics, error) {
	timestamp := time.Now().Unix()
	metrics := []*utils.Metrics{}
	var err error
	// We are interested in first error as others might come because of it
	for _, provider := range catalogMetrics.TapCatalogMetricsProviders {
		tmpMetric, tmpErr := provider(catalogMC.client)
		if tmpErr != nil {
			log.Println("Error while loading metric: ", tmpErr)
			if err == nil {
				err = tmpErr
			}
		}
		metrics = append(metrics, tmpMetric...)
	}
	for i := range metrics {
		metrics[i].Timestamp = timestamp
	}
	return metrics, err
}

func getTAPCatalogMetricsCollector() (*TAPCatalogMetricsCollector, error) {
	client, err := tapImageFactory.GetCatalogConnector()
	if err != nil {
		log.Println("Error while getting TAP Catalog Connector", err)
		return nil, err
	}
	return &TAPCatalogMetricsCollector{client}, nil
}

func main() {
	log.Println("Starting TAP Catalog Metrics Collector")
	schedule := os.Getenv("SCHEDULE")
	if schedule == "" {
		schedule = defaultSchedule
	}

	var collector utils.MetricsCollector
	collector, err := getTAPCatalogMetricsCollector()
	if err != nil {
		log.Fatal("Error creating TAP Catalog Metrics Collector: ", err)
	}

	collectingEngine, err := utils.NewCollectingEngine(collectorAddress, metricsBuffer)
	if err != nil {
		log.Fatal("Error creating CollectingEngine: ", err)
	}

	collectingEngine.AddMetricsCollector(schedule, &collector)
	log.Println("Starting collecting metrics")
	collectingEngine.StartAndBlockTillSigInt()
	log.Println("Stoped collecting metrics")
}
