package dto

import (
	"time"
)

type NodeMetricsResponse struct {
	Timestamp      time.Time `json:"timestamp"`
	NodeName       string    `json:"node_name"`
	CpuMillicores  float64   `json:"cpu_millicores"`
	MemoryBytes    int64     `json:"memory_bytes"`
	DiskReadBytes  int64     `json:"disk_read_bytes"`
	DiskWriteBytes int64     `json:"disk_write_bytes"`
	NetworkRxBytes int64     `json:"network_rx_bytes"`
	NetworkTxBytes int64     `json:"network_tx_bytes"`
}
