package platform

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTask001PostgresSchemasAndUserIsolation(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "crm_system",
			"POSTGRES_USER":     "crm_admin",
			"POSTGRES_PASSWORD": "crm_admin_dev_password",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start postgres testcontainer: %v", err)
	}
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Fatalf("terminate postgres testcontainer: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("container mapped port: %v", err)
	}
	adminDSN := fmt.Sprintf("postgres://crm_admin:crm_admin_dev_password@%s:%s/crm_system?sslmode=disable", host, port.Port())
	adminDB := openDB(t, adminDSN)
	defer adminDB.Close()

	applyInitialMigrations(t, adminDB)

	var schemaCount int
	if err := adminDB.QueryRow(`
		SELECT count(*)
		FROM information_schema.schemata
		WHERE schema_name IN (
			'identity_authz','lead','account','opportunity','commercial',
			'work','audit_history','reporting','import_export'
		)
	`).Scan(&schemaCount); err != nil {
		t.Fatalf("count schemas: %v", err)
	}
	if schemaCount != 9 {
		t.Fatalf("expected 9 service schemas, got %d", schemaCount)
	}

	leadDSN := fmt.Sprintf("postgres://crm_lead_user:crm_lead_dev_password@%s:%s/crm_system?sslmode=disable&search_path=lead", host, port.Port())
	leadDB := openDB(t, leadDSN)
	defer leadDB.Close()
	if _, err := leadDB.Exec(`CREATE TABLE lead.task001_owned_probe (id integer PRIMARY KEY)`); err != nil {
		t.Fatalf("lead user should create in own schema: %v", err)
	}
	if _, err := leadDB.Exec(`CREATE TABLE identity_authz.task001_forbidden_probe (id integer PRIMARY KEY)`); err == nil {
		t.Fatal("lead user unexpectedly created a table in identity_authz schema")
	}
}

func applyInitialMigrations(t *testing.T, db *sql.DB) {
	t.Helper()
	services := []string{
		"identity-authz",
		"lead",
		"account",
		"opportunity",
		"commercial",
		"work",
		"audit-history",
		"reporting",
		"import-export",
	}
	for _, service := range services {
		path := filepath.Join("..", "..", "..", "..", "services", service, "migrations", "0001_init_schema.up.sql")
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read migration %s: %v", service, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", service, err)
		}
	}
}

func openDB(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping db: %v", err)
	}
	return db
}
