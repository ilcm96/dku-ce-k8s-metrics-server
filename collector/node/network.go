package node

import (
	"encoding/json"

	"github.com/shirou/gopsutil/v4/net"
)

type NodeNetworkMetric struct {
	RxBytes uint64 `json:"rxBytes"`
	TxBytes uint64 `json:"txBytes"`
}

func (n NodeNetworkMetric) String() string {
	s, _ := json.Marshal(n)
	return string(s)
}

func CollectNodeNetworkMetric() (NodeNetworkMetric, error) {
	netIOCounters, err := net.IOCounters(false)
	if err != nil {
		return NodeNetworkMetric{}, err
	}

	netIOCounter := netIOCounters[0]

	networkMetric := NodeNetworkMetric{
		RxBytes: netIOCounter.BytesRecv,
		TxBytes: netIOCounter.BytesSent,
	}

	return networkMetric, nil
}
