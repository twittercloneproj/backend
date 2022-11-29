package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type UserRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

const (
	DATABASE               = "auth"
	CREDENTIALS_COLLECTION = "credentials"
)

type AuthRepoMongoDb struct {
	credentials mongo.Collection
}

func New(ctx context.Context, logger *log.Logger) (*UserRepo, error) {
	db := os.Getenv("AUTH_DB_HOST")
	dbport := os.Getenv("AUTH_DB_PORT")

	host := fmt.Sprintf("%s:%s", db, dbport)
	client, err := mongo.NewClient(options.Client().ApplyURI(`mongodb://` + host))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &UserRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// GetOneUser TODO
func (pr *UserRepo) GetOneUser(username string) (*User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	patientsCollection := pr.getCollection()

	var user User
	//usrname, _ := primitive.ObjectIDFromHex(username)
	err := patientsCollection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&user)
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return &user, nil
}

//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//defer cancel()
//
//patientsCollection := pr.getCollection()
//
//var patients Users
//patientsCursor, err := patientsCollection.Find(ctx, bson.M{"username": username})
//if err != nil {
//	pr.logger.Println(err)
//	return nil, err
//}
//if err = patientsCursor.All(ctx, &patients); err != nil {
//	pr.logger.Println(err)
//	return nil, err
//}
//return patients, nil
//}

func (store *AuthRepoMongoDb) filterOne(filter interface{}) (user *User, err error) {
	result := store.credentials.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func (pr *UserRepo) GetAll() (Users, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := pr.getCollection()

	var users Users
	usersCursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	if err = usersCursor.All(ctx, &users); err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return users, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (pr *UserRepo) Post(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usersCollection := pr.getCollection()

	pass := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	if user.Email == "" && user.Firm == "" && user.Website == "" {
		user.Role = Regular
	} else if user.Email != "" && user.Firm != "" && user.Website != "" {
		user.Role = Business
	} else {
		user.Role = Regular
	}

	result, err := usersCollection.InsertOne(ctx, &user)
	if err != nil {
		pr.logger.Println(err)
		return err
	}
	pr.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

// Disconnect from database
func (pr *UserRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *UserRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println(databases)
}

func (pr *UserRepo) getCollection() *mongo.Collection {
	userDatabase := pr.cli.Database("mongodb")
	usersCollection := userDatabase.Collection("users")
	return usersCollection
}
