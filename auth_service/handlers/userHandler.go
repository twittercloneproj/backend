package handlers

import (
	social_graph "auth_service/client/social-graph"
	"auth_service/data"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

type KeyUser struct{}

type UsersHandler struct {
	logger      *log.Logger
	repo        *data.UserRepo
	socialGraph social_graph.Client
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

// Injecting the logger makes this code much more testable.
func NewUsersHandler(l *log.Logger, r *data.UserRepo, socialGraph social_graph.Client) *UsersHandler {
	return &UsersHandler{l, r, socialGraph}
}

func sendMailSimple(subject string, body string, to []string) {
	auth := smtp.PlainAuth(
		"",
		"oliver.kojic22@gmail.com",
		"tdejbdyydokiprsz",
		"smtp.gmail.com",
	)

	msg := "Subject: " + subject + "\n" + body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"oliver.kojic22@gmail.com",
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
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
	err := p.repo.Post(usr)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	err = p.socialGraph.CreateUser(usr)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
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
		timestamp := time.Now().Add(time.Hour * 1).Format("02-Jan-2006 15:04:05")
		service.logger.Printf("Login failed! Date and Time: %v, Username: %v tried to login", timestamp, credential.Username)
		fmt.Println(err)
		return "", err
	}

	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password))
	if passError != nil {
		timestamp := time.Now().Add(time.Hour * 1).Format("02-Jan-2006 15:04:05")
		service.logger.Printf("Login failed! Date and Time: %v, Username: %v tried to login with wrong password", timestamp, credential.Username)
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

	timestamp := time.Now().Add(time.Hour * 1).Format("02-Jan-2006 15:04:05")
	service.logger.Printf("Login Successful! Date and Time: %v, Username: %v, Password: %v", timestamp, user.Username, user.Password)

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
