package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

func (p *TweetsHandler) PostTweet(rw http.ResponseWriter, h *http.Request) {
	//#Todo
	//tweet := h.Context().Value(KeyTweet{}).(*data.Tweet)
	//p.repo.Post(tweet)
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
