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

// NamespaceMetricsResponse 는 NamespaceMetrics의 별칭입니다.
type NamespaceMetricsResponse = NamespaceMetrics

// NamespaceTimeSeriesResponse 는 Namespace 시계열 조회 API의 응답 구조체입니다.
// 지정된 시간 구간 동안의 요약된 메트릭을 제공합니다.
type NamespaceTimeSeriesResponse struct {
	NamespaceName    string    `json:"namespace_name"`
	Window           string    `json:"window"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	AvgCpuMillicores float64   `json:"avg_cpu_millicores"`
	AvgMemoryBytes   int64     `json:"avg_memory_bytes"`
	AvgDiskReadRate  float64   `json:"avg_disk_read_rate"`  // bytes/sec
	AvgDiskWriteRate float64   `json:"avg_disk_write_rate"` // bytes/sec
	AvgNetworkRxRate float64   `json:"avg_network_rx_rate"` // bytes/sec
	AvgNetworkTxRate float64   `json:"avg_network_tx_rate"` // bytes/sec
}
