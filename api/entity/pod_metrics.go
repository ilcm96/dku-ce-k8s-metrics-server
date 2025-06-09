package entity

import (
	"database/sql"
	"time"
)

type PodMetrics struct {
	ID             uint64         `db:"id"`
	Timestamp      time.Time      `db:"timestamp"`
	PodName        string         `db:"pod_name"`
	UID            string         `db:"uid"`
	CPUUsageUsec   int64          `db:"cpu_usage_usec"`
	MemoryUsage    int64          `db:"memory_usage"`
	DiskReadBytes  int64          `db:"disk_read_bytes"`
	DiskWriteBytes int64          `db:"disk_write_bytes"`
	NetworkRxBytes int64          `db:"network_rx_bytes"`
	NetworkTxBytes int64          `db:"network_tx_bytes"`
	NamespaceName  string         `db:"namespace_name"`
	DeploymentName sql.NullString `db:"deployment_name"`
	NodeName       string         `db:"node_name"`
}
