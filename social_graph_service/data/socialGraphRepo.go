package data

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
)

type SocialGraphRepo struct {
	// Thread-safe instance which maintains a database connection pool
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func New(logger *log.Logger) (*SocialGraphRepo, error) {
	// Local instance
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &SocialGraphRepo{
		driver: driver,
		logger: logger,
	}, nil
}

// Check if connection is established
func (mr *SocialGraphRepo) CheckConnection() {
	ctx := context.Background()
	err := mr.driver.VerifyConnectivity(ctx)
	if err != nil {
		mr.logger.Panic(err)
		return
	}
	// Print Neo4J server address
	mr.logger.Printf(`Neo4J server address: %s`, mr.driver.Target().Host)
}

// Disconnect from database
func (mr *SocialGraphRepo) CloseDriverConnection(ctx context.Context) {
	mr.driver.Close(ctx)
}

func (mr *SocialGraphRepo) WritePerson(user *User) error {
	// Neo4J Sessions are lightweight so we create one for each transaction (Cassandra sessions are not lightweight!)
	// Sessions are NOT thread safe
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// ExecuteWrite for write transactions (Create/Update/Delete)
	savedUser, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (u:User) SET u.username = $username RETURN u.username + ', from node ' + id(u)",
				map[string]any{"username": user.Username})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		mr.logger.Println("Error inserting user:", err)
		return err
	}
	mr.logger.Println(savedUser.(string))
	return nil
}
