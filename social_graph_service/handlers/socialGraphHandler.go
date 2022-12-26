package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social_graph_service/data"
)

type KeyProduct struct{}

type SocialGraphHandler struct {
	logger *log.Logger
	// NoSQL: injecting movie repository
	repo *data.SocialGraphRepo
}

// Injecting the logger makes this code much more testable.
func NewMoviesHandler(l *log.Logger, r *data.SocialGraphRepo) *SocialGraphHandler {
	return &SocialGraphHandler{l, r}
}

func (m *SocialGraphHandler) CreateUser(rw http.ResponseWriter, h *http.Request) {

	var postUser data.User
	eerr := json.NewDecoder(h.Body).Decode(&postUser)

	if eerr != nil {
		http.Error(rw, "Cannot unmarshal body", 500)
		return
	}
	fmt.Printf(postUser.Username)
	err := m.repo.WritePerson(&postUser)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (m *SocialGraphHandler) Follow(rw http.ResponseWriter, h *http.Request) {

	// #TODO
	//to := mux.Vars(h)["username"]
	//from := jwtTOken // uzimanje iz jwt

	var postUser data.User
	eerr := json.NewDecoder(h.Body).Decode(&postUser)

	if eerr != nil {
		http.Error(rw, "Cannot unmarshal body", 500)
		return
	}
	fmt.Printf(postUser.Username)
	err := m.repo.WritePerson(&postUser)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
