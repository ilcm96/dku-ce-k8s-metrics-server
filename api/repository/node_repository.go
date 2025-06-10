package repository

import (
	"time"

	"github.com/ilcm96/dku-ce-k8s-metrics-server/api/entity"
	"github.com/jmoiron/sqlx"
)

type NodeRepository interface {
	FindAll() ([]*entity.NodeMetrics, error)
	FindByNodeName(nodeName string) ([]*entity.NodeMetrics, error)
	FindByNodeNameInTimeWindow(nodeName string, startTime, endTime time.Time) ([]*entity.NodeMetrics, error)
}

type nodeRepository struct {
	db *sqlx.DB
}

func NewNodeRepository(db *sqlx.DB) NodeRepository {
	return &nodeRepository{
		db: db,
	}
}

// FindAll 는 모든 노드들에 대해 가장 최근의 2개의 메트릭을 조회합니다.
func (r *nodeRepository) FindAll() ([]*entity.NodeMetrics, error) {
	query := `
		WITH ranked AS (
		    SELECT
		        *,
		        ROW_NUMBER() OVER (PARTITION BY node_name ORDER BY timestamp DESC) AS rn
		    FROM node_metrics
		)
		SELECT
			id, timestamp, node_name, cpu_total, cpu_busy, cpu_count,
			memory_total, memory_available, memory_used,
			disk_read_bytes, disk_write_bytes, network_rx_bytes, network_tx_bytes
		FROM ranked
		WHERE rn <= 2
		ORDER BY node_name, timestamp DESC;
	`

	var metrics []*entity.NodeMetrics
	err := r.db.Select(&metrics, query)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// FindByNodeName 은 주어진 노드명에 대하여 가장 최근의 2개의 메트릭을 조회합니다.
func (r *nodeRepository) FindByNodeName(nodeName string) ([]*entity.NodeMetrics, error) {
	query := `
		SELECT *
		FROM node_metrics
		WHERE node_name = $1
		ORDER BY timestamp DESC
		LIMIT 2;
	`

	var metrics []*entity.NodeMetrics
	err := r.db.Select(&metrics, query, nodeName)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// FindByNodeNameInTimeWindow 는 주어진 노드명과 시간 범위에 대한 메트릭을 조회합니다.
func (r *nodeRepository) FindByNodeNameInTimeWindow(nodeName string, startTime, endTime time.Time) ([]*entity.NodeMetrics, error) {
	query := `
		SELECT
			id, timestamp, node_name, cpu_total, cpu_busy, cpu_count,
			memory_total, memory_available, memory_used,
			disk_read_bytes, disk_write_bytes, network_rx_bytes, network_tx_bytes
		FROM node_metrics
		WHERE node_name = $1
		  AND timestamp >= $2
		  AND timestamp <= $3
		ORDER BY timestamp DESC;
	`

	var metrics []*entity.NodeMetrics
	err := r.db.Select(&metrics, query, nodeName, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
