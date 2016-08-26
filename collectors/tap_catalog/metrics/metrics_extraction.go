package metrics

import (
	utils "github.com/trustedanalytics/metrics/collectors"
	//tapCatalogModels "github.com/trustedanalytics/tapng-catalog/models"
	"log"

	tapCatalogClient "github.com/trustedanalytics/tapng-catalog/client"
)

const (
	PlatformMeasurement  = "platform_measurement"
	ComponentMeasurement = "component_measurement"

	OrganizationsCount     = "organizations_count"
	UsersCount             = "users_count"
	ApplicationsCount      = "applications_count"
	ServicesInstancesCount = "services_instances_count"
	MemoryUsageBytes       = "memory_usage_bytes"
)

type TapCatalogMetricProvider func(tapCatalogClient.TapCatalogApi) ([]*utils.Metrics, error)

func PlatformMetrics(client tapCatalogClient.TapCatalogApi) ([]*utils.Metrics, error) {
	var err error
	applications, status, err := client.ListApplications()
	if err != nil {
		log.Println("Error when getting list of applications", status, err)
	}
	services, status, err := client.ListServicesInstances()
	if err != nil {
		log.Println("Error when getting list of services ", status, err)
	}
	metric := &utils.Metrics{
		Measurement: PlatformMeasurement,
		Fields: map[string]interface{}{
			OrganizationsCount:     1,  // AFAIK there are no multiple orgs
			UsersCount:             -1, // TODO it is in user management: catalog doesn't provide number of users
			ApplicationsCount:      len(applications),
			ServicesInstancesCount: len(services),
		},
		Tags: map[string]string{},
	}
	return []*utils.Metrics{metric}, err
}

func ComponentsMetrics(client tapCatalogClient.TapCatalogApi) ([]*utils.Metrics, error) {
	var err error
	applications, status, err := client.ListInstances()
	if err != nil {
		log.Println("Error when getting list of applications", status, err)
	}
	var metrics []*utils.Metrics
	for i := range applications {
		app := &applications[i]
		metric := &utils.Metrics{
			Measurement: ComponentMeasurement,
			Fields: map[string]interface{}{
				"status": string(app.State),
			},
			Tags: map[string]string{
				"id":   app.Id,
				"name": app.Name,
				"type": string(app.Type),
			},
		}
		metrics = append(metrics, metric)
	}
	return metrics, err
}

var TapCatalogMetricsProviders []TapCatalogMetricProvider = []TapCatalogMetricProvider{
	PlatformMetrics, ComponentsMetrics,
}
