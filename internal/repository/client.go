package repository

import "database/sql"

// это для чего?
type Client struct {
	DB *sql.DB
}
