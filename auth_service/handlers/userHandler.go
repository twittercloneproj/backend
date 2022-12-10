package handlers

import (
	"auth_service/application"
	"auth_service/domain"
	"auth_service/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type KeyUser struct{}

type UsersHandler struct {
	logger *log.Logger
	//repo    *data.UserRepo
	service *application.AuthService
	store   *store.AuthMongoDBStore
}

func NewAuthHandler(service *application.AuthService) *UsersHandler {
	return &UsersHandler{
		service: service,
	}
}

func (handler *UsersHandler) Init(router *mux.Router) {

	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/register", handler.PostUsers).Methods("POST")
	router.HandleFunc("/all", handler.GetAllUsers).Methods("GET")
	http.Handle("/", router)
	//log.Fatal(http.ListenAndServe("/", router))

}

func (p *UsersHandler) GetAllUsers(rw http.ResponseWriter, h *http.Request) {
	allUsers, err := p.service.GetAll()
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		p.logger.Fatal("Database exception: ", err)
	}
	jsonResponse(allUsers, rw)

}

func (p *UsersHandler) PostUsers(rw http.ResponseWriter, h *http.Request) {
	var user domain.User
	err := json.NewDecoder(h.Body).Decode(&user)

	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	token, err := p.service.Post(&user)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	jsonResponse(token, rw)
}

func (handler *UsersHandler) Login(writer http.ResponseWriter, req *http.Request) {
	var request domain.User
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := handler.service.LoginHelp(&request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Write([]byte(token))
}

func jsonResponse(object interface{}, w http.ResponseWriter) {
	resp, err := json.Marshal(object)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
