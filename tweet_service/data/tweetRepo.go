package data

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
)

type TweetRepo struct {
	logger *log.Logger
	db     *gocql.Session
}

func New(logger *log.Logger) (*TweetRepo, error) {
	dbport := os.Getenv("DBPORT")
	db := os.Getenv("DB")
	host := fmt.Sprintf("%s:%s", db, dbport)
	cluster := gocql.NewCluster(host)
	cluster.ProtoVersion = 4
	cluster.Keyspace = "tweet_db"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &TweetRepo{db: session}, nil
}

func (s *TweetRepo) GetAll() ([]Tweet, error) {
	var tweet Tweet
	var tweets []Tweet
	iter := s.db.Query(`SELECT id, text, created_on FROM tweets`).Iter()
	for iter.Scan(&tweet.ID, &tweet.Text, &tweet.CreatedOn, &tweet.User) {
		tweets = append(tweets, tweet)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return tweets, nil
}

func (s *TweetRepo) SaveTweet(tweet *Tweet) error {
	err := s.db.Query("INSERT INTO tweets(id, text, created_on, user) VALUES(?, ?, ?, ?)").
		Bind(tweet.ID, tweet.Text, tweet.CreatedOn, tweet.User). // u bazi ako cuvam ID uuid, onda ide sa gocql.UUID, umesto stringa.
		Exec()

	return err
}

//func (pr *TweetRepo) Post(tweet *Tweet) (*Tweet, error) {
//	kv := pr.cli.KV()
//
//	tweet.CreatedOn = time.Now().UTC().String()
//
//	dbId, id := generateKey()
//	tweet.ID = id
//
//	data, err := json.Marshal(tweet)
//	if err != nil {
//		return nil, err
//	}
