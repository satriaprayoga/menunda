package menunda

import "database/sql"

type config struct {
	port string
}

type Database struct {
	DbType string
	Pool   *sql.DB
}
