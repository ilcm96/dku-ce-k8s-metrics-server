package dto

import "time"

type NamespaceMetrics struct {
	NamespaceName  string    `db:"namespace_name" json:"namespace_name"`
	Timestamp      time.Time `db:"timestamp" json:"timestamp"`
	CpuMillicores  float64   `db:"cpu_millicores" json:"cpu_millicores"`
	MemoryBytes    int64     `db:"memory_bytes" json:"memory_bytes"`
	DiskReadBytes  int64     `db:"disk_read_bytes" json:"disk_read_bytes"`
	DiskWriteBytes int64     `db:"disk_write_bytes" json:"disk_write_bytes"`
	NetworkRxBytes int64     `db:"network_rx_bytes" json:"network_rx_bytes"`
	NetworkTxBytes int64     `db:"network_tx_bytes" json:"network_tx_bytes"`
	PodCount       int       `db:"pod_count" json:"pod_count"`
}

// NamespaceMetricsResponse는 NamespaceMetrics의 별칭입니다.
// API 응답용으로 사용되지만 실제로는 동일한 구조체입니다.
type NamespaceMetricsResponse = NamespaceMetrics
