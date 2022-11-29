package handlers

import (
	"auth_service/data"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type KeyUser struct{}

type UsersHandler struct {
	logger *log.Logger
	repo   *data.UserRepo
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

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

func (handler *UsersHandler) Login(writer http.ResponseWriter, req *http.Request) {
	var request options.Credential
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := handler.LoginHelp(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(token, writer)
}

func (service *UsersHandler) LoginHelp(credential options.Credential) (string, error) {
	user, err := service.repo.GetOneUser(credential.Username)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password))
	if passError != nil {
		fmt.Println(passError)
		return "", err
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &data.Claims{
		ID:       user.ID,
		Username: user.Username, //menjanje za userID
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//service.GetID(service.GetClaims(tokenString))

	return tokenString, nil
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
