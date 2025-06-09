package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/service"
)

type PodController interface {
	GetMetricsList(ctx *fiber.Ctx) error
	GetMetricsByPodName(ctx *fiber.Ctx) error
}

type podController struct {
	podService service.PodService
}

func NewPodController(podService service.PodService) PodController {
	return &podController{
		podService: podService,
	}
}

// GetMetricsList 는 모든 파드의 최신 메트릭을 제공합니다.
func (c *podController) GetMetricsList(ctx *fiber.Ctx) error {
	metrics, err := c.podService.FindAll()
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(metrics)
}

// GetMetricsByPodName 는 특정 파드의 최신 메트릭을 제공합니다.
// window 쿼리 파라미터가 있으면 시계열 조회, 없으면 실시간 조회를 수행합니다.
func (c *podController) GetMetricsByPodName(ctx *fiber.Ctx) error {
	podName := ctx.Params("podName")
	window := ctx.Query("window")

	// window 파라미터가 있으면 시계열 조회
	if window != "" {
		timeSeriesMetrics, err := c.podService.FindTimeSeriesByPodName(podName, window)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if timeSeriesMetrics == nil {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		return ctx.JSON(timeSeriesMetrics)
	}

	// window 파라미터가 없으면 기존 실시간 조회
	metrics, err := c.podService.FindByPodName(podName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}
