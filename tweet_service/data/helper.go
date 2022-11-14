package data

import (
	"encoding/json"
	"fmt"
	"io"

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

func DecodeTweetBody(r io.Reader) (*Tweet, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var tweet *Tweet
	if err := dec.Decode(&tweet); err != nil {
		return nil, err
	}
	return tweet, nil
}
