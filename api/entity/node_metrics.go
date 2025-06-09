package entity

import "time"

type NodeMetrics struct {
	ID              uint64    `db:"id"`
	Timestamp       time.Time `db:"timestamp"`
	NodeName        string    `db:"node_name"`
	CPUTotal        float64   `db:"cpu_total"`
	CPUBusy         float64   `db:"cpu_busy"`
	MemoryTotal     int64     `db:"memory_total"`
	MemoryAvailable int64     `db:"memory_available"`
	MemoryUsed      int64     `db:"memory_used"`
	DiskReadBytes   int64     `db:"disk_read_bytes"`
	DiskWriteBytes  int64     `db:"disk_write_bytes"`
	NetworkRxBytes  int64     `db:"network_rx_bytes"`
	NetworkTxBytes  int64     `db:"network_tx_bytes"`
}
