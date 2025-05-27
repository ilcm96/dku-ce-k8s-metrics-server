package pod

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/collector/metadata"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/shared/types"
)

func CollectPodMetrics() ([]types.PodMetric, error) {
	podPaths, err := getPodCgroupPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to get pod cgroup paths: %w", err)
	}

	var podMetrics []types.PodMetric
	for _, podPath := range podPaths {
		metric, err := collectSinglePodMetric(podPath)
		if err != nil {
			log.Printf("failed to collect metrics for pod %s: %v", podPath, err)
			continue
		}
		podMetrics = append(podMetrics, metric)
	}

	return podMetrics, nil
}

// collectSinglePodMetric은 단일 파드의 메트릭을 수집합니다
func collectSinglePodMetric(podPath string) (types.PodMetric, error) {
	resource := podPath[14:]

	// cgroup 매니저 로드
	manager, err := cgroup2.Load(resource)
	if err != nil {
		return types.PodMetric{}, fmt.Errorf("failed to load cgroup manager: %w", err)
	}

	// 기본 메트릭 수집
	metrics, err := manager.Stat()
	if err != nil {
		return types.PodMetric{}, fmt.Errorf("failed to get cgroup stats: %w", err)
	}

	// 디스크 메트릭 수집
	diskReadBytes, diskWriteBytes := CollectPodDiskMetric(metrics.Io.Usage)

	// 네트워크 메트릭 수집
	networkRxBytes, networkTxBytes, err := CollectPodNetworkMetric(podPath)
	if err != nil {
		return types.PodMetric{}, fmt.Errorf("failed to collect network metrics: %w", err)
	}

	// UID 추출
	uid := extractPodUID(resource)

	return types.PodMetric{
		Namespace:      metadata.Namespace,
		UID:            uid,
		CPUUsageUsec:   metrics.CPU.UsageUsec,
		MemoryUsage:    metrics.Memory.Usage,
		DiskReadBytes:  diskReadBytes,
		DiskWriteBytes: diskWriteBytes,
		NetworkRxBytes: networkRxBytes,
		NetworkTxBytes: networkTxBytes,
	}, nil
}

// extractPodUID는 resource 경로에서 파드 UID를 추출합니다
func extractPodUID(resource string) string {
	if len(resource) >= 36 {
		return resource[len(resource)-36:]
	}
	return ""
}

func getPodCgroupPaths() ([]string, error) {
	var pods []string
	err := filepath.Walk("/sys/fs/cgroup/kubepods/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && len(info.Name()) > 3 && info.Name()[:3] == "pod" {
			pods = append(pods, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return pods, nil
}
