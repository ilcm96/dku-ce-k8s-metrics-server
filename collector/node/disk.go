package node

import (
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

type NodeDiskMetric struct {
	Timestamp  time.Time `json:"timestamp"`
	ReadBytes  uint64    `json:"readBytes"`
	WriteBytes uint64    `json:"writeBytes"`
}

func (d NodeDiskMetric) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func CollectNodeDiskMetric() (NodeDiskMetric, error) {
	diskIOCounters, err := disk.IOCounters("vda")
	if err != nil {
		return NodeDiskMetric{}, err
	}

	stat := diskIOCounters["vda"]

	return NodeDiskMetric{
		Timestamp:  time.Now(),
		ReadBytes:  stat.ReadBytes,
		WriteBytes: stat.WriteBytes,
	}, nil
}
