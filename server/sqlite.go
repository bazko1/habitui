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
}

func NewSQLiteController(dataSource string) *SQLiteController {
	return &SQLiteController{DataSource: dataSource, pool: nil}
}

func (c *SQLiteController) Initialize() error {
	pool, err := sql.Open("sqlite3", c.DataSource)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}

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

	return nil
}

func (c SQLiteController) GetUserByName(name string) (UserModel, error) {
	var userModel UserModel

	err := c.pool.QueryRow("select email, password, habits from users where name = '?'", name).Scan(userModel)
	if err != nil {
		return UserModel{}, fmt.Errorf("failed prepare select statement %w", err)
	}

	return userModel, nil
}

func (c SQLiteController) CreateNewUser(user UserModel) (UserModel, error) {
	insertStmt, err := c.pool.Prepare("insert into users(username, email, password, habits) values(?, ?, ?, ?)")
	if err != nil {
		return UserModel{}, fmt.Errorf("failed prepare insert statement %w", err)
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec(user.Username, user.Email, user.Password, user.Habits)
	if err != nil {
		return UserModel{}, fmt.Errorf("failed execute insert statement %w", err)
	}

	return UserModel{}, nil
}

func (c SQLiteController) UpdateUserHabits(user UserModel, habits habit.TaskList) error {
	_, err := c.pool.Exec(`update users set habits = ? where username == '?'`, habits, &user.Username)
	if err != nil {
		return fmt.Errorf("failed execute update statement %w", err)
	}

	return nil
}

func (c SQLiteController) GetUserHabits(user UserModel) (habit.TaskList, error) {
	taskList := habit.TaskList{}

	err := c.pool.QueryRow("select habits from users where name = '?'", user.Username).Scan(&taskList)
	if err != nil {
		return habit.TaskList{}, fmt.Errorf("failed execute select statement %w", err)
	}

	return taskList, nil
}

func (c SQLiteController) IsValid(user UserModel) bool {
	var exists bool

	err := c.pool.QueryRow("select exists(select 1 from users where username = '?' and password = '?')",
		user.Username, user.Password).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}
