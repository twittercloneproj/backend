package data

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"time"
)

// TweetRepo struct encapsulating Consul api client
type TweetRepo struct {
	cli    *api.Client
	logger *log.Logger
}

// Constructor which reads db configuration from environment
func New(logger *log.Logger) (*TweetRepo, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &TweetRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Returns all tweets
func (pr *TweetRepo) GetAll() (Tweets, error) {
	kv := pr.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	tweets := Tweets{}
	for _, pair := range data {
		tweet := &Tweet{}
		err = json.Unmarshal(pair.Value, tweet)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

// NoSQL: Returns Tweet by id
func (pr *TweetRepo) Get(id string) (*Tweet, error) {
	kv := pr.cli.KV()

	pair, _, err := kv.Get(constructKey(id), nil)
	if err != nil {
		return nil, err
	}
	// If pair is nil -> no object found for given id -> return nil
	if pair == nil {
		return nil, nil
	}

	tweet := &Tweet{}
	err = json.Unmarshal(pair.Value, tweet)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

// Saves Tweet to DB
func (pr *TweetRepo) Post(tweet *Tweet) (*Tweet, error) {
	kv := pr.cli.KV()

	tweet.CreatedOn = time.Now().UTC().String()

	dbId, id := generateKey()
	tweet.ID = id

	data, err := json.Marshal(tweet)
	if err != nil {
		return nil, err
	}

	tweetKeyValue := &api.KVPair{Key: dbId, Value: data}
	_, err = kv.Put(tweetKeyValue, nil)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}
