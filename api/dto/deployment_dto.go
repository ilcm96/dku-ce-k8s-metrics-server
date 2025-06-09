package dto

import "time"

type DeploymentMetricsResponse struct {
	DeploymentName string    `json:"deployment_name"`
	NamespaceName  string    `json:"namespace_name"`
	Timestamp      time.Time `json:"timestamp"`
	CpuMillicores  float64   `json:"cpu_millicores"`
	MemoryBytes    int64     `json:"memory_bytes"`
	DiskReadBytes  int64     `json:"disk_read_bytes"`
	DiskWriteBytes int64     `json:"disk_write_bytes"`
	NetworkRxBytes int64     `json:"network_rx_bytes"`
	NetworkTxBytes int64     `json:"network_tx_bytes"`
	PodCount       int       `json:"pod_count"`
}
