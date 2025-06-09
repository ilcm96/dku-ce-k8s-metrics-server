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
func (c *podController) GetMetricsByPodName(ctx *fiber.Ctx) error {
	podName := ctx.Params("podName")
	metrics, err := c.podService.FindByPodName(podName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}
