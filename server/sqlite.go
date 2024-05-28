package server

import (
	"database/sql"
	"fmt"
	"log"
)

// SQLiteController is a controller with sqlite backend.
type SQLiteController struct {
	DataSource string
	pool       *sql.DB
}

func (c *SQLiteController) Initialize() error {
	pool, err := sql.Open("sqlite3", c.DataSource)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}
	defer pool.Close()

	c.pool = pool

	createStmt := `
	create table users (username text not null primary key,
	email text,
	password text,
	habits text
	);
	delete from users;
	`

	_, err = pool.Exec(createStmt)
	if err != nil {
		return fmt.Errorf("err executing %s: %w", createStmt, err)
	}

	return nil
}
