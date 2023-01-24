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

	user, err := m.repo.GetUser(to)
	if err != nil {
		http.Error(rw, "Cannot check user privacy, try again later", 500)
		return
	}

	var dberr error
	if user.Privacy == "Private" {
		exists, _ := m.repo.CheckIfRelationshipExists(from, to, "REQUEST")
		if !exists {
			dberr = m.repo.FollowPerson(from, to, "REQUEST")
		}
	} else {
		exists, _ := m.repo.CheckIfRelationshipExists(from, to, "FOLLOW")
		if !exists {
			dberr = m.repo.FollowPerson(from, to, "FOLLOW")
		}
	}
	if dberr != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (m *SocialGraphHandler) CheckFollow(rw http.ResponseWriter, h *http.Request) {

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

	follow, dberr := m.repo.CheckIfRelationshipExists(from, to, "FOLLOW")

	if dberr != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if follow {
		rw.WriteHeader(http.StatusOK)
		return
	} else {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

}

func (m *SocialGraphHandler) CanAccessTweet(rw http.ResponseWriter, h *http.Request) {

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

	user, err := m.repo.GetUser(to)

	if err != nil {
		http.Error(rw, "user doesnt exists", 500)
		return
	}

	if user.Privacy == "Public" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	follow, dberr := m.repo.CheckIfRelationshipExists(from, to, "FOLLOW")

	if dberr != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if follow {
		rw.WriteHeader(http.StatusOK)
		return
	} else {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

}

func (m *SocialGraphHandler) RemoveFollow(rw http.ResponseWriter, h *http.Request) {

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

	dberr := m.repo.RemoveFollow(from, to, "FOLLOW")

	if dberr != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (m *SocialGraphHandler) AcceptRejectRequest(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	from := vars["username"]

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
	authUsername := claims["username"]

	var approved data.ApproveRequest
	err = json.NewDecoder(h.Body).Decode(&approved)
	if err != nil {
		http.Error(rw, "Invalid body", 500)
		return
	}

	dberr := m.repo.RemoveFollow(from, authUsername, "REQUEST")
	exists, dberr := m.repo.CheckIfRelationshipExists(from, authUsername, "FOLLOW")

	if approved.Approved && !exists {
		dberr = m.repo.FollowPerson(from, authUsername, "FOLLOW")
	}

	if dberr != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (m *SocialGraphHandler) GetFollowRequests(rw http.ResponseWriter, h *http.Request) {
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

	users, err := m.repo.GetFollowRequests(username)

	if err != nil {
		http.Error(rw, "Cannot find requests", 500)
		return
	}
	jsonResponse(users, rw)
}

func (m *SocialGraphHandler) GetFollowersForUser(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	username := vars["username"]

	users, err := m.repo.GetFollowersForUser(username)

	if err != nil {
		http.Error(rw, "Cannot find requests", 500)
		return
	}
	jsonResponse(users, rw)
}

func (m *SocialGraphHandler) GetFollowingUsers(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	username := vars["username"]

	users, err := m.repo.GetFollowingUsers(username)

	if err != nil {
		http.Error(rw, "Cannot find requests", 500)
		return
	}
	jsonResponse(users, rw)
}

func (m *SocialGraphHandler) ChangePrivacy(rw http.ResponseWriter, h *http.Request) {
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

	var updatePrivacy data.UpdatePrivacy
	err = json.NewDecoder(h.Body).Decode(&updatePrivacy)

	err = m.repo.ChangePrivacy(username, updatePrivacy.Privacy)

	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (m *SocialGraphHandler) GetSuggestedUsers(rw http.ResponseWriter, h *http.Request) {
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

	following, err := m.repo.GetFollowingUsers(username)
	if len(following) == 0 {
		users, err := m.repo.GetUsersFromSameTown(username)

		if len(users) == 0 {
			users, err = m.repo.GetRandomUsers(username)
		}

		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		if len(users) > 0 {
			jsonResponse(users, rw)
			return
		}
	}

	suggestedUsers, err := m.repo.GetSuggestionsForUser(username)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	if len(suggestedUsers) == 0 {
		users, err := m.repo.GetUsersFromSameTown(username)

		if len(users) == 0 {
			users, err = m.repo.GetRandomUsers(username)
		}

		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		if len(users) > 0 {
			jsonResponse(users, rw)
			return
		}
	}
	jsonResponse(suggestedUsers, rw)

}

func GetMapClaims(tokenBytes []byte) map[string]string {
	var claims map[string]string

	err := jwt.ParseClaims(tokenBytes, verifier, &claims)
	if err != nil {
		log.Println(err)
	}

	return claims
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
