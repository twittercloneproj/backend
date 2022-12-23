package handlers

import (
	"auth_service/data"
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

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

type KeyProduct struct{}

func (p *UsersHandler) ChangePrivacy(rw http.ResponseWriter, h *http.Request) {

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
	username := claims["username"]

	// kopirati negde ako zatreba #DECODE #UNMARSHALL
	var updatePrivacy data.UpdatePrivacy
	eerr := json.NewDecoder(h.Body).Decode(&updatePrivacy)

	if eerr != nil {
		fmt.Println(err)
		http.Error(rw, "Cannot unmarshal body", 500)
		return
	}
	p.repo.Update(username, updatePrivacy.Privacy)
	rw.WriteHeader(http.StatusOK)
}

func GetMapClaims(tokenBytes []byte) map[string]string {
	var claims map[string]string

	err := jwt.ParseClaims(tokenBytes, verifier, &claims)
	if err != nil {
		log.Println(err)
	}

	return claims
}
