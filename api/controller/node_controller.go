package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/service"
)

type NodeController interface {
	GetMetricsList(ctx *fiber.Ctx) error
	GetMetricsByNodeName(ctx *fiber.Ctx) error
	GetPodMetricsListByNodeName(ctx *fiber.Ctx) error
}

type nodeController struct {
	nodeService service.NodeService
	podService  service.PodService
}

func NewNodeController(nodeService service.NodeService, podService service.PodService) NodeController {
	return &nodeController{
		nodeService: nodeService,
		podService:  podService,
	}
}

// GetMetricsList 는 모든 노드의 최신 메트릭을 제공합니다.
func (c *nodeController) GetMetricsList(ctx *fiber.Ctx) error {
	metrics, err := c.nodeService.FindAll()
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(metrics)
}

// GetMetricsByNodeName 은 특정 노드의 최신 메트릭을 제공합니다.
// window 쿼리 파라미터가 있으면 시계열 조회, 없으면 실시간 조회를 수행합니다.
func (c *nodeController) GetMetricsByNodeName(ctx *fiber.Ctx) error {
	nodeName := ctx.Params("nodeName")
	window := ctx.Query("window")

	// window 파라미터가 있으면 시계열 조회
	if window != "" {
		timeSeriesMetrics, err := c.nodeService.FindTimeSeriesByNodeName(nodeName, window)
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
	metrics, err := c.nodeService.FindByNodeName(nodeName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}

// GetPodMetricsListByNodeName 은 특정 노드에 존재하는 모든 파드의 최신 메트릭을 조회합니다.
func (c *nodeController) GetPodMetricsListByNodeName(ctx *fiber.Ctx) error {
	nodeName := ctx.Params("nodeName")
	metrics, err := c.podService.FindByNodeName(nodeName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}
