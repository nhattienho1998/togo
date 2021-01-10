package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/manabie-com/togo/internal/storages"
	"github.com/manabie-com/togo/internal/storages/postgres"
	"log"
	"time"
)


type ToDoUseCase struct {
	JWTKey string
	Store  *postgres.DBPostgres
}

func (s *ToDoUseCase) GetAuthToken(ctx context.Context, userID, password sql.NullString) (*string, error) {
	if !s.Store.ValidateUser(ctx, userID, password) {
		return nil, errors.New("incorrect user_id/pwd")
	}

	token, err := s.createToken(userID.String)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *ToDoUseCase) ListTasks(ctx context.Context, userID, createDate sql.NullString)([]*storages.Task, error) {
	tasks, err := s.Store.RetrieveTasks(
		ctx,
		userID,
		createDate,
	)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *ToDoUseCase) AddTask(ctx context.Context, t *storages.Task, userID string)  error {

	now := time.Now()
	t.ID = uuid.New().String()
	t.UserID = userID
	t.CreatedDate = now.Format("2006-01-02")

	err := s.validateAddTasks(ctx, userID, t.CreatedDate)
	if err != nil {
		return err
	}

	err = s.Store.AddTask(ctx, t)
	if err != nil {
		return err
	}

 	return nil
}


func (s *ToDoUseCase) createToken(id string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(s.JWTKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *ToDoUseCase) ValidToken(token string) (*string, bool) {

	claims := make(jwt.MapClaims)
	t, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(s.JWTKey), nil
	})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	if !t.Valid {
		return nil, false
	}

	id, ok := claims["user_id"].(string)
	if !ok {
		return nil, false
	}

	return &id, true
}


func (s *ToDoUseCase) validateAddTasks(ctx context.Context, userID, createDate string) error {
	numberOfTasks, _ := s.Store.CountTasks(
		ctx,
		sql.NullString{
			String: userID,
			Valid:  true,
		},
		sql.NullString{
			String: createDate,
			Valid:  true,
		},
	)
	maxToDo, err := s.Store.GetMaxToDo(
		ctx,
		sql.NullString{
			String: userID,
			Valid:  true,
		},
	)

	if err != nil {
		return err
	}

	if numberOfTasks >= maxToDo {
		return errors.New("You exceeded the limit of tasks per day!")
	}

	return nil
}