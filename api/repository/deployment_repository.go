package repository

import (
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/jmoiron/sqlx"
)

type DeploymentRepository interface {
	FindByNamespaceName(namespaceName string) ([]*entity.PodMetrics, error)
	FindByDeploymentName(namespaceName, deploymentName string) ([]*entity.PodMetrics, error)
}

type deploymentRepository struct {
	db *sqlx.DB
}

func NewDeploymentRepository(db *sqlx.DB) DeploymentRepository {
	return &deploymentRepository{
		db: db,
	}
}

// FindByNamespaceName 는 특정 네임스페이스의 디플로이먼트들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
func (r *deploymentRepository) FindByNamespaceName(namespaceName string) ([]*entity.PodMetrics, error) {
	query := `
		WITH ranked AS (
			SELECT
				*,
				ROW_NUMBER() OVER (PARTITION BY pod_name ORDER BY timestamp DESC) AS rn
			FROM pod_metrics
			WHERE namespace_name = $1 AND deployment_name IS NOT NULL
		)
		SELECT
			id, timestamp, pod_name, uid, cpu_usage_usec, memory_usage,
			disk_read_bytes, disk_write_bytes, network_rx_bytes, network_tx_bytes,
			namespace_name, deployment_name, node_name
		FROM ranked
		WHERE rn <= 2
		ORDER BY deployment_name, pod_name, timestamp DESC;
	`

	var metrics []*entity.PodMetrics
	err := r.db.Select(&metrics, query, namespaceName)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// FindByDeploymentName 는 특정 디플로이먼트의 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
func (r *deploymentRepository) FindByDeploymentName(namespaceName, deploymentName string) ([]*entity.PodMetrics, error) {
	query := `
		WITH ranked AS (
			SELECT
				*,
				ROW_NUMBER() OVER (PARTITION BY pod_name ORDER BY timestamp DESC) AS rn
			FROM pod_metrics
			WHERE namespace_name = $1 AND deployment_name = $2
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
	err := r.db.Select(&metrics, query, namespaceName, deploymentName)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
