package repository

import (
	"fmt"
	"time"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/jmoiron/sqlx"
)

type PodRepository interface {
	FindAll() ([]*entity.PodMetrics, error)
	FindByPodName(podName string) ([]*entity.PodMetrics, error)
	FindByNodeName(nodeName string) ([]*entity.PodMetrics, error)
	FindByPodNameInTimeWindow(podName string, startTime, endTime time.Time) ([]*entity.PodMetrics, error)
}

type podRepository struct {
	db *sqlx.DB
}

func NewPodRepository(db *sqlx.DB) PodRepository {
	return &podRepository{
		db: db,
	}
}

// FindAll 은 모든 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
func (r *podRepository) FindAll() ([]*entity.PodMetrics, error) {
	query := `
		WITH ranked AS (
			SELECT
				*,
				ROW_NUMBER() OVER (PARTITION BY pod_name ORDER BY timestamp DESC) AS rn
			FROM pod_metrics
		)
		SELECT
			id, timestamp, pod_name, uid, cpu_usage_usec, memory_usage,
			disk_read_bytes, disk_write_bytes, network_rx_bytes, network_tx_bytes,
			namespace_name, deployment_name, node_name
		FROM ranked
		WHERE rn <= 2
		ORDER BY pod_name, timestamp DESC;
	`

	var metrics []*entity.PodMetrics
	err := r.db.Select(&metrics, query)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// FindByPodName 은 주어진 파드명에 대하여 가장 최근의 2개의 메트릭을 조회합니다.
func (r *podRepository) FindByPodName(podName string) ([]*entity.PodMetrics, error) {
	query := `
		SELECT *
		FROM pod_metrics
		WHERE pod_name = $1
		ORDER BY timestamp DESC
		LIMIT 2;
	`

	var metrics []*entity.PodMetrics
	err := r.db.Select(&metrics, query, podName)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// FindByNodeName 은 주어진 노드명을 가진 모든 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
func (r *podRepository) FindByNodeName(nodeName string) ([]*entity.PodMetrics, error) {
	query := `
        WITH ranked AS (
            SELECT
                *,
                ROW_NUMBER() OVER (PARTITION BY pod_name ORDER BY timestamp DESC) AS rn
            FROM pod_metrics
            WHERE node_name = $1
        )
        SELECT
			id, timestamp, pod_name, uid, cpu_usage_usec, memory_usage,
			disk_read_bytes, disk_write_bytes, network_rx_bytes, network_tx_bytes,
			namespace_name, deployment_name, node_name
        FROM ranked
        WHERE rn <= 2
        ORDER BY pod_name, timestamp DESC;
    `

	var metrics []*entity.PodMetrics
	err := r.db.Select(&metrics, query, nodeName)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// FindByPodNameInTimeWindow 는 주어진 파드명과 시간 범위에 대한 메트릭을 조회합니다.
func (r *podRepository) FindByPodNameInTimeWindow(podName string, startTime, endTime time.Time) ([]*entity.PodMetrics, error) {
	fmt.Println("Finding metrics for pod:", podName, "from", startTime, "to", endTime)

	query := `
		SELECT
			id, timestamp, pod_name, uid, cpu_usage_usec, memory_usage,
			disk_read_bytes, disk_write_bytes, network_rx_bytes, network_tx_bytes,
			namespace_name, deployment_name, node_name
		FROM pod_metrics
		WHERE pod_name = $1
		  AND timestamp >= $2
		  AND timestamp <= $3
		ORDER BY timestamp DESC;
	`

	var metrics []*entity.PodMetrics
	err := r.db.Select(&metrics, query, podName, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
