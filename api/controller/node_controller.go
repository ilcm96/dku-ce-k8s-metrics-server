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
func (c *nodeController) GetMetricsByNodeName(ctx *fiber.Ctx) error {
	nodeName := ctx.Params("nodeName")
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
