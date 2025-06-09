package service

import (
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/dto"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/repository"
	"log/slog"
)

type PodService interface {
	FindAll() ([]*dto.PodMetricsResponse, error)
	FindByPodName(podName string) (*dto.PodMetricsResponse, error)
	FindByNodeName(nodeName string) ([]*dto.PodMetricsResponse, error)
}

type podService struct {
	podRepository repository.PodRepository
}

func NewPodService(podRepository repository.PodRepository) PodService {
	return &podService{
		podRepository: podRepository,
	}
}

// FindAll 는 모든 파드들에 대해 최신 메트릭을 제공합니다.
func (s *podService) FindAll() ([]*dto.PodMetricsResponse, error) {
	// 모든 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	metrics, err := s.podRepository.FindAll()
	if err != nil {
		slog.Error("failed to get pod metrics list", "error", err)
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

		response := &dto.PodMetricsResponse{}
		response.Build(cpuMillicores, latest)

		responses = append(responses, response)
	}

	return responses, nil
}

// FindByPodName 는 주어진 파드명에 대해 최신 메트릭을 제공합니다.
func (s *podService) FindByPodName(podName string) (*dto.PodMetricsResponse, error) {
	metrics, err := s.podRepository.FindByPodName(podName)
	if err != nil {
		slog.Error("failed to get pod metrics by pod name", "pod", podName, "error", err)
	}
	if len(metrics) < 2 {
		return nil, nil
	}

	latest := metrics[0]
	previous := metrics[1]

	cpuMillicores := calculatePodCpuMillicores(latest, previous)

	response := &dto.PodMetricsResponse{}
	response.Build(cpuMillicores, latest)

	return response, nil
}

// FindByNodeName 는 주어진 노드명에 대해 모든 파드의 최신 메트릭을 제공합니다.
func (s *podService) FindByNodeName(nodeName string) ([]*dto.PodMetricsResponse, error) {
	// 주어진 노드명을 가지는 모든 파드에 대해 가장 최근의 2개의 메트릭을 조회합니다.
	metrics, err := s.podRepository.FindByNodeName(nodeName)
	if err != nil {
		slog.Error("failed to get pod metrics by node name", "nodeName", nodeName, "error", err)
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

		response := &dto.PodMetricsResponse{}
		response.Build(cpuMillicores, latest)

		responses = append(responses, response)
	}

	return responses, nil
}

// calculateCpuMillicores 는 이전 메트릭과 최신 메트릭을 비교하여 CPU 밀리코어를 계산합니다.
func calculatePodCpuMillicores(latest, previous *entity.PodMetrics) float64 {
	if latest == nil || previous == nil {
		return 0.0
	}

	deltaCpuUsage := latest.CPUUsageUsec - previous.CPUUsageUsec
	interval := latest.Timestamp.Sub(previous.Timestamp).Seconds()

	if interval <= 0 {
		return 0.0
	}

	return float64(deltaCpuUsage) / (interval * 1e3)
}
