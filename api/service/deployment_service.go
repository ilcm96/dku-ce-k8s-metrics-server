package service

import (
	"log/slog"
	"time"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/dto"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/repository"
)

type DeploymentService interface {
	FindByNamespaceName(namespaceName string) ([]*dto.DeploymentMetricsResponse, error)
	FindByDeploymentName(namespaceName, deploymentName string) (*dto.DeploymentMetricsResponse, error)
	FindPodsByDeploymentName(namespaceName, deploymentName string) ([]*dto.PodMetricsResponse, error)
}

type deploymentService struct {
	deploymentRepository repository.DeploymentRepository
}

func NewDeploymentService(deploymentRepository repository.DeploymentRepository) DeploymentService {
	return &deploymentService{
		deploymentRepository: deploymentRepository,
	}
}

// FindByNamespaceName 는 주어진 네임스페이스의 모든 디플로이먼트에 대해 집계된 메트릭을 제공합니다.
func (s *deploymentService) FindByNamespaceName(namespaceName string) ([]*dto.DeploymentMetricsResponse, error) {
	// 특정 네임스페이스의 디플로이먼트 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	allPodMetrics, err := s.deploymentRepository.FindByNamespaceName(namespaceName)
	if err != nil {
		slog.Error("failed to get deployment pod metrics by namespace name", "namespaceName", namespaceName, "error", err)
		return nil, err
	}
	if len(allPodMetrics) == 0 {
		return nil, nil
	}

	// 디플로이먼트별로 파드 메트릭을 그룹화합니다.
	deploymentMetricsMap := make(map[string][]*entity.PodMetrics)
	for _, metric := range allPodMetrics {
		if metric.DeploymentName.Valid && metric.DeploymentName.String != "" {
			deploymentMetricsMap[metric.DeploymentName.String] = append(deploymentMetricsMap[metric.DeploymentName.String], metric)
		}
	}

	// 각 디플로이먼트에 대해 집계 계산을 수행합니다.
	var responses []*dto.DeploymentMetricsResponse
	for deploymentName, podMetrics := range deploymentMetricsMap {
		aggregatedMetrics := calculateDeploymentMetrics(namespaceName, deploymentName, podMetrics)
		if aggregatedMetrics != nil {
			responses = append(responses, aggregatedMetrics)
		}
	}

	return responses, nil
}

// FindByDeploymentName 는 주어진 디플로이먼트에 대해 집계된 메트릭을 제공합니다.
func (s *deploymentService) FindByDeploymentName(namespaceName, deploymentName string) (*dto.DeploymentMetricsResponse, error) {
	// 특정 디플로이먼트의 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	podMetrics, err := s.deploymentRepository.FindByDeploymentName(namespaceName, deploymentName)
	if err != nil {
		slog.Error("failed to get pod metrics by deployment name", "namespaceName", namespaceName, "deploymentName", deploymentName, "error", err)
		return nil, err
	}
	if len(podMetrics) == 0 {
		return nil, nil
	}

	// 디플로이먼트 집계 계산을 수행합니다.
	aggregatedMetrics := calculateDeploymentMetrics(namespaceName, deploymentName, podMetrics)
	return aggregatedMetrics, nil
}

// FindPodsByDeploymentName 는 주어진 디플로이먼트의 모든 파드에 대해 최신 메트릭을 제공합니다.
func (s *deploymentService) FindPodsByDeploymentName(namespaceName, deploymentName string) ([]*dto.PodMetricsResponse, error) {
	// 주어진 디플로이먼트의 모든 파드에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	metrics, err := s.deploymentRepository.FindByDeploymentName(namespaceName, deploymentName)
	if err != nil {
		slog.Error("failed to get pod metrics by deployment name", "namespaceName", namespaceName, "deploymentName", deploymentName, "error", err)
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

// calculateDeploymentMetrics 는 디플로이먼트의 파드 메트릭들을 집계하여 디플로이먼트 메트릭을 계산합니다.
func calculateDeploymentMetrics(namespaceName, deploymentName string, podMetrics []*entity.PodMetrics) *dto.DeploymentMetricsResponse {
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

	// 디플로이먼트 메트릭 응답 생성
	return &dto.DeploymentMetricsResponse{
		DeploymentName: deploymentName,
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
