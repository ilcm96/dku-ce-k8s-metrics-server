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

// PodTimeSeriesResponse 는 Pod 시계열 조회 API의 응답 구조체입니다.
// 지정된 시간 구간 동안의 요약된 메트릭을 제공합니다.
type PodTimeSeriesResponse struct {
	PodName          string    `json:"pod_name"`
	DeploymentName   *string   `json:"deployment_name,omitempty"`
	NamespaceName    string    `json:"namespace_name"`
	NodeName         string    `json:"node_name"`
	UID              string    `json:"uid"`
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
