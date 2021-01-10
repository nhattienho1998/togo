package main

import (
	"github.com/manabie-com/togo/internal/services/transport"
	"github.com/manabie-com/togo/internal/services/usecase"
	"github.com/manabie-com/togo/internal/storages/postgres"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port,_ := strconv.Atoi(os.Getenv("PORT"))
	store := &postgres.DBPostgres{
		Host:     os.Getenv("HOST"),
		Port:     port,
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}
	store.Connect()
	defer store.Close()

	todoUseCase := usecase.ToDoUseCase{
		Store: store,
		JWTKey: "wqGyEBBfPK9w3Lxw",
	}

	http.ListenAndServe(":5050", &transport.ToDoTransport{
		ToDoUseCase: todoUseCase,
	})
}
