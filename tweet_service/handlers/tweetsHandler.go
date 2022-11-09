package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tweet_service/data"
)

type KeyTweet struct{}

type TweetsHandler struct {
	logger *log.Logger
	// NoSQL: injecting product repository
	repo *data.TweetRepo
}

func NewTweetsHandler(l *log.Logger, r *data.TweetRepo) *TweetsHandler {
	return &TweetsHandler{l, r}
}

func (p *TweetsHandler) GetAllTweets(rw http.ResponseWriter, h *http.Request) {
	allTweets, err := p.repo.GetAll()
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		p.logger.Fatal("Database exception: ", err)
	}

	err = allTweets.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *TweetsHandler) GetOneTweet(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	tweet, err := p.repo.Get(id)
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		p.logger.Fatal("Database exception: ", err)
	}

	if tweet == nil {
		http.Error(rw, "Tweet with given id not found", http.StatusNotFound)
		p.logger.Printf("Tweet with id: '%s' not found", id)
		return
	}

	err = tweet.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *TweetsHandler) PostTweet(rw http.ResponseWriter, h *http.Request) {
	tweet := h.Context().Value(KeyTweet{}).(*data.Tweet)
	p.repo.Post(tweet)
	// NoSQL: Bonus - create inside nested key structure
	rw.WriteHeader(http.StatusCreated)
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

func (p *TweetsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		p.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
