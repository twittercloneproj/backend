package data

import (
	"encoding/json"
	"io"
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

type Tweet struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	CreatedOn string `json:"created_on"`
}
type Tweets struct {
	tweets []Tweet
}

func (p *Tweets) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Tweet) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Tweet) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
