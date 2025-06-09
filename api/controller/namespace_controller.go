package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/service"
)

type NamespaceController interface {
	GetMetricsList(ctx *fiber.Ctx) error
	GetMetricsByNamespaceName(ctx *fiber.Ctx) error
	GetPodMetricsListByNamespaceName(ctx *fiber.Ctx) error
}

type namespaceController struct {
	namespaceService service.NamespaceService
}

func NewNamespaceController(namespaceService service.NamespaceService) NamespaceController {
	return &namespaceController{
		namespaceService: namespaceService,
	}
}

// GetMetricsList 는 모든 네임스페이스의 집계된 메트릭을 제공합니다.
func (c *namespaceController) GetMetricsList(ctx *fiber.Ctx) error {
	metrics, err := c.namespaceService.FindAll()
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(metrics)
}

// GetMetricsByNamespaceName 는 특정 네임스페이스의 집계된 메트릭을 제공합니다.
func (c *namespaceController) GetMetricsByNamespaceName(ctx *fiber.Ctx) error {
	namespaceName := ctx.Params("namespaceName")
	metrics, err := c.namespaceService.FindByNamespaceName(namespaceName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}

// GetPodMetricsListByNamespaceName 는 특정 네임스페이스에 존재하는 모든 파드의 최신 메트릭을 조회합니다.
func (c *namespaceController) GetPodMetricsListByNamespaceName(ctx *fiber.Ctx) error {
	namespaceName := ctx.Params("namespaceName")
	metrics, err := c.namespaceService.FindPodsByNamespaceName(namespaceName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}
