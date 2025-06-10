package service

import (
	"fmt"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/dto"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/utils"
)

type TimeSeriesCalculator interface {
	CalculateNodeTimeSeries(nodeName string, metrics []*entity.NodeMetrics, window *utils.WindowSpec) (*dto.NodeTimeSeriesResponse, error)
	CalculatePodTimeSeries(podName string, metrics []*entity.PodMetrics, window *utils.WindowSpec) (*dto.PodTimeSeriesResponse, error)
	CalculateNamespaceTimeSeries(namespaceName string, metrics []*entity.PodMetrics, window *utils.WindowSpec) (*dto.NamespaceTimeSeriesResponse, error)
}

type timeSeriesCalculator struct{}

func NewTimeSeriesCalculator() TimeSeriesCalculator {
	return &timeSeriesCalculator{}
}

// CalculateNodeTimeSeries 는 노드 메트릭들로부터 시계열 데이터를 계산합니다.
func (c *timeSeriesCalculator) CalculateNodeTimeSeries(nodeName string, metrics []*entity.NodeMetrics, window *utils.WindowSpec) (*dto.NodeTimeSeriesResponse, error) {
	if len(metrics) < 2 {
		return nil, fmt.Errorf("insufficient data points for time series calculation (need at least 2, got %d)", len(metrics))
	}

	// 가장 최근 시간을 endTime으로 설정
	endTime := metrics[0].Timestamp
	for _, metric := range metrics {
		if metric.Timestamp.After(endTime) {
			endTime = metric.Timestamp
		}
	}
	startTime := window.GetStartTime(endTime)

	// 평균값 계산
	avgCpuMillicores, avgMemoryBytes, avgDiskReadRate, avgDiskWriteRate, avgNetworkRxRate, avgNetworkTxRate, err := c.calculateNodeAverages(metrics, window)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate averages: %w", err)
	}

	response := &dto.NodeTimeSeriesResponse{
		NodeName:         nodeName,
		Window:           window.String(),
		StartTime:        startTime,
		EndTime:          endTime,
		AvgCpuMillicores: avgCpuMillicores,
		AvgMemoryBytes:   avgMemoryBytes,
		AvgDiskReadRate:  avgDiskReadRate,
		AvgDiskWriteRate: avgDiskWriteRate,
		AvgNetworkRxRate: avgNetworkRxRate,
		AvgNetworkTxRate: avgNetworkTxRate,
	}

	return response, nil
}

// calculateNodeAverages 는 노드 메트릭들로부터 평균값들을 계산합니다.
func (c *timeSeriesCalculator) calculateNodeAverages(metrics []*entity.NodeMetrics, window *utils.WindowSpec) (float64, int64, float64, float64, float64, float64, error) {
	if len(metrics) < 2 {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("need at least 2 metrics for calculation")
	}

	windowDuration := window.ToDuration()
	var totalCpuMillicores float64
	var totalMemoryBytes int64
	var totalDiskReadRate float64
	var totalDiskWriteRate float64
	var totalNetworkRxRate float64
	var totalNetworkTxRate float64
	var count int

	// 각 연속된 메트릭 쌍에 대해 계산
	for i := 0; i < len(metrics)-1; i++ {
		current := metrics[i]
		next := metrics[i+1]

		// 윈도우 내의 데이터만 사용
		if current.Timestamp.Sub(next.Timestamp) > windowDuration {
			continue
		}

		// CPU 계산
		cpuMillicores := c.calculateNodeCpuMillicores(current, next)
		totalCpuMillicores += cpuMillicores

		// 메모리 계산 (현재값 사용)
		memoryBytes := current.MemoryTotal - current.MemoryAvailable
		totalMemoryBytes += memoryBytes

		// 디스크 I/O 속도 계산
		timeDiffSeconds := current.Timestamp.Sub(next.Timestamp).Seconds()
		if timeDiffSeconds <= 0 {
			timeDiffSeconds = 1
		}

		diskReadRate := float64(current.DiskReadBytes-next.DiskReadBytes) / timeDiffSeconds
		diskWriteRate := float64(current.DiskWriteBytes-next.DiskWriteBytes) / timeDiffSeconds
		networkRxRate := float64(current.NetworkRxBytes-next.NetworkRxBytes) / timeDiffSeconds
		networkTxRate := float64(current.NetworkTxBytes-next.NetworkTxBytes) / timeDiffSeconds

		// 음수 값 처리
		if diskReadRate < 0 {
			diskReadRate = 0
		}
		if diskWriteRate < 0 {
			diskWriteRate = 0
		}
		if networkRxRate < 0 {
			networkRxRate = 0
		}
		if networkTxRate < 0 {
			networkTxRate = 0
		}

		totalDiskReadRate += diskReadRate
		totalDiskWriteRate += diskWriteRate
		totalNetworkRxRate += networkRxRate
		totalNetworkTxRate += networkTxRate

		count++
	}

	if count == 0 {
		return 0, 0, 0, 0, 0, 0, nil
	}

	// 평균 계산
	avgCpuMillicores := totalCpuMillicores / float64(count)
	avgMemoryBytes := totalMemoryBytes / int64(count)
	avgDiskReadRate := totalDiskReadRate / float64(count)
	avgDiskWriteRate := totalDiskWriteRate / float64(count)
	avgNetworkRxRate := totalNetworkRxRate / float64(count)
	avgNetworkTxRate := totalNetworkTxRate / float64(count)

	return avgCpuMillicores, avgMemoryBytes, avgDiskReadRate, avgDiskWriteRate, avgNetworkRxRate, avgNetworkTxRate, nil
}

// calculateNodeCpuMillicores 는 두 개의 NodeMetrics 객체를 비교하여 CPU 사용량을 밀리코어 단위로 계산합니다.
func (c *timeSeriesCalculator) calculateNodeCpuMillicores(latest, previous *entity.NodeMetrics) float64 {
	if latest == nil || previous == nil {
		return 0.0
	}

	deltaCpuBusy := latest.CPUBusy - previous.CPUBusy
	deltaCpuTotal := latest.CPUTotal - previous.CPUTotal

	if deltaCpuTotal <= 0 {
		return 0.0
	}

	return (deltaCpuBusy / deltaCpuTotal) * 1000 * float64(latest.CPUCount)
}

// CalculatePodTimeSeries 는 파드 메트릭들로부터 시계열 데이터를 계산합니다.
func (c *timeSeriesCalculator) CalculatePodTimeSeries(podName string, metrics []*entity.PodMetrics, window *utils.WindowSpec) (*dto.PodTimeSeriesResponse, error) {
	if len(metrics) < 2 {
		return nil, fmt.Errorf("insufficient data points for time series calculation (need at least 2, got %d)", len(metrics))
	}

	// 가장 최근 시간을 endTime으로 설정
	endTime := metrics[0].Timestamp
	for _, metric := range metrics {
		if metric.Timestamp.After(endTime) {
			endTime = metric.Timestamp
		}
	}
	startTime := window.GetStartTime(endTime)

	// 평균값 계산
	avgCpuMillicores, avgMemoryBytes, avgDiskReadRate, avgDiskWriteRate, avgNetworkRxRate, avgNetworkTxRate, err := c.calculatePodAverages(metrics, window)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate averages: %w", err)
	}

	// 첫 번째 메트릭에서 추가 정보 가져오기
	firstMetric := metrics[0]

	// DeploymentName을 처리 (sql.NullString -> *string)
	var deploymentName *string
	if firstMetric.DeploymentName.Valid {
		deploymentName = &firstMetric.DeploymentName.String
	}

	response := &dto.PodTimeSeriesResponse{
		PodName:          podName,
		DeploymentName:   deploymentName,
		NamespaceName:    firstMetric.NamespaceName,
		NodeName:         firstMetric.NodeName,
		UID:              firstMetric.UID,
		Window:           window.String(),
		StartTime:        startTime,
		EndTime:          endTime,
		AvgCpuMillicores: avgCpuMillicores,
		AvgMemoryBytes:   avgMemoryBytes,
		AvgDiskReadRate:  avgDiskReadRate,
		AvgDiskWriteRate: avgDiskWriteRate,
		AvgNetworkRxRate: avgNetworkRxRate,
		AvgNetworkTxRate: avgNetworkTxRate,
	}

	return response, nil
}

// calculatePodAverages 는 파드 메트릭들로부터 평균값들을 계산합니다.
func (c *timeSeriesCalculator) calculatePodAverages(metrics []*entity.PodMetrics, window *utils.WindowSpec) (float64, int64, float64, float64, float64, float64, error) {
	if len(metrics) < 2 {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("need at least 2 metrics for calculation")
	}

	windowDuration := window.ToDuration()
	var totalCpuMillicores float64
	var totalMemoryBytes int64
	var totalDiskReadRate float64
	var totalDiskWriteRate float64
	var totalNetworkRxRate float64
	var totalNetworkTxRate float64
	var count int

	// 각 연속된 메트릭 쌍에 대해 계산
	for i := 0; i < len(metrics)-1; i++ {
		current := metrics[i]
		next := metrics[i+1]

		// 윈도우 내의 데이터만 사용
		if current.Timestamp.Sub(next.Timestamp) > windowDuration {
			continue
		}

		// CPU 계산
		cpuMillicores := c.calculatePodCpuMillicores(current, next)
		totalCpuMillicores += cpuMillicores

		// 메모리 계산 (현재값 사용)
		totalMemoryBytes += current.MemoryUsage

		// 디스크 I/O 속도 계산
		timeDiffSeconds := current.Timestamp.Sub(next.Timestamp).Seconds()
		if timeDiffSeconds <= 0 {
			timeDiffSeconds = 1
		}

		diskReadRate := float64(current.DiskReadBytes-next.DiskReadBytes) / timeDiffSeconds
		diskWriteRate := float64(current.DiskWriteBytes-next.DiskWriteBytes) / timeDiffSeconds
		networkRxRate := float64(current.NetworkRxBytes-next.NetworkRxBytes) / timeDiffSeconds
		networkTxRate := float64(current.NetworkTxBytes-next.NetworkTxBytes) / timeDiffSeconds

		// 음수 값 처리
		if diskReadRate < 0 {
			diskReadRate = 0
		}
		if diskWriteRate < 0 {
			diskWriteRate = 0
		}
		if networkRxRate < 0 {
			networkRxRate = 0
		}
		if networkTxRate < 0 {
			networkTxRate = 0
		}

		totalDiskReadRate += diskReadRate
		totalDiskWriteRate += diskWriteRate
		totalNetworkRxRate += networkRxRate
		totalNetworkTxRate += networkTxRate

		count++
	}

	if count == 0 {
		return 0, 0, 0, 0, 0, 0, nil
	}

	// 평균 계산
	avgCpuMillicores := totalCpuMillicores / float64(count)
	avgMemoryBytes := totalMemoryBytes / int64(count)
	avgDiskReadRate := totalDiskReadRate / float64(count)
	avgDiskWriteRate := totalDiskWriteRate / float64(count)
	avgNetworkRxRate := totalNetworkRxRate / float64(count)
	avgNetworkTxRate := totalNetworkTxRate / float64(count)

	return avgCpuMillicores, avgMemoryBytes, avgDiskReadRate, avgDiskWriteRate, avgNetworkRxRate, avgNetworkTxRate, nil
}

// calculatePodCpuMillicores 는 두 개의 PodMetrics 객체를 비교하여 CPU 사용량을 밀리코어 단위로 계산합니다.
// Pod의 경우 CPU 사용량이 마이크로초(usec) 단위로 저장되므로 다른 계산 방식을 사용합니다.
func (c *timeSeriesCalculator) calculatePodCpuMillicores(latest, previous *entity.PodMetrics) float64 {
	if latest == nil || previous == nil {
		return 0.0
	}

	deltaCpuUsage := latest.CPUUsageUsec - previous.CPUUsageUsec
	interval := latest.Timestamp.Sub(previous.Timestamp).Seconds()

	if interval <= 0 {
		return 0.0
	}

	// CPU 사용량 계산: (cpu_usage_usec_end - cpu_usage_usec_start) / time_diff_seconds / 1000
	return float64(deltaCpuUsage) / (interval * 1000)
}

// CalculateNamespaceTimeSeries 는 네임스페이스의 파드 메트릭들로부터 시계열 데이터를 계산합니다.
func (c *timeSeriesCalculator) CalculateNamespaceTimeSeries(namespaceName string, metrics []*entity.PodMetrics, window *utils.WindowSpec) (*dto.NamespaceTimeSeriesResponse, error) {
	if len(metrics) < 2 {
		return nil, fmt.Errorf("insufficient data points for time series calculation (need at least 2, got %d)", len(metrics))
	}

	// 가장 최근 시간을 endTime으로 설정
	endTime := metrics[0].Timestamp
	for _, metric := range metrics {
		if metric.Timestamp.After(endTime) {
			endTime = metric.Timestamp
		}
	}
	startTime := window.GetStartTime(endTime)

	// 평균값 계산
	avgCpuMillicores, avgMemoryBytes, avgDiskReadRate, avgDiskWriteRate, avgNetworkRxRate, avgNetworkTxRate, err := c.calculateNamespaceAverages(metrics, window)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate averages: %w", err)
	}

	response := &dto.NamespaceTimeSeriesResponse{
		NamespaceName:    namespaceName,
		Window:           window.String(),
		StartTime:        startTime,
		EndTime:          endTime,
		AvgCpuMillicores: avgCpuMillicores,
		AvgMemoryBytes:   avgMemoryBytes,
		AvgDiskReadRate:  avgDiskReadRate,
		AvgDiskWriteRate: avgDiskWriteRate,
		AvgNetworkRxRate: avgNetworkRxRate,
		AvgNetworkTxRate: avgNetworkTxRate,
	}

	return response, nil
}

// calculateNamespaceAverages 는 네임스페이스의 파드 메트릭들로부터 평균값들을 계산합니다.
func (c *timeSeriesCalculator) calculateNamespaceAverages(metrics []*entity.PodMetrics, window *utils.WindowSpec) (float64, int64, float64, float64, float64, float64, error) {
	if len(metrics) < 2 {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("need at least 2 metrics for calculation")
	}

	// 파드별로 메트릭을 그룹화
	podMetricsMap := make(map[string][]*entity.PodMetrics)
	for _, metric := range metrics {
		podMetricsMap[metric.PodName] = append(podMetricsMap[metric.PodName], metric)
	}

	var totalCpuMillicores float64
	var totalMemoryBytes int64
	var totalDiskReadRate float64
	var totalDiskWriteRate float64
	var totalNetworkRxRate float64
	var totalNetworkTxRate float64
	var count int

	// 각 파드에 대해 평균값 계산 후 집계
	for _, podMetrics := range podMetricsMap {
		if len(podMetrics) < 2 {
			continue // 최소 2개의 메트릭이 있어야 계산 가능
		}

		// 파드별 평균값 계산
		podAvgCpu, podAvgMem, podAvgDiskRead, podAvgDiskWrite, podAvgNetRx, podAvgNetTx, err := c.calculatePodAverages(podMetrics, window)
		if err != nil {
			continue
		}

		totalCpuMillicores += podAvgCpu
		totalMemoryBytes += podAvgMem
		totalDiskReadRate += podAvgDiskRead
		totalDiskWriteRate += podAvgDiskWrite
		totalNetworkRxRate += podAvgNetRx
		totalNetworkTxRate += podAvgNetTx
		count++
	}

	if count == 0 {
		return 0, 0, 0, 0, 0, 0, nil
	}

	// 네임스페이스 전체 평균값은 각 파드의 합계를 사용 (파드별 평균의 합)
	return totalCpuMillicores, totalMemoryBytes, totalDiskReadRate, totalDiskWriteRate, totalNetworkRxRate, totalNetworkTxRate, nil
}
