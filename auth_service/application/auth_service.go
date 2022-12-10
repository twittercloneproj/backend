package application

import (
	"auth_service/domain"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type AuthService struct {
	store domain.AuthStore
}

func NewAuthService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (service *AuthService) GetAll() ([]*domain.User, error) {
	return service.store.GetAll()
}

func (service *AuthService) LoginHelp(credential *domain.User) (string, error) {

	user, err := service.store.GetOneUser(credential.Username)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password))
	if passError != nil {
		fmt.Println(passError)
		return "", err
	}

	tokenString, err := GenerateJWT(user)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
func GenerateJWT(user *domain.User) (string, error) {

	key := []byte(os.Getenv("SECRET_KEY"))
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		log.Println(err)
	}

	builder := jwt.NewBuilder(signer)

	claims := &domain.Claims{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(time.Minute * 60),
	}

	token, err := builder.Build(claims)
	if err != nil {
		log.Println(err)
	}

	return token.String(), nil
}

func (service *AuthService) Post(user *domain.User) (*domain.User, error) {
	user.ID = primitive.NewObjectID()
	pass := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)

	if user.Email == "" && user.Firm == "" && user.Website == "" {
		user.Role = "Regular"
	} else if user.Email != "" && user.Firm != "" && user.Website != "" {
		user.Role = "Business"
	} else {
		user.Role = "Regular"
	}

	newUser := domain.User{
		Username: user.Username,
		Password: user.Password,
	}
	err = service.store.Post(&newUser)

	if err != nil {
		return nil, err
	}

	return user, nil

}
