package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bazko1/habitui/habit"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteController is a controller with sqlite backend.
type SQLiteController struct {
	DataSource string
	pool       *sql.DB
	insertStmt *sql.Stmt
}

func NewSQLiteController(dataSource string) *SQLiteController {
	return &SQLiteController{DataSource: dataSource, pool: nil}
}

func (c *SQLiteController) Initialize() error {
	pool, err := sql.Open("sqlite3", c.DataSource)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}
	defer pool.Close()

	c.pool = pool

	// TODO: Not sure if i want jsonb or json for habits
	createStmt := `
	create table if not exists users (username text not null primary key,
	email text,
	password text,
	habits jsonb
	);
	delete from users;
	`

	_, err = pool.Exec(createStmt)
	if err != nil {
		return fmt.Errorf("err executing %s: %w", createStmt, err)
	}

	// TODO the statement will need to be closed thus not sure if this is correct way to do that
	c.insertStmt, err = c.pool.Prepare("insert into users(username, email, password, habits) values(?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed prepare insert statement %w", err)
	}

	return nil
}

func (c SQLiteController) GetUserByName(name string) (UserModel, error) {
	stmt, err := c.pool.Prepare("select email, password, habits from users where name = ?")
	if err != nil {
		// TODO move to init or this method will need to also return error
		// return fmt.Errorf("failed prepare select statement %w", err)
		return UserModel{}, fmt.Errorf("failed prepare select statement %w", err)
	}
	defer stmt.Close()

	var userModel UserModel
	if err := stmt.QueryRow().Scan(userModel); err != nil {
		return UserModel{}, ErrUsernameDoesNotExist
	}

	return userModel, nil
}

func (c SQLiteController) CreateNewUser(user UserModel) (UserModel, error) {
	return UserModel{}, nil
}

func (c SQLiteController) UpdateUserHabits(user UserModel, habits habit.TaskList) error {
	return nil
}

func (c SQLiteController) GetUserHabits(user UserModel) (habit.TaskList, error) {
	return habit.TaskList{}, nil
}

func (c SQLiteController) IsValid(user UserModel) bool {
	return false
}
