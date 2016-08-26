package model

type OrganizationMetrics struct {
	OrganizationID           string  `json:"organization_id"`
	Name                     string  `json:"name"`
	Status                   string  `json:"status"`
	ApplicationsRunningCount int     `json:"applications_running_count"`
	ApplicationsFailedCount  int     `json:"applications_failed_count"`
	ServicesCount            int     `json:"services_count"`
	ServicesUsagePercentage  float32 `json:"services_usage_percentage"`
	UsersCount               int     `json:"users_count"`
	MemoryUsage              int     `json:"memory_usage"`
	MemoryUsagePercentage    float32 `json:"memory_usage_percentage"`
	CpuUsage                 int     `json:"cpu_usage"`
	CpuUsagePercentage       float32 `json:"cpu_usage_percentage"`
	PublicDatasetsCount      int     `json:"public_datasets_count"`
	PrivateDatasetsCount     int     `json:"private_datasets_count"`
}

type PlatformMetrics struct {
	OrganizationsCount    int                `json:"organizations_count"`
	ApplicationsCount     int                `json:"applications_count"`
	ServiceInstancesCount int                `json:"service_instances_count"`
	MemoryUsage           int                `json:"memory_usage"`
	LatestEvents          int                `json:"latest_eventes"`
	Nodes                 []NodeMetrics      `json:"nodes"`
	Components            []ComponentMetrics `json:"components"`
}

type NodeMetrics struct {
	Name             string               `json:"name"`
	Capacity         NodeResources        `json:"capacity"`
	Allocatable      NodeResources        `json:"allocatable"`
	Status           string               `json:"status"`
	SoftwareVersions NodeSoftwareVersions `json:"software_versions"`
	ImagesCount      int                  `json:"images_count"`
}

type NodeResources struct {
	NumberOfCores int `json:"number_of_cores"`
	Memory        int `json:"memory"`
	MaxPodCount   int `json:"max_pod_count"` // ???
}

type NodeSoftwareVersions struct {
	Kernel    string `json:"kernel"`
	OsImage   string `json:"os_image"`
	Kubelet   string `json:"kubelet"`
	KubeProxy string `json:"kube_proxy"`
}

type ComponentMetrics struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RawMetrics struct {
	Metrics []RawMetric `json:"metrics"`
}

type RawMetric struct {
	Name   string           `json:"name"`
	Values []RawMetricValue `json:"values"`
}

type RawMetricValue struct {
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}
