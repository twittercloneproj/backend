package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"social_graph_service/data"
	"strings"
)

type KeyProduct struct{}

type SocialGraphHandler struct {
	logger *log.Logger
	// NoSQL: injecting movie repository
	repo *data.SocialGraphRepo
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

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

	vars := mux.Vars(h)
	to := vars["username"]

	bearer := h.Header.Get("Authorization")
	bearerToken := strings.Split(bearer, "Bearer ")
	tokenString := bearerToken[1]

	token, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Cannot parse token", 403)
		return
	}

	claims := GetMapClaims(token.Bytes())
	from := claims["username"]

	errr := m.repo.FollowPerson(from, to)
	if errr != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func GetMapClaims(tokenBytes []byte) map[string]string {
	var claims map[string]string

	err := jwt.ParseClaims(tokenBytes, verifier, &claims)
	if err != nil {
		log.Println(err)
	}

	return claims
}
