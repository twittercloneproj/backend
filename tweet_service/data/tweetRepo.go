package data

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"io"
	"log"
	"net/http"
	"os"
)

type TweetRepo struct {
	logger *log.Logger
	db     *gocql.Session
}

var Session *gocql.Session

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
	for iter.Scan(&tweet.ID, &tweet.Text, &tweet.CreatedOn) {
		fmt.Println("ID= "+tweet.ID, "Tekst= "+tweet.Text, "Kreiran(datum):"+tweet.CreatedOn)
		tweets = append(tweets, tweet)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return tweets, nil
}

func Post(w http.ResponseWriter, r *http.Request) {
	var Newtweet Tweet
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "wrong data")
	}
	json.Unmarshal(reqBody, &Newtweet)
	if err := Session.Query("INSERT INTO tweets(id, text, created_on) VALUES(?, ?, ?)",
		Newtweet.ID, Newtweet.Text, Newtweet.CreatedOn).Exec(); err != nil {
		fmt.Println("Error while inserting")
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	Conv, _ := json.MarshalIndent(Newtweet, "", " ")
	fmt.Fprintf(w, "%s", string(Conv))

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
