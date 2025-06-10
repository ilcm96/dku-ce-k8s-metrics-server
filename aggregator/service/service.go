package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/aggregator/db"
	"github.com/ilcm96/dku-ce-k8s-metrics-server/aggregator/kube"
	sharedTypes "github.com/ilcm96/dku-ce-k8s-metrics-server/shared/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

func SaveMetrics() {
	log.Println("SaveMetrics() executed")

	podUIDToDeploymentNameMap, podUIDToNamespaceNameMap, podUIDToPodMap := getResourceInfo()
	collectorIps := getCollectorIps()
	metrics := fetchMetrics(collectorIps)

	ctx := context.Background()

	for _, m := range metrics {
		_, err := db.Pool.Exec(ctx, `
			INSERT INTO node_metrics (
				timestamp,
				node_name,
				cpu_total,
				cpu_busy,
				cpu_count,
				memory_total,
				memory_available,
				memory_used,
				disk_read_bytes,
				disk_write_bytes,
				network_rx_bytes,
				network_tx_bytes
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`,
			m.Timestamp,
			m.NodeMetric.NodeName,
			m.NodeMetric.CPUTotal,
			m.NodeMetric.CPUBusy,
			m.NodeMetric.CPUCount,
			m.NodeMetric.MemoryTotal,
			m.NodeMetric.MemoryAvailable,
			m.NodeMetric.MemoryUsed,
			m.NodeMetric.DiskReadBytes,
			m.NodeMetric.DiskWriteBytes,
			m.NodeMetric.NetworkRxBytes,
			m.NodeMetric.NetworkTxBytes,
		)
		if err != nil {
			log.Println("Failed to insert node metric for node", m.NodeMetric.NodeName, "Error:", err)
		}

		for _, p := range m.PodMetric {
			podName := podUIDToPodMap[types.UID(p.UID)].Name
			deploymentName := podUIDToDeploymentNameMap[types.UID(p.UID)]
			namespaceName := podUIDToNamespaceNameMap[types.UID(p.UID)]
			nodeName := m.NodeMetric.NodeName

			var namespaceParam any = namespaceName
			if namespaceName == "" {
				namespaceParam = nil
			}

			var deploymentParam any = deploymentName
			if deploymentName == "" {
				deploymentParam = nil
			}

			_, err := db.Pool.Exec(ctx, `
				INSERT INTO pod_metrics (
					timestamp,
					pod_name,
					uid,
					cpu_usage_usec,
					memory_usage,
					disk_read_bytes,
					disk_write_bytes,
					network_rx_bytes,
					network_tx_bytes,
					namespace_name,
					deployment_name,
					node_name
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
				)
			`, m.Timestamp,
				podName,
				p.UID,
				p.CPUUsageUsec,
				p.MemoryUsage,
				p.DiskReadBytes,
				p.DiskWriteBytes,
				p.NetworkRxBytes,
				p.NetworkTxBytes,
				namespaceParam,
				deploymentParam,
				nodeName,
			)
			if err != nil {
				log.Println("Failed to insert pod metric for pod UID", p.UID, "Error:", err)
				continue
			}
		}
	}
}

func getResourceInfo() (map[types.UID]string, map[types.UID]string, map[types.UID]*v1.Pod) {
	podUIDToDeploymentNameMap := make(map[types.UID]string)
	podUIDToNamespaceNameMap := make(map[types.UID]string)
	podUIDToPodMap := make(map[types.UID]*v1.Pod)

	pods, _ := kube.PodLister.Pods("").List(labels.Everything())
	replicasets, _ := kube.ReplicaSetLister.ReplicaSets("").List(labels.Everything())

	rsNameToDeploymentName := make(map[string]string)
	for _, rs := range replicasets {
		for _, ownerRef := range rs.OwnerReferences {
			if ownerRef.APIVersion == "apps/v1" && ownerRef.Kind == "Deployment" && ownerRef.Name != "" {
				rsNameToDeploymentName[rs.Name] = ownerRef.Name
				break
			}
		}
	}

	for _, pod := range pods {
		if pod.OwnerReferences != nil {
			found := false
			for _, ownerRef := range pod.OwnerReferences {
				if ownerRef.APIVersion == "apps/v1" && ownerRef.Kind == "ReplicaSet" && ownerRef.Name != "" {
					podUIDToDeploymentNameMap[pod.UID] = rsNameToDeploymentName[ownerRef.Name]
					found = true
					break
				}
			}
			if !found {
				podUIDToDeploymentNameMap[pod.UID] = ""
			}
		} else {
			podUIDToDeploymentNameMap[pod.UID] = ""
		}
		podUIDToNamespaceNameMap[pod.UID] = pod.Namespace
		podUIDToPodMap[pod.UID] = pod
	}

	return podUIDToDeploymentNameMap, podUIDToNamespaceNameMap, podUIDToPodMap
}

func getCollectorIps() []string {
	var ips []string

	selector := labels.SelectorFromSet(labels.Set{"kubernetes.io/service-name": "metrics-collector-headless-svc"})
	endpointSlices, _ := kube.EndpointSliceLister.EndpointSlices("metrics-server-ns").List(selector)
	for _, endpointSlice := range endpointSlices {
		for _, endpoint := range endpointSlice.Endpoints {
			if len(endpoint.Addresses) > 0 {
				ips = append(ips, endpoint.Addresses[0])
			}
		}
	}

	return ips
}

func fetchMetrics(ips []string) []sharedTypes.Metric {
	var metrics []sharedTypes.Metric
	for _, ip := range ips {
		resp, err := http.Get(fmt.Sprintf("http://%s:9000/metrics", ip))
		if err != nil {
			log.Println("Failed to fetch metrics from", ip, "Error:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("Failed to fetch metrics from", ip, "Status Code:", resp.StatusCode)
			continue
		}

		var metric sharedTypes.Metric
		if err := json.NewDecoder(resp.Body).Decode(&metric); err != nil {
			log.Println("Failed to decode metrics from", ip, "Error:", err)
			continue
		}
		metrics = append(metrics, metric)
	}
	return metrics
}
