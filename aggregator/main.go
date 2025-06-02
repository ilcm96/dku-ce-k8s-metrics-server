package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-co-op/gocron/v2"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/aggregator/db"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/aggregator/kube"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/aggregator/service"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		log.Println("Loading environment variables from .env file")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to load .env file:", err)
		} else {
			log.Println("Environment variables loaded successfully from .env file")
		}

	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("Failed to create scheduler:", err)
	}
	defer s.Shutdown()
	log.Println("Scheduler created successfully")

	db.Connect()
	defer db.Pool.Close()
	db.Migrate()

	kube.InitKubeConfig()
	kube.InitClientset()
	kube.ListCollectorIP()

	stopCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		close(stopCh)
	}()
	kube.InitLister(stopCh)

	job, err := s.NewJob(
		gocron.CronJob("*/1 * * * *", false), // Every minute
		gocron.NewTask(service.SaveMetrics),
	)
	if err != nil {
		log.Fatal("Failed to create job:", err)
	}
	log.Println("Job created successfully:", job.ID())
	s.Start()

	<-stopCh
	log.Println("Shutting down gracefully.")
}
