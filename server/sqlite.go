package server

import (
	"database/sql"
	"errors"
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

	createStmt := `
	create table if not exists users (username text not null primary key,
	email text,
	password text,
	habits jsonb  
	);
	`

	_, err = pool.Exec(createStmt)
	if err != nil {
		return fmt.Errorf("Initialize err executing create statement: %w", err)
	}

	return nil
}

func (c SQLiteController) GetUserByName(name string) (UserModel, error) {
	var user UserModel

	err := c.pool.QueryRow("select username, email, password, habits from users where username = ?",
		name).Scan(&user.Username, &user.Email, &user.Password, &user.Habits)

	if errors.Is(err, sql.ErrNoRows) {
		return UserModel{}, ErrUsernameDoesNotExist
	}

	if err != nil {
		return UserModel{}, fmt.Errorf("GetUserByName failed to execute select statement %w", err)
	}

	return user, nil
}

func (c SQLiteController) CreateNewUser(user UserModel) (UserModel, error) {
	res, err := c.pool.Exec("insert or ignore into users(username, email, password, habits) values(?, ?, ?, ?)",
		user.Username, user.Email, user.Password, user.Habits)
	if err != nil {
		return UserModel{}, fmt.Errorf("CreateNewUser failed to execute insert statement %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return UserModel{}, fmt.Errorf("CreateNewUser failed to check affected rows %w", err)
	}

	if affected == 0 {
		return UserModel{}, ErrUsernameAlreadyExists
	}

	return UserModel{}, nil
}

func (c SQLiteController) UpdateUserHabits(user UserModel, habits habit.TaskList) error {
	_, err := c.pool.Exec(`update users set habits == ? where username == ?`, habits, &user.Username)
	if err != nil {
		return fmt.Errorf("UpdateUserHabits failed to execute update statement %w", err)
	}

	return nil
}

func (c SQLiteController) GetUserHabits(user UserModel) (habit.TaskList, error) {
	taskList := habit.TaskList{}

	err := c.pool.QueryRow("select habits from users where username == ?", user.Username).Scan(&taskList)
	if err != nil {
		return habit.TaskList{}, fmt.Errorf("GetUserHabits failed to execute select statement %w", err)
	}

	if taskList == nil {
		taskList = habit.TaskList{}
	}

	return taskList, nil
}

func (c SQLiteController) IsValid(user UserModel) (bool, error) {
	var exists int

	err := c.pool.QueryRow("select exists(select 1 from users where username = ? and password = ?)",
		user.Username, user.Password).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("IsValid failed to execute select statement %w", err)
	}

	return exists == 1, nil
}

func (c SQLiteController) Finalize() error {
	if err := c.pool.Close(); err != nil {
		return fmt.Errorf("Finalize sqlite db close error: %w", err)
	}

	return nil
}
