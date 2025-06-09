package dto

import (
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
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

func (p *PodMetricsResponse) Build(cpuMillicores float64, metrics *entity.PodMetrics) {
	deploymentName := metrics.DeploymentName
	if deploymentName.Valid {
		p.DeploymentName = &deploymentName.String
	} else {
		p.DeploymentName = nil
	}

	p.Timestamp = metrics.Timestamp
	p.PodName = metrics.PodName
	p.NamespaceName = metrics.NamespaceName
	p.NodeName = metrics.NodeName
	p.UID = metrics.UID
	p.CpuMillicores = cpuMillicores
	p.MemoryBytes = metrics.MemoryUsage
	p.DiskReadBytes = metrics.DiskReadBytes
	p.DiskWriteBytes = metrics.DiskWriteBytes
	p.NetworkRxBytes = metrics.NetworkRxBytes
	p.NetworkTxBytes = metrics.NetworkTxBytes
}
