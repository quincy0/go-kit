package postgres

import (
	"go-kit/core/stores/sqlx"
	// imports the driver, don't remove this comment, golint requires.
	_ "github.com/lib/pq"
)

const postgresDriverName = "postgres"

// New returns a postgres connection.
func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(postgresDriverName, datasource, opts...)
}
