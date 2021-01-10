package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/manabie-com/togo/internal/services/usecase"
	"github.com/manabie-com/togo/internal/storages"
	"log"
	"net/http"
)

type ToDoTransport struct {
	ToDoUseCase usecase.ToDoUseCase
}

func (s *ToDoTransport) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Headers", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "*")

	if req.Method == http.MethodOptions {
		resp.WriteHeader(http.StatusOK)
		return
	}

	switch req.URL.Path {
	case "/login":
		s.getAuthToken(resp, req)
		return
	case "/tasks":
		token := req.Header.Get("Authorization")
		var ok bool
		id, ok := s.ToDoUseCase.ValidToken(token)
		if !ok {
			resp.WriteHeader(http.StatusUnauthorized)
			return
		}
		req = req.WithContext(context.WithValue(req.Context(), userAuthKey(0), *id))
		switch req.Method {
		case http.MethodGet:
			s.listTasks(resp, req)
		case http.MethodPost:
			s.addTask(resp, req)
		}
		return
	}
}

func (s *ToDoTransport) getAuthToken(resp http.ResponseWriter, req *http.Request) {
	id := value(req, "user_id")
	password := value(req, "password")
	token, err := s.ToDoUseCase.GetAuthToken(req.Context(), id, password)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	resp.Header().Set("Content-Type", "application/json")

	json.NewEncoder(resp).Encode(map[string]string{
		"data": *token,
	})
}

func (s *ToDoTransport) listTasks(resp http.ResponseWriter, req *http.Request) {
	id, _ := userIDFromCtx(req.Context())
	tasks, err := s.ToDoUseCase.ListTasks(
		req.Context(),
		sql.NullString{
			String: id,
			Valid:  true,
		},
		value(req, "created_date"),
	)

	resp.Header().Set("Content-Type", "application/json")

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(resp).Encode(map[string][]*storages.Task{
		"data": tasks,
	})
}

func (s *ToDoTransport) addTask(resp http.ResponseWriter, req *http.Request) {
	t := &storages.Task{}
	err := json.NewDecoder(req.Body).Decode(t)
	defer req.Body.Close()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, _ := userIDFromCtx(req.Context())

	resp.Header().Set("Content-Type", "application/json")

	err = s.ToDoUseCase.AddTask(req.Context(), t,userID)
	if err != nil {
		resp.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(resp).Encode(map[string]*storages.Task{
		"data": t,
	})
}

func value(req *http.Request, p string) sql.NullString {
	return sql.NullString{
		String: req.FormValue(p),
		Valid:  true,
	}
}



type userAuthKey int8

func userIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(userAuthKey(0))
	id, ok := v.(string)
	return id, ok
}
