package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name,omitempty" json:"name" `
	Surname  string             `bson:"surname,omitempty" json:"surname"`
	Username string             `bson:"username,omitempty" json:"username"`
	Password string             `bson:"password,omitempty" json:"password"`
	Sex      string             `bson:"sex,omitempty" json:"sex"`
	Age      float32            `bson:"age,omitempty" json:"age"`
	Town     string             `bson:"town,omitempty" json:"town"`
	Email    string             `bson:"email,omitempty" json:"email"`
	Firm     string             `bson:"firm,omitempty" json:"firm"`
	Website  string             `bson:"website,omitempty" json:"website"`
	Role     string             `bson:"role,omitempty" json:"role"`
}

//type BusinessUser struct {
//	ID       primitive.ObjectID `bson:"_id" json:"id"`
//	Username string             `bson:"username" json:"username"`
//	Password string             `bson:"password" json:"password"`
//	Email    string             `bson:"email" json:"email"`
//	Firm     string             `bson:"firm" json:"firm"`
//	Website  string             `bson:"website" json:"website"`
//}

//type BusinessUsers []*BusinessUser

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

//func (p *BusinessUsers) ToJSON(w io.Writer) error {
//	e := json.NewEncoder(w)
//	return e.Encode(p)
//}
//
//func (p *BusinessUser) ToJSON(w io.Writer) error {
//	e := json.NewEncoder(w)
//	return e.Encode(p)
//}
//
//func (p *BusinessUser) FromJSON(r io.Reader) error {
//	d := json.NewDecoder(r)
//	return d.Decode(p)
//}

//func DecodeBody(r io.Reader) (*BusinessUser, error) {
//	dec := json.NewDecoder(r)
//	dec.DisallowUnknownFields()
//
//	var businessUser *BusinessUser
//	if err := dec.Decode(&businessUser); err != nil {
//		return nil, err
//	}
//	return businessUser, nil
//}

func RenderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
