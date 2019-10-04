package relational

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/database/schema"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/logging"

	//looping in postgres driver
	_ "github.com/lib/pq"
)

const (
	clientTimeout = 5 * time.Minute
	maxIdleConns  = 1
)

// Global constants for primary and replica databases.
var (
	Primary *sql.DB
	Replica *sql.DB
)

// Init sets up the database connections
func Init(dbconfig config.DatabaseConfiguration) error {

	var err error

	if Primary, err = connect(dbconfig.Primary); err != nil {
		return err
	}

	if Replica, err = connect(dbconfig.Replica); err != nil {
		return err
	}

	return nil

}

// PreMigrate runs the premigration process for the database.
func PreMigrate(loader loaders.ResourceLoader, dbconfig config.DatabaseConfiguration, schemaFiles []string) error {

	err := Init(dbconfig)

	if err != nil {
		return err
	}

	migrator := schema.DefaultMigrator{
		Datasource: Primary,
		Dialecter:  &schema.PostgresDialect{},
	}
	return migrator.Migrate(loader, schemaFiles, dbconfig.Primary.Schema, nil, nil)

}

// PostMigrate runs the postmigration process for the database.
func PostMigrate(loader loaders.ResourceLoader, dbconfig config.DatabaseConfiguration, schemaFiles []string) error {

	err := Init(dbconfig)

	if err != nil {
		return err
	}

	return nil
	/*
		migrator := schema.DefaultMigrator{}
		return migrator.Migrate(loader, schemaFiles, nil, nil)
	*/
}

// connect opens a configured connection pool to a database.
func connect(cfg config.RelationalDatasource) (*sql.DB, error) {
	logging.Infof("Connecting to database: %v", cfg.Schema)

	connStr := fmt.Sprintf(
		"host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		cfg.Hostname,
		cfg.Portnumber,
		cfg.Schema,
		cfg.Username,
		cfg.Password,
	)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(clientTimeout)
	db.SetMaxIdleConns(maxIdleConns)

	return db, nil
}
