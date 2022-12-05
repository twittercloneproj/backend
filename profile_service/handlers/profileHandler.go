package handlers

import (
	"auth_service/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type UsersHandler struct {
	logger *log.Logger
	repo   *data.UserRepo
}

func NewUsersHandler(l *log.Logger, r *data.UserRepo) *UsersHandler {
	return &UsersHandler{l, r}
}

func (p *UsersHandler) GetUserByUsername(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	username := vars["username"]

	patient, err := p.repo.GetOneUser(username)
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if patient == nil {
		http.Error(rw, "Patient with given id not found", http.StatusNotFound)
		p.logger.Printf("Patient with id: '%s' not found", username)
		return
	}

	err = patient.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}
