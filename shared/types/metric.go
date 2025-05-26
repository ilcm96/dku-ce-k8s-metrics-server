package types

import (
	"encoding/json"
	"time"
)

type Metric struct {
	Timestamp  time.Time   `json:"timestamp"`
	NodeMetric NodeMetric  `json:"nodeMetric"`
	PodMetric  []PodMetric `json:"podMetric"`
}

func (m Metric) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

type NodeMetric struct {
	NodeName        string  `json:"nodeName"`
	CPUTotal        float64 `json:"cpuTotal"`
	CPUBusy         float64 `json:"cpuBusy"`
	MemoryTotal     uint64  `json:"memoryTotal"`
	MemoryAvailable uint64  `json:"memoryAvailable,omitempty"`
	MemoryUsed      uint64  `json:"memoryUsed"`
	DiskReadBytes   uint64  `json:"diskReadBytes"`
	DiskWriteBytes  uint64  `json:"diskWriteBytes"`
	NetworkRxBytes  uint64  `json:"networkRxBytes"`
	NetworkTxBytes  uint64  `json:"networkTxBytes"`
}

func (n NodeMetric) String() string {
	s, _ := json.Marshal(n)
	return string(s)
}

type PodMetric struct {
	Namespace      string `json:"namespace"`
	UID            string `json:"uid"`
	CPUUsageUsec   uint64 `json:"cpuUsageUsec"`
	MemoryUsage    uint64 `json:"memoryUsage"`
	DiskReadBytes  uint64 `json:"diskReadBytes"`
	DiskWriteBytes uint64 `json:"diskWriteBytes"`
	NetworkRxBytes uint64 `json:"networkRxBytes"`
	NetworkTxBytes uint64 `json:"networkTxBytes"`
}

func (p PodMetric) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}
