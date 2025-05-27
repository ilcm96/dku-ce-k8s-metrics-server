package node

import (
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
)

type NodeMemoryMetric struct {
	Timestamp time.Time `json:"timestamp"`
	Total     uint64    `json:"total"`
	Available uint64    `json:"available,omitempty"`
	Used      uint64    `json:"used"`
}

func (m NodeMemoryMetric) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

func CollectNodeMemoryMetric() (NodeMemoryMetric, error) {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return NodeMemoryMetric{}, err
	}
	return NodeMemoryMetric{
		Timestamp: time.Now(),
		Total:     virtualMemory.Total,
		Available: virtualMemory.Available,
		Used:      virtualMemory.Used,
	}, nil
}
