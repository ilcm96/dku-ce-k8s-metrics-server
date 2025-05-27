package pod

import "github.com/containerd/cgroups/v3/cgroup2/stats"

func CollectPodDiskMetric(entries []*stats.IOEntry) (readBytes, writeBytes uint64) {
	for _, entry := range entries {
		readBytes += entry.Rbytes
		writeBytes += entry.Wbytes
	}
	return readBytes, writeBytes
}
