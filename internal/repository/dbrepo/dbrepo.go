package dbrepo

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresRepo struct {
	DB *sql.DB
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// `String` returns the Postgres DB connection string.
func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

// `Open`` opens a SQL connection with the provided Postgres DB. Callers should
// ensure that the connection is eventually closed via the `db.Close()` method.
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	// Checks the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to the database!")
	return db, nil
}

// `Migrate` runs the goose migration to setup the database
func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

// `MigrateFS` sets the base file system, runs the migration, and sets the file
// system back to `nil`
func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	// Sets the file system and immediately sets it back to `nil`
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}
