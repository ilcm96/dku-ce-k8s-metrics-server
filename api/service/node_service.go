package service

import (
	"log/slog"
	"time"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/dto"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/repository"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/utils"
)

type NodeService interface {
	FindAll() ([]*dto.NodeMetricsResponse, error)
	FindByNodeName(nodeName string) (*dto.NodeMetricsResponse, error)
	FindTimeSeriesByNodeName(nodeName, window string) (*dto.NodeTimeSeriesResponse, error)
}

type nodeService struct {
	nodeRepository       repository.NodeRepository
	timeSeriesCalculator TimeSeriesCalculator
}

func NewNodeService(nodeRepository repository.NodeRepository) NodeService {
	return &nodeService{
		nodeRepository:       nodeRepository,
		timeSeriesCalculator: NewTimeSeriesCalculator(),
	}
}

// FindAll 는 모든 노드들에 대해 최신 메트릭을 제공합니다.
func (s *nodeService) FindAll() ([]*dto.NodeMetricsResponse, error) {
	// 모든 노드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	metrics, err := s.nodeRepository.FindAll()
	if err != nil {
		slog.Error("failed to get node metrics list", "error", err)
		return nil, err
	}
	if len(metrics) == 0 {
		return nil, nil
	}

	// 노드 이름별로 메트릭을 그룹화합니다.
	metricsMap := make(map[string][]*entity.NodeMetrics)
	for _, metric := range metrics {
		metricsMap[metric.NodeName] = append(metricsMap[metric.NodeName], metric)
	}

	// 각 노드에 대해 가장 최근의 2개의 메트릭을 비교하여 응답을 생성합니다.
	var responses []*dto.NodeMetricsResponse
	for _, nodeMetrics := range metricsMap {
		if len(nodeMetrics) < 2 {
			continue // 최소 2개의 메트릭이 있어야 비교 가능
		}

		latest := nodeMetrics[0]
		previous := nodeMetrics[1]

		cpuMillicores := calculateNodeCpuMillicores(latest, previous)
		memoryBytes := latest.MemoryTotal - latest.MemoryAvailable

		response := &dto.NodeMetricsResponse{
			Timestamp:      latest.Timestamp,
			NodeName:       latest.NodeName,
			CpuMillicores:  cpuMillicores,
			MemoryBytes:    memoryBytes,
			DiskReadBytes:  latest.DiskReadBytes,
			DiskWriteBytes: latest.DiskWriteBytes,
			NetworkRxBytes: latest.NetworkRxBytes,
			NetworkTxBytes: latest.NetworkTxBytes,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// FindByNodeName 는 주어진 노드명에 대해 최신 메트릭을 제공합니다.
func (s *nodeService) FindByNodeName(nodeName string) (*dto.NodeMetricsResponse, error) {
	metrics, err := s.nodeRepository.FindByNodeName(nodeName)
	if err != nil {
		slog.Error("failed to get node metrics by node name", "nodeName", nodeName, "error", err)
		return nil, err
	}
	if len(metrics) < 2 {
		return nil, nil // 최소 2개의 메트릭이 있어야 비교 가능
	}

	latest := metrics[0]
	previous := metrics[1]

	cpuMillicores := calculateNodeCpuMillicores(latest, previous)
	memoryBytes := latest.MemoryTotal - latest.MemoryAvailable

	response := &dto.NodeMetricsResponse{
		Timestamp:      latest.Timestamp,
		NodeName:       latest.NodeName,
		CpuMillicores:  cpuMillicores,
		MemoryBytes:    memoryBytes,
		DiskReadBytes:  latest.DiskReadBytes,
		DiskWriteBytes: latest.DiskWriteBytes,
		NetworkRxBytes: latest.NetworkRxBytes,
		NetworkTxBytes: latest.NetworkTxBytes,
	}

	return response, nil
}

// FindTimeSeriesByNodeName 는 주어진 노드명과 윈도우에 대해 시계열 메트릭을 제공합니다.
func (s *nodeService) FindTimeSeriesByNodeName(nodeName, window string) (*dto.NodeTimeSeriesResponse, error) {
	// 윈도우 파라미터 파싱
	windowSpec, err := utils.ParseWindow(window)
	if err != nil {
		slog.Error("failed to parse window parameter", "window", window, "error", err)
		return nil, err
	}

	// 시간 범위 계산 (UTC 변환)
	endTime := time.Now().UTC()
	startTime := windowSpec.GetStartTime(endTime)

	// 시간 범위 내의 메트릭 조회 (UTC 시간으로 조회)
	metrics, err := s.nodeRepository.FindByNodeNameInTimeWindow(nodeName, startTime, endTime)
	if err != nil {
		slog.Error("failed to get node metrics in time window", "nodeName", nodeName, "startTime", startTime, "endTime", endTime, "error", err)
		return nil, err
	}

	if len(metrics) == 0 {
		return nil, nil
	}

	// 시계열 계산
	response, err := s.timeSeriesCalculator.CalculateNodeTimeSeries(nodeName, metrics, windowSpec)
	if err != nil {
		slog.Error("failed to calculate node time series", "nodeName", nodeName, "error", err)
		return nil, err
	}

	return response, nil
}

// calculateCpuMillicores 는 두 개의 NodeMetrics 객체를 비교하여 CPU 사용량을 밀리코어 단위로 계산합니다.
func calculateNodeCpuMillicores(latest, previous *entity.NodeMetrics) float64 {
	if latest == nil || previous == nil {
		return 0.0
	}

	deltaCpuBusy := latest.CPUBusy - previous.CPUBusy
	deltaCpuTotal := latest.CPUTotal - previous.CPUTotal

	if deltaCpuTotal <= 0 {
		return 0.0
	}

	return (deltaCpuBusy / deltaCpuTotal) * 1000
}
