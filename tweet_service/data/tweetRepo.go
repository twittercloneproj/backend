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
	iter := s.db.Query(`SELECT id, text, posted_by FROM tweet_by_user`).Iter()
	for iter.Scan(&tweet.ID, &tweet.Text, &tweet.PostedBy) {
		tweets = append(tweets, tweet)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return tweets, nil
}

func (s *TweetRepo) GetTweetListByUsername(username string) ([]Tweet, error) {
	var tweet Tweet
	var tweets []Tweet
	iter := s.db.Query(`SELECT id, text, posted_by, retweet, original_posted_by FROM tweet_by_user WHERE posted_by = ?`).Bind(username).Iter()
	for iter.Scan(&tweet.ID, &tweet.Text, &tweet.PostedBy, &tweet.Retweet, &tweet.OriginalPostedBy) {
		tweets = append(tweets, tweet)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return tweets, nil
}

func (s *TweetRepo) SaveTweet(tweet *Tweet, usernames []string) (*Tweet, error) {
	err := s.db.Query("INSERT INTO tweet_by_user(id, text, posted_by, retweet, original_posted_by) VALUES(?, ?, ?, ?, ?)").
		Bind(tweet.ID, tweet.Text, tweet.PostedBy, tweet.Retweet, tweet.OriginalPostedBy).
		Exec()

	for _, username := range usernames {
		err = s.db.Query("INSERT INTO feed_by_user(id, text, posted_by, username, retweet, original_posted_by) VALUES(?, ?, ?, ?, ?, ?)").
			Bind(tweet.ID, tweet.Text, tweet.PostedBy, username, tweet.Retweet, tweet.OriginalPostedBy).
			Exec()

	}

	if err != nil {
		println(err)
		return nil, err
	}

	return tweet, nil
}

func (s *TweetRepo) LikeTweett(like *Likes) (*Likes, error) {
	err := s.db.Query("INSERT INTO likes(id, username) VALUES(?, ?)").
		Bind(like.ID, like.Username).
		Exec()
	if err != nil {
		println(err)
		return nil, err
	}
	return like, nil
}

func (s *TweetRepo) UnlikeTweet(like *Likes) (*Likes, error) {
	err := s.db.Query("DELETE FROM tweet_db.likes WHERE id=? and username=?").
		Bind(like.ID, like.Username).
		Exec()
	if err != nil {
		println(err)
		return nil, err
	}
	return like, nil
}

func (s *TweetRepo) GetUsersWhoLikedTweet(id gocql.UUID) ([]Likes, error) {
	var like Likes
	var likes []Likes
	iter := s.db.Query(`SELECT * FROM likes WHERE id = ?`).Bind(id).Iter()
	for iter.Scan(&like.ID, &like.Username) {
		likes = append(likes, like)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return likes, nil
}

func (s *TweetRepo) GetHomeFeed(username string) ([]Tweet, error) {
	var tweet Tweet
	var tweets []Tweet
	iter := s.db.Query(`SELECT id, posted_by, text, retweet, original_posted_by FROM feed_by_user WHERE username = ?`).Bind(username).Iter()
	for iter.Scan(&tweet.ID, &tweet.PostedBy, &tweet.Text, &tweet.Retweet, &tweet.OriginalPostedBy) {
		tweets = append(tweets, tweet)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return tweets, nil
}

func (s *TweetRepo) GetTweetById(id string) (Tweet, error) {
	var tweet Tweet
	err := s.db.Query(`SELECT id, posted_by, text, retweet, original_posted_by FROM tweet_by_user WHERE id = ?`).
		Bind(id).
		Scan(&tweet.ID, &tweet.PostedBy, &tweet.Text, &tweet.Retweet, &tweet.OriginalPostedBy)

	return tweet, err
}
