package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	social_graph "tweet_service/client/social-graph"
	"tweet_service/data"
)

type KeyTweet struct{}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

type TweetsHandler struct {
	logger *log.Logger
	// NoSQL: injecting product repository
	repo        *data.TweetRepo
	socialGraph social_graph.Client
}

func NewTweetsHandler(l *log.Logger, r *data.TweetRepo, socialGraph social_graph.Client) *TweetsHandler {
	return &TweetsHandler{l, r, socialGraph}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (p *TweetsHandler) GetAllTweets(rw http.ResponseWriter, h *http.Request) {
	allTweets, err := p.repo.GetAll()
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		p.logger.Fatal("Database exception: ", err)
	}

	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
	renderJSON(rw, allTweets)
}

func (p *TweetsHandler) GetAllUserTweets(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	username := vars["username"]

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

	access, err := p.socialGraph.CanAccessTweet(username, tokenString)

	if err != nil {
		http.Error(rw, "Error 123", http.StatusInternalServerError)
		return
	}

	if username != authUsername && !access {
		http.Error(rw, "cannot access profile tweets", http.StatusForbidden)
		return
	}

	allTweets, err := p.repo.GetTweetListByUsername(username)

	for i, tweet := range allTweets {
		if tweet.Retweet {
			canSee, sgerr := p.socialGraph.CanAccessTweet(tweet.OriginalPostedBy, tokenString)
			if sgerr != nil {
				allTweets[i].Text = ""
			}

			if username != authUsername && !canSee {
				allTweets[i].Text = ""
			}

		}
	}

	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		p.logger.Fatal("Database exception: ", err)
	}

	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
	renderJSON(rw, allTweets)
}

func (p *TweetsHandler) GetUsersWhoLikedTweet(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	id := vars["id"]
	id2, _ := gocql.ParseUUID(id)
	allTweets, err := p.repo.GetUsersWhoLikedTweet(id2)

	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		p.logger.Fatal("Database exception: ", err)
	}

	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
	renderJSON(rw, allTweets)

}

func (p *TweetsHandler) HomeFeed(rw http.ResponseWriter, h *http.Request) {

	bearer := h.Header.Get("Authorization")
	bearerToken := strings.Split(bearer, "Bearer ")
	tokenString := bearerToken[1]
	fmt.Println(tokenString)

	token, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Cannot parse token", 403)
		return
	}

	claims := GetMapClaims(token.Bytes())
	username := claims["username"]

	feedTweets, err := p.repo.GetHomeFeed(username)

	for i, tweet := range feedTweets {
		if tweet.Retweet {
			access, sgerr := p.socialGraph.CanAccessTweet(tweet.OriginalPostedBy, tokenString)
			if sgerr != nil {
				feedTweets[i].Text = ""
			}

			if tweet.OriginalPostedBy != username && !access {
				feedTweets[i].Text = ""
			}

		}
	}

	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		p.logger.Fatal("Database exception: ", err)
	}

	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
	renderJSON(rw, feedTweets)

}

func (p *TweetsHandler) PostTweet(rw http.ResponseWriter, h *http.Request) {

	var request data.Tweet
	err := json.NewDecoder(h.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	bearer := h.Header.Get("Authorization")
	bearerToken := strings.Split(bearer, "Bearer ")
	tokenString := bearerToken[1]
	fmt.Println(tokenString)

	token, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Cannot parse token", 403)
		return
	}

	claims := GetMapClaims(token.Bytes())

	authUsername := claims["username"]

	request.ID = gocql.TimeUUID()
	request.PostedBy = authUsername
	request.Retweet = false
	request.OriginalPostedBy = ""

	usernames, err := p.socialGraph.GetFollowers(authUsername)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	usernames = append(usernames, authUsername)

	tweet, err := p.repo.SaveTweet(&request, usernames)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
	jsonResponse(tweet, rw)
}

func (p *TweetsHandler) Retweet(rw http.ResponseWriter, h *http.Request) {
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

	vars := mux.Vars(h)
	tweetId := vars["id"]

	tweet, err := p.repo.GetTweetById(tweetId)

	access, err := p.socialGraph.CanAccessTweet(tweet.PostedBy, tokenString)

	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	if !access {
		http.Error(rw, "cannot access this tweet", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	retweet := data.Tweet{
		ID:               gocql.TimeUUID(),
		Text:             tweet.Text,
		PostedBy:         authUsername,
		Retweet:          true,
		OriginalPostedBy: tweet.PostedBy,
	}

	usernames, err := p.socialGraph.GetFollowers(authUsername)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	usernames = append(usernames, authUsername)

	responseTweet, err := p.repo.SaveTweet(&retweet, usernames)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	jsonResponse(responseTweet, rw)
}

func (p *TweetsHandler) LikeTweet(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	id := vars["id"]

	var request data.Likes

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

	request.Username = username
	request.ID, _ = gocql.ParseUUID(id)

	tweet, err := p.repo.LikeTweett(&request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	jsonResponse(tweet, rw)
}

func (p *TweetsHandler) UnlikeTweet(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	id := vars["id"]

	var request data.Likes

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

	request.Username = username
	request.ID, _ = gocql.ParseUUID(id)

	tweet, err := p.repo.UnlikeTweet(&request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	jsonResponse(tweet, rw)
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

func GetMapClaims(tokenBytes []byte) map[string]string {
	var claims map[string]string

	err := jwt.ParseClaims(tokenBytes, verifier, &claims)
	if err != nil {
		log.Println(err)
	}

	return claims
}

func (p *TweetsHandler) MiddlewareTweetValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		tweet := &data.Tweet{}
		err := tweet.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		if err != nil {
			p.logger.Println("Error validating product", err)
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(h.Context(), KeyTweet{}, tweet)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}
