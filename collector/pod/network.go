package pod

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shirou/gopsutil/v4/net"
)

func CollectPodNetworkMetric(podPath string) (rxBytes, txBytes uint64, err error) {
	containerPid, err := getContainerPID(podPath)
	if err != nil {
		return 0, 0, err
	}

	path := fmt.Sprintf("/proc/%d/net/dev", containerPid)
	netIOCounters, err := net.IOCountersByFileWithContext(context.Background(), false, path)
	if err != nil {
		return 0, 0, err
	}
	return netIOCounters[0].BytesRecv, netIOCounters[0].BytesSent, nil
}

// getContainerPID는 파드에 속한 컨테이너 중 하나의 PID를 반환합니다
func getContainerPID(podPath string) (int, error) {
	entries, err := os.ReadDir(podPath)
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 64 {
			procsFile := filepath.Join(podPath, entry.Name(), "cgroup.procs")
			data, err := os.ReadFile(procsFile)
			if err != nil {
				continue
			}
			var pid int
			_, err = fmt.Sscanf(string(data), "%d", &pid)
			if err == nil && pid > 0 {
				return pid, nil
			}
		}
	}
	return 0, fmt.Errorf("no container PID found in pod cgroup: %s", podPath)
}
