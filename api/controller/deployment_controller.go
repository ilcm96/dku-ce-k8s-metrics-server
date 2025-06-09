package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/service"
)

type DeploymentController interface {
	GetDeploymentsByNamespaceName(ctx *fiber.Ctx) error
	GetMetricsByDeploymentName(ctx *fiber.Ctx) error
	GetPodMetricsByDeploymentName(ctx *fiber.Ctx) error
}

type deploymentController struct {
	deploymentService service.DeploymentService
}

func NewDeploymentController(deploymentService service.DeploymentService) DeploymentController {
	return &deploymentController{
		deploymentService: deploymentService,
	}
}

// GetDeploymentsByNamespaceName 는 특정 네임스페이스의 모든 디플로이먼트와 리소스 사용량을 제공합니다.
func (c *deploymentController) GetDeploymentsByNamespaceName(ctx *fiber.Ctx) error {
	namespaceName := ctx.Params("namespaceName")
	metrics, err := c.deploymentService.FindByNamespaceName(namespaceName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}

// GetMetricsByDeploymentName 는 특정 디플로이먼트의 리소스 사용량을 제공합니다.
func (c *deploymentController) GetMetricsByDeploymentName(ctx *fiber.Ctx) error {
	namespaceName := ctx.Params("namespaceName")
	deploymentName := ctx.Params("deploymentName")
	metrics, err := c.deploymentService.FindByDeploymentName(namespaceName, deploymentName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}

// GetPodMetricsByDeploymentName 는 특정 디플로이먼트의 모든 파드 목록과 리소스 사용량을 제공합니다.
func (c *deploymentController) GetPodMetricsByDeploymentName(ctx *fiber.Ctx) error {
	namespaceName := ctx.Params("namespaceName")
	deploymentName := ctx.Params("deploymentName")
	metrics, err := c.deploymentService.FindPodsByDeploymentName(namespaceName, deploymentName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if metrics == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(metrics)
}
