package dto

import (
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
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

func (n *NodeMetricsResponse) Build(cpuMillicores float64, memoryBytes int64, metrics *entity.NodeMetrics) {
	n.Timestamp = metrics.Timestamp
	n.NodeName = metrics.NodeName
	n.CpuMillicores = cpuMillicores
	n.MemoryBytes = memoryBytes
	n.DiskReadBytes = metrics.DiskReadBytes
	n.DiskWriteBytes = metrics.DiskWriteBytes
	n.NetworkRxBytes = metrics.NetworkRxBytes
	n.NetworkTxBytes = metrics.NetworkTxBytes
}
