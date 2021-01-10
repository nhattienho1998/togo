package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/manabie-com/togo/internal/storages"
	"log"
)

type DBPostgres struct {
	Db *sql.DB
	Host string
	Port int
	Username string
	Password string
	Dbname string
}

func (d *DBPostgres) Connect(){
	dataSource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.Username, d.Password, d.Dbname)
	var err error
	d.Db, err = sql.Open("postgres", dataSource)
	if err != nil {
		log.Fatalf("error in connect database: %v\n", err)
	}

	err = d.Db.Ping()
	if err != nil {
		log.Fatalf("error in connect database: %v\n", err)
	}
}

func (d *DBPostgres) Close(){
	d.Db.Close()
}


func (d *DBPostgres) RetrieveTasks(ctx context.Context, userID, createdDate sql.NullString) ([]*storages.Task, error) {
	stmt := `SELECT id, content, user_id, created_date FROM tasks WHERE user_id = $1 AND created_date = $2`
	rows, err := d.Db.QueryContext(ctx, stmt, userID, createdDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*storages.Task
	for rows.Next() {
		t := &storages.Task{}
		err := rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// AddTask adds a new task to DB
func (d *DBPostgres) AddTask(ctx context.Context, t *storages.Task) error {
	stmt := `INSERT INTO tasks (id, content, user_id, created_date) VALUES ($1, $2, $3, $4)`
	_, err := d.Db.ExecContext(ctx, stmt, &t.ID, &t.Content, &t.UserID, &t.CreatedDate)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUser returns tasks if match userID AND password
func (d *DBPostgres) ValidateUser(ctx context.Context, userID, pwd sql.NullString) bool {
	stmt := `SELECT id FROM users WHERE id = $1 AND password = $2`
	row := d.Db.QueryRow(stmt, userID, pwd)
	u := &storages.User{}
	err := row.Scan(&u.ID)
	if err != nil {
		return false
	}

	return true
}

// GetMaxToDo returns max to do task per day if match userID
func (d *DBPostgres) GetMaxToDo(ctx context.Context, userID sql.NullString) (uint, error) {
	stmt := `SELECT max_todo FROM users WHERE id = $1`
	row := d.Db.QueryRowContext(ctx, stmt,userID)
	u := &storages.User{}
	err := row.Scan(&u.MaxTodo)
	if err != nil {
		return 0, err
	}
	return u.MaxTodo, nil
}

// CountTasks return number of tasks of an user if match create_date
func (d *DBPostgres) CountTasks(ctx context.Context, userID, createDate sql.NullString)(uint, error){
	stmt := `SELECT COUNT(id) FROM tasks WHERE user_id = $1 AND created_date = $2`
	rows, err := d.Db.QueryContext(ctx, stmt, userID, createDate)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var numberOfTasks uint
	rows.Next()
	err = rows.Scan(&numberOfTasks)
	if err != nil {
		return 0, err
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	}
	return numberOfTasks, nil
}