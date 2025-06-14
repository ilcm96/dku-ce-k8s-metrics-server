package service

import (
	"log/slog"
	"time"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/dto"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/repository"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/utils"
)

type NamespaceService interface {
	FindAll() ([]*dto.NamespaceMetricsResponse, error)
	FindByNamespaceName(namespaceName string) (*dto.NamespaceMetricsResponse, error)
	FindPodsByNamespaceName(namespaceName string) ([]*dto.PodMetricsResponse, error)
	FindTimeSeriesByNamespaceName(namespaceName, window string) (*dto.NamespaceTimeSeriesResponse, error)
}

type namespaceService struct {
	namespaceRepository  repository.NamespaceRepository
	timeSeriesCalculator TimeSeriesCalculator
}

func NewNamespaceService(namespaceRepository repository.NamespaceRepository) NamespaceService {
	return &namespaceService{
		namespaceRepository:  namespaceRepository,
		timeSeriesCalculator: NewTimeSeriesCalculator(),
	}
}

// FindAll 는 모든 네임스페이스들에 대해 집계된 메트릭을 제공합니다.
func (s *namespaceService) FindAll() ([]*dto.NamespaceMetricsResponse, error) {
	// 모든 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	allPodMetrics, err := s.namespaceRepository.FindAll()
	if err != nil {
		slog.Error("failed to get all pod metrics", "error", err)
		return nil, err
	}
	if len(allPodMetrics) == 0 {
		return nil, nil
	}

	// 네임스페이스별로 파드 메트릭을 그룹화합니다.
	namespaceMetricsMap := make(map[string][]*entity.PodMetrics)
	for _, metric := range allPodMetrics {
		if metric.NamespaceName != "" {
			namespaceMetricsMap[metric.NamespaceName] = append(namespaceMetricsMap[metric.NamespaceName], metric)
		}
	}

	// 각 네임스페이스에 대해 집계 계산을 수행합니다.
	var responses []*dto.NamespaceMetricsResponse
	for namespaceName, podMetrics := range namespaceMetricsMap {
		aggregatedMetrics := calculateNamespaceMetrics(namespaceName, podMetrics)
		if aggregatedMetrics != nil {
			responses = append(responses, aggregatedMetrics)
		}
	}

	return responses, nil
}

// FindByNamespaceName 는 주어진 네임스페이스명에 대해 집계된 메트릭을 제공합니다.
func (s *namespaceService) FindByNamespaceName(namespaceName string) (*dto.NamespaceMetricsResponse, error) {
	// 특정 네임스페이스의 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	podMetrics, err := s.namespaceRepository.FindByNamespaceName(namespaceName)
	if err != nil {
		slog.Error("failed to get pod metrics by namespace name", "namespaceName", namespaceName, "error", err)
		return nil, err
	}
	if len(podMetrics) == 0 {
		return nil, nil
	}

	// 네임스페이스 집계 계산을 수행합니다.
	aggregatedMetrics := calculateNamespaceMetrics(namespaceName, podMetrics)
	return aggregatedMetrics, nil
}

// FindPodsByNamespaceName 는 주어진 네임스페이스명을 가진 모든 파드의 최신 메트릭을 제공합니다.
func (s *namespaceService) FindPodsByNamespaceName(namespaceName string) ([]*dto.PodMetricsResponse, error) {
	// 주어진 네임스페이스명을 가지는 모든 파드에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	metrics, err := s.namespaceRepository.FindByNamespaceName(namespaceName)
	if err != nil {
		slog.Error("failed to get pod metrics by namespace name", "namespaceName", namespaceName, "error", err)
		return nil, err
	}
	if len(metrics) == 0 {
		return nil, nil
	}

	// 파드 이름별로 메트릭을 그룹화합니다.
	metricsMap := make(map[string][]*entity.PodMetrics)
	for _, metric := range metrics {
		metricsMap[metric.PodName] = append(metricsMap[metric.PodName], metric)
	}

	// 각 파드에 대해 가장 최근의 2개의 메트릭을 비교하여 응답을 생성합니다.
	var responses []*dto.PodMetricsResponse
	for _, podMetrics := range metricsMap {
		if len(podMetrics) < 2 {
			continue // 최소 2개의 메트릭이 있어야 비교 가능
		}

		latest := podMetrics[0]
		previous := podMetrics[1]

		cpuMillicores := calculatePodCpuMillicores(latest, previous)

		var deploymentName *string
		if latest.DeploymentName.Valid {
			deploymentName = &latest.DeploymentName.String
		}

		response := &dto.PodMetricsResponse{
			Timestamp:      latest.Timestamp,
			PodName:        latest.PodName,
			DeploymentName: deploymentName,
			NamespaceName:  latest.NamespaceName,
			NodeName:       latest.NodeName,
			UID:            latest.UID,
			CpuMillicores:  cpuMillicores,
			MemoryBytes:    latest.MemoryUsage,
			DiskReadBytes:  latest.DiskReadBytes,
			DiskWriteBytes: latest.DiskWriteBytes,
			NetworkRxBytes: latest.NetworkRxBytes,
			NetworkTxBytes: latest.NetworkTxBytes,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// calculateNamespaceMetrics 는 네임스페이스의 파드 메트릭들을 집계하여 네임스페이스 메트릭을 계산합니다.
func calculateNamespaceMetrics(namespaceName string, podMetrics []*entity.PodMetrics) *dto.NamespaceMetricsResponse {
	if len(podMetrics) == 0 {
		return nil
	}

	// 파드 이름별로 메트릭을 그룹화합니다.
	podMetricsMap := make(map[string][]*entity.PodMetrics)
	for _, metric := range podMetrics {
		podMetricsMap[metric.PodName] = append(podMetricsMap[metric.PodName], metric)
	}

	var totalCpuMillicores float64
	var totalMemoryBytes int64
	var totalDiskReadBytes int64
	var totalDiskWriteBytes int64
	var totalNetworkRxBytes int64
	var totalNetworkTxBytes int64
	var latestTimestamp time.Time
	var activePodCount int

	// 각 파드에 대해 메트릭을 계산하고 집계합니다.
	for _, podMetricsSlice := range podMetricsMap {
		if len(podMetricsSlice) < 2 {
			continue // 최소 2개의 메트릭이 있어야 비교 가능
		}

		latest := podMetricsSlice[0]
		previous := podMetricsSlice[1]

		// CPU 밀리코어 계산
		cpuMillicores := calculatePodCpuMillicores(latest, previous)
		totalCpuMillicores += cpuMillicores

		// 최신 메트릭 값들을 집계
		totalMemoryBytes += latest.MemoryUsage
		totalDiskReadBytes += latest.DiskReadBytes
		totalDiskWriteBytes += latest.DiskWriteBytes
		totalNetworkRxBytes += latest.NetworkRxBytes
		totalNetworkTxBytes += latest.NetworkTxBytes

		// 최신 타임스탬프 추적
		if latest.Timestamp.After(latestTimestamp) {
			latestTimestamp = latest.Timestamp
		}

		activePodCount++
	}

	// 네임스페이스 메트릭 응답 생성
	return &dto.NamespaceMetricsResponse{
		NamespaceName:  namespaceName,
		Timestamp:      latestTimestamp,
		CpuMillicores:  totalCpuMillicores,
		MemoryBytes:    totalMemoryBytes,
		DiskReadBytes:  totalDiskReadBytes,
		DiskWriteBytes: totalDiskWriteBytes,
		NetworkRxBytes: totalNetworkRxBytes,
		NetworkTxBytes: totalNetworkTxBytes,
		PodCount:       activePodCount,
	}
}

// FindTimeSeriesByNamespaceName 는 주어진 네임스페이스명과 윈도우에 대해 시계열 메트릭을 제공합니다.
func (s *namespaceService) FindTimeSeriesByNamespaceName(namespaceName, window string) (*dto.NamespaceTimeSeriesResponse, error) {
	// 윈도우 파라미터 파싱
	windowSpec, err := utils.ParseWindow(window)
	if err != nil {
		slog.Error("failed to parse window parameter", "window", window, "error", err)
		return nil, err
	}

	// 시간 범위 계산 (UTC 변환)
	endTime := time.Now().UTC()
	startTime := windowSpec.GetStartTime(endTime)

	// 시간 범위 내의 네임스페이스 파드 메트릭 조회 (UTC 시간으로 조회)
	metrics, err := s.namespaceRepository.FindByNamespaceNameInTimeWindow(namespaceName, startTime, endTime)
	if err != nil {
		slog.Error("failed to get namespace pod metrics in time window", "namespaceName", namespaceName, "startTime", startTime, "endTime", endTime, "error", err)
		return nil, err
	}

	if len(metrics) == 0 {
		return nil, nil
	}

	// 네임스페이스 시계열 계산
	response, err := s.timeSeriesCalculator.CalculateNamespaceTimeSeries(namespaceName, metrics, windowSpec)
	if err != nil {
		slog.Error("failed to calculate namespace time series", "namespaceName", namespaceName, "error", err)
		return nil, err
	}

	return response, nil
}
