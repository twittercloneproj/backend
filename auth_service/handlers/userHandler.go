package handlers

import (
	"auth_service/data"
	"context"
	"log"
	"net/http"
)

type KeyUser struct{}

type UsersHandler struct {
	logger *log.Logger
	repo   *data.UserRepo
}

// Injecting the logger makes this code much more testable.
func NewUsersHandler(l *log.Logger, r *data.UserRepo) *UsersHandler {
	return &UsersHandler{l, r}
}

func (p *UsersHandler) GetAllUsers(rw http.ResponseWriter, h *http.Request) {
	allUsers, err := p.repo.GetAll()
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		p.logger.Fatal("Database exception: ", err)
	}

	err = allUsers.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *UsersHandler) PostUsers(rw http.ResponseWriter, h *http.Request) {
	usr := h.Context().Value(KeyUser{}).(*data.User)
	p.repo.Post(usr)
	rw.WriteHeader(http.StatusCreated)
}

func (p *UsersHandler) MiddlewarePatientDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		user := &data.User{}
		err := user.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyUser{}, user)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}
