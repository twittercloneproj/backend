package domain

type AuthStore interface {
	GetAll() ([]*User, error)
	Post(user *User) error
	GetOneUser(username string) (*User, error)
}
