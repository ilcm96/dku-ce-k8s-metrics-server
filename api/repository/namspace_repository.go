package repository

import (
	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/jmoiron/sqlx"
)

type NamespaceRepository interface {
	FindAll() ([]*entity.PodMetrics, error)
	FindByNamespaceName(namespaceName string) ([]*entity.PodMetrics, error)
}

type namespaceRepository struct {
	db *sqlx.DB
}

func NewNamespaceRepository(db *sqlx.DB) NamespaceRepository {
	return &namespaceRepository{
		db: db,
	}
}

// FindAll 는 모든 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
func (r *namespaceRepository) FindAll() ([]*entity.PodMetrics, error) {
	query := `
		WITH ranked AS (
			SELECT
				*,
				ROW_NUMBER() OVER (PARTITION BY pod_name ORDER BY timestamp DESC) AS rn
			FROM pod_metrics
			WHERE namespace_name IS NOT NULL
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

// FindByNamespaceName 는 특정 네임스페이스의 파드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
func (r *namespaceRepository) FindByNamespaceName(namespaceName string) ([]*entity.PodMetrics, error) {
	query := `
		WITH ranked AS (
			SELECT
				*,
				ROW_NUMBER() OVER (PARTITION BY pod_name ORDER BY timestamp DESC) AS rn
			FROM pod_metrics
			WHERE namespace_name = $1
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
	err := r.db.Select(&metrics, query, namespaceName)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
