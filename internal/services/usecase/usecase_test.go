package usecase

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/manabie-com/togo/internal/storages"
	"github.com/manabie-com/togo/internal/storages/postgres"
	"os"
	"testing"
)

var (
	todoUseCase ToDoUseCase
	store *postgres.DBPostgres
)

func init () {
	store := &postgres.DBPostgres{
		Host:     os.Getenv("HOST"),
		Port:     5432,
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}
	store.Connect()
	todoUseCase = ToDoUseCase{
		Store: store,
		JWTKey: "wqGyEBBfPK9w3Lxw",
	}
}

func TestGetAuthToken(t *testing.T) {
	userID := "firstUser"
	password := "example"
	_, err := todoUseCase.GetAuthToken(
		context.Background(),
		sql.NullString{
			String: userID,
			Valid: true,
		},
		sql.NullString{
			String: password,
			Valid: true,
		},
	)
	if err != nil {
		t.Errorf("Output expect nil instead of %v", err)
	}
}

func TestListTasks(t *testing.T) {
	userID := "firstUser"
	createDate := "2021-01-10"
	_, err := todoUseCase.ListTasks (
		context.Background(),
		sql.NullString{
			String: userID,
			Valid: true,
		},
		sql.NullString{
			String: createDate,
			Valid: true,
		},
	)
	if err != nil {
		t.Errorf("Output expect nil instead of %v", err)
	}
}

func TestAddTask(t *testing.T) {
	task := &storages.Task{
		Content: "testing",
	}
	userID := "firstUser"
	err := todoUseCase.AddTask(context.Background(),task,userID)
	if err != nil && err.Error() != "You exceeded the limit of tasks per day!" {
		t.Errorf("Output expect nil instead of %v", err)
	}
}

func TestCreateToken(t *testing.T) {
	token := "jshdgajsgd"
	_, ok := todoUseCase.ValidToken(token)
	if ok == true {
		t.Errorf("Output expect false instead of %v", ok)
	}
}