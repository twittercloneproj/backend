package data

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	tweets = "tweets/%s"
	all    = "tweets"
)

func generateKey() (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(tweets, id), id
}

func constructKey(id string) string {
	return fmt.Sprintf(tweets, id)
}
