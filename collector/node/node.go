package node

import (
	"fmt"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/collector/metadata"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/shared/types"
)

func CollectNodeMetric() (types.NodeMetric, error) {
	// CPU Metric
	cpuMetric, err := CollectNodeCpuMetric()
	if err != nil {
		return types.NodeMetric{}, fmt.Errorf("failed to collect node CPU metric: %w", err)
	}

	// Memory Metric
	memoryMetric, err := CollectNodeMemoryMetric()
	if err != nil {
		return types.NodeMetric{}, fmt.Errorf("failed to collect node memory metric: %w", err)
	}

	// Disk Metric
	diskMetric, err := CollectNodeDiskMetric()
	if err != nil {
		return types.NodeMetric{}, fmt.Errorf("failed to collect node disk metric: %w", err)
	}

	// Network Metric
	networkMetric, err := CollectNodeNetworkMetric()
	if err != nil {
		return types.NodeMetric{}, fmt.Errorf("failed to collect node network metric: %w", err)
	}

	return types.NodeMetric{
		NodeName:        metadata.NodeName,
		CPUTotal:        cpuMetric.Total,
		CPUBusy:         cpuMetric.Busy,
		CPUCount:        cpuMetric.Count,
		MemoryTotal:     memoryMetric.Total,
		MemoryAvailable: memoryMetric.Available,
		MemoryUsed:      memoryMetric.Used,
		DiskReadBytes:   diskMetric.ReadBytes,
		DiskWriteBytes:  diskMetric.WriteBytes,
		NetworkRxBytes:  networkMetric.RxBytes,
		NetworkTxBytes:  networkMetric.TxBytes,
	}, nil
}
