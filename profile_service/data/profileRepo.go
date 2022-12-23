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

func (pr *UserRepo) Update(username string, privacy string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	patientsCollection := pr.getCollection()

	//objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{
		"privacy": privacy,
	}}
	result, err := patientsCollection.UpdateOne(ctx, filter, update)
	pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		pr.logger.Println(err)
		return err
	}
	return nil
}

func (pr *UserRepo) getCollection() *mongo.Collection {
	userDatabase := pr.cli.Database("mongodb")
	usersCollection := userDatabase.Collection("users")
	return usersCollection
}

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
