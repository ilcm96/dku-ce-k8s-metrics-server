package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	dsn := createDsn()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Failed to parse DSN:", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 10 * time.Minute
	config.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	Pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to PostgreSQL database successfully")
}

func Migrate() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := Pool.Begin(ctx)
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS node_metrics (
	  id                SERIAL PRIMARY KEY,
	  timestamp         TIMESTAMP NOT NULL,
	  node_name         TEXT      NOT NULL,
	  cpu_total         REAL      NOT NULL,
	  cpu_busy          REAL      NOT NULL,
	  memory_total      BIGINT    NOT NULL,
	  memory_available  BIGINT,
	  memory_used       BIGINT    NOT NULL,
	  disk_read_bytes   BIGINT    NOT NULL,
	  disk_write_bytes  BIGINT    NOT NULL,
	  network_rx_bytes  BIGINT    NOT NULL,
	  network_tx_bytes  BIGINT    NOT NULL
	);

	CREATE TABLE IF NOT EXISTS pod_metrics (
	  id                SERIAL PRIMARY KEY,
	  timestamp         TIMESTAMP NOT NULL,
	  pod_name          TEXT      NOT NULL,
	  uid               TEXT      NOT NULL,
	  cpu_usage_usec    BIGINT    NOT NULL,
	  memory_usage      BIGINT    NOT NULL,
	  disk_read_bytes   BIGINT    NOT NULL,
	  disk_write_bytes  BIGINT    NOT NULL,
	  network_rx_bytes  BIGINT    NOT NULL,
	  network_tx_bytes  BIGINT    NOT NULL,
	  namespace_name    TEXT      NOT NULL,
	  deployment_name   TEXT,
	  node_name         TEXT      NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_pod_metrics_namespace ON pod_metrics (namespace_name);
	CREATE INDEX IF NOT EXISTS idx_pod_metrics_deployment ON pod_metrics (deployment_name);
	`

	if _, err := tx.Exec(ctx, schema); err != nil {
		log.Fatal("Failed to execute migration schema:", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}

	log.Println("Database migration completed successfully")
}

func createDsn() string {
	return "postgres://" + os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
}
