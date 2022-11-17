package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

//func New(clientmongo.Client) data.AuthRepo {
//	credentials := client.Database(DATABASE).Collection(CREDENTIALS_COLLECTION)
//	credentials.Drop(context.TODO())
//	return &AuthRepoMongoDb{
//		credentials: credentials,
//	}
//}

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

func (pr *UserRepo) Post(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usersCollection := pr.getCollection()

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
