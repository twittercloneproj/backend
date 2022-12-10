package data

//
//import (
//	"auth_service/domain"
//	"context"
//	"fmt"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"go.mongodb.org/mongo-driver/mongo/readpref"
//	"golang.org/x/crypto/bcrypt"
//	"log"
//	"os"
//	"strings"
//	"time"
//)
//
//type UserRepo struct {
//	cli         *mongo.Client
//	logger      *log.Logger
//	credentials mongo.Collection
//}
//
//const (
//	DATABASE               = "auth"
//	CREDENTIALS_COLLECTION = "credentials"
//)
//
//type AuthRepoMongoDb struct {
//	credentials mongo.Collection
//}
//
//type Mail struct {
//	senderId string
//	toIds    []string
//	subject  string
//	body     string
//}
//
//type SmtpServer struct {
//	host string
//	port string
//}
//
//func (s *SmtpServer) ServerName() string {
//	return s.host + ":" + s.port
//}
//
//func (mail *Mail) BuildMessage() string {
//	message := ""
//	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
//	if len(mail.toIds) > 0 {
//		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
//	}
//
//	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
//	message += "\r\n" + mail.body
//
//	return message
//}
//
//func New(ctx context.Context, logger *log.Logger) (*UserRepo, error) {
//	db := os.Getenv("AUTH_DB_HOST")
//	dbport := os.Getenv("AUTH_DB_PORT")
//
//	host := fmt.Sprintf("%s:%s", db, dbport)
//	client, err := mongo.NewClient(options.Client().ApplyURI(`mongodb://` + host))
//	if err != nil {
//		return nil, err
//	}
//
//	err = client.Connect(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	return &UserRepo{
//		cli:    client,
//		logger: logger,
//	}, nil
//}
//
//func (store *UserRepo) GetOneUser(username string) (*domain.User, error) {
//	filter := bson.M{"username": username}
//
//	user, err := store.filterOne(filter)
//	if err != nil {
//		return nil, err
//	}
//
//	return user, nil
//}
//
//func (store *UserRepo) filterOne(filter interface{}) (user *domain.User, err error) {
//	result := store.credentials.FindOne(context.TODO(), filter)
//	err = result.Decode(&user)
//	return
//}
//
//func (pr *UserRepo) GetAll() (domain.Users, error) {
//	// Initialise context (after 5 seconds timeout, abort operation)
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	usersCollection := pr.getCollection()
//
//	var users domain.Users
//	usersCursor, err := usersCollection.Find(ctx, bson.M{})
//	if err != nil {
//		pr.logger.Println(err)
//		return nil, err
//	}
//	if err = usersCursor.All(ctx, &users); err != nil {
//		pr.logger.Println(err)
//		return nil, err
//	}
//	return users, nil
//}
//
//func HashPassword(password string) (string, error) {
//	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
//	return string(bytes), err
//}
//
//func CheckPasswordHash(password, hash string) bool {
//	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//	return err == nil
//}
//
//// Disconnect from database
//func (pr *UserRepo) Disconnect(ctx context.Context) error {
//	err := pr.cli.Disconnect(ctx)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// Check database connection
//func (pr *UserRepo) Ping() {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	// Check connection -> if no error, connection is established
//	err := pr.cli.Ping(ctx, readpref.Primary())
//	if err != nil {
//		pr.logger.Println(err)
//	}
//
//	// Print available databases
//	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
//	if err != nil {
//		pr.logger.Println(err)
//	}
//	fmt.Println(databases)
//}
//
//func (pr *UserRepo) getCollection() *mongo.Collection {
//	userDatabase := pr.cli.Database("mongodb")
//	usersCollection := userDatabase.Collection("users")
//	return usersCollection
//}
