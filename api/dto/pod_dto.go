package dto

import (
	"time"
)

type PodMetricsResponse struct {
	Timestamp      time.Time `json:"timestamp"`
	PodName        string    `json:"pod_name"`
	DeploymentName *string   `json:"deployment_name,omitempty"`
	NamespaceName  string    `json:"namespace_name"`
	NodeName       string    `json:"node_name"`
	UID            string    `json:"uid"`
	CpuMillicores  float64   `json:"cpu_millicores"`
	MemoryBytes    int64     `json:"memory_bytes"`
	DiskReadBytes  int64     `json:"disk_read_bytes"`
	DiskWriteBytes int64     `json:"disk_write_bytes"`
	NetworkRxBytes int64     `json:"network_rx_bytes"`
	NetworkTxBytes int64     `json:"network_tx_bytes"`
}
