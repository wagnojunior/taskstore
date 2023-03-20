package taskstore

import (
	"database/sql"
	"fmt"
)

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
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s"+
		"sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
		cfg.SSLMode)
}

// `Open`` opens a SQL connection with the provided Postgres DB. Callers should
// ensure that the connection is eventually closed via the `db.Close()` method.
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return db, nil
}
