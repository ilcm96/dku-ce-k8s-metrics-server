package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/controller"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/database"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/repository"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/service"
	"github.com/joho/godotenv"
	slogfiber "github.com/samber/slog-fiber"
	"log"
	"log/slog"
	"os"
)

func main() {
	// Logger 설정
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 환경변수 설정
	if os.Getenv("ENV") != "PRODUCTION" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file")
		}
		slog.Info("environment variables loaded successfully from .env file")
	}

	// 데이터베이스 연결
	database.Connect()
	db := database.GetConnection()

	// Fiber 앱 설정
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(slogfiber.New(logger))
	app.Use(recover.New())

	// 의존성 생성 및 주입
	nodeRepository := repository.NewNodeRepository(db)
	podRepository := repository.NewPodRepository(db)

	nodeService := service.NewNodeService(nodeRepository)
	podService := service.NewPodService(podRepository)

	nodeController := controller.NewNodeController(nodeService, podService)
	podController := controller.NewPodController(podService)

	// 라우트 설정
	app.Get("/api/nodes", nodeController.GetMetricsList)
	app.Get("/api/nodes/:nodeName", nodeController.GetMetricsByNodeName)
	app.Get("/api/nodes/:nodeName/pods", nodeController.GetPodMetricsListByNodeName)

	app.Get("/api/pods", podController.GetMetricsList)
	app.Get("/api/pods/:podName", podController.GetMetricsByPodName)

	// 실행
	err := app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
