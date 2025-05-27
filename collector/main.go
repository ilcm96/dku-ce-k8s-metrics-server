package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"log/slog"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/collector/node"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/collector/pod"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/shared/types"
)

type Metric struct {
	Timestamp  time.Time         `json:"timestamp"`
	NodeMetric types.NodeMetric  `json:"nodeMetric"`
	PodMetric  []types.PodMetric `json:"podMetric"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(log.Writer(), nil)))
	http.HandleFunc("/metrics", loggingMiddleware(collect))
	http.ListenAndServe(":9000", nil)
}

func collect(w http.ResponseWriter, r *http.Request) {
	nodeMetric, err := node.CollectNodeMetric()
	if err != nil {
		nodeMetric = types.NodeMetric{}
	}

	podMetrics, err := pod.CollectPodMetrics()
	if err != nil {
		podMetrics = []types.PodMetric{}
	}

	metric := Metric{
		Timestamp:  time.Now(),
		NodeMetric: nodeMetric,
		PodMetric:  podMetrics,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).String()
		slog.Info("", "method", r.Method, "path", r.URL.Path, "duration", duration)
	}
}
