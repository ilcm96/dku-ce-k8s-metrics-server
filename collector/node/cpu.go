package node

import (
	"encoding/json"

	"github.com/shirou/gopsutil/v4/cpu"
)

type NodeCpuMetric struct {
	Total float64 `json:"total"`
	Busy  float64 `json:"busy"`
	Count int     `json:"count"`
}

func (c NodeCpuMetric) String() string {
	s, _ := json.Marshal(c)
	return string(s)
}

func CollectNodeCpuMetric() (NodeCpuMetric, error) {
	cpuTimes, err := cpu.Times(false)
	if err != nil {
		return NodeCpuMetric{}, err
	}

	cpuTime := cpuTimes[0]

	total := cpuTime.User +
		cpuTime.System +
		cpuTime.Idle +
		cpuTime.Nice +
		cpuTime.Iowait +
		cpuTime.Irq +
		cpuTime.Softirq

	busy := total - cpuTime.Idle - cpuTime.Iowait

	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return NodeCpuMetric{}, err
	}

	return NodeCpuMetric{
		Total: total,
		Busy:  busy,
		Count: cpuCount,
	}, nil
}
