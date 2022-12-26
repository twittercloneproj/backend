package data

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name,omitempty" json:"name" `                            // o
	Surname  string             `bson:"surname,omitempty" json:"surname"`                       // o
	Username string             `bson:"username,omitempty" json:"username" validate:"required"` // o
	Password string             `bson:"password,omitempty" json:"password" validate:"required"` // o
	Sex      string             `bson:"sex,omitempty" json:"sex"`                               // o
	Age      float32            `bson:"age,omitempty" json:"age"`                               // o
	Town     string             `bson:"town,omitempty" json:"town"`                             // o

	Email   string `bson:"email,omitempty" json:"email"`     // b
	Firm    string `bson:"firm,omitempty" json:"firm"`       // b
	Website string `bson:"website,omitempty" json:"website"` // b

	Role    Role   `bson:"role,omitempty" json:"role"`
	Privacy string `bson:"privacy,omitempty" json:"privacy"`
}

const (
	Private = "Private"
	Public  = "Public"
)

type Claims struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Role     Role               `json:"role"`
	jwt.RegisteredClaims
}

type Role string

const (
	Regular  = "Regular"
	Business = "Business"
)

type Users []*User

func (p *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func RenderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
