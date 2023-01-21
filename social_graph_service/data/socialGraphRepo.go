package data

import (
	"context"
	"fmt"
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
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// ExecuteWrite for write transactions (Create/Update/Delete)
	savedUser, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (u:User) SET u.username = $username, u.sex = $sex, u.age = $age, u.town = $town, u.privacy = $privacy RETURN u.username + ', from node ' + id(u)",
				map[string]any{"username": user.Username, "sex": user.Sex, "age": user.Age, "town": user.Town, "privacy": user.Privacy})
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

func (mr *SocialGraphRepo) FollowPerson(from string, to string, relationship string) error {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := fmt.Sprintf("MATCH (a:User), (b:User) WHERE a.username = $from AND b.username = $to CREATE (a)-[r:%s]->(b) RETURN type(r)", relationship)

	session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx, query, map[string]interface{}{"from": from, "to": to})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})

	return nil

}

func (mr *SocialGraphRepo) RemoveFollow(from string, to string, relationship string) error {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := fmt.Sprintf("MATCH (f {username: $from})-[r:%s]->(t {username: $to})DELETE r", relationship)

	session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx, query, map[string]interface{}{"from": from, "to": to})
			if err != nil {
				return nil, err
			}

			return nil, nil
		})

	return nil

}

func (mr *SocialGraphRepo) CheckIfRelationshipExists(usernameFrom, usernameTo, relationship string) (bool, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := fmt.Sprintf("MATCH (f:User {username: $from }), (t:User {username: $to}) RETURN EXISTS( (f)-[:%s]->(t)) as exists", relationship)

	res, _ := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx, query, map[string]interface{}{"from": usernameFrom, "to": usernameTo})
			if err != nil {
				return false, err
			}
			result.Next(ctx)
			r := result.Record()
			if r == nil {
				return false, nil
			}
			res, _ := r.Get("exists")
			return res, nil
		})
	return res.(bool), nil
}

func (mr *SocialGraphRepo) GetUser(username string) (*User, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := "MATCH (user {username: $username}) RETURN user.username as username, user.privacy as privacy"

	user, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx, query, map[string]interface{}{"username": username})

			if err != nil {
				return nil, err
			}
			result.Next(ctx)
			r := result.Record()

			if r == nil {
				return nil, nil
			}

			privacy, _ := r.Get("privacy")
			u, _ := r.Get("username")

			return &User{
				Privacy:  privacy.(string),
				Username: u.(string),
			}, nil
		})

	if err != nil {
		return nil, err
	}

	return user.(*User), nil

}

func (mr *SocialGraphRepo) GetFollowRequests(username string) ([]User, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := "MATCH (user:User)<-[:REQUEST]-(request) WHERE user.username = $username RETURN request.username as username"

	users, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx, query, map[string]interface{}{"username": username})

			if err != nil {
				return nil, err
			}
			var users []User
			for result.Next(ctx) {
				r := result.Record()
				u, _ := r.Get("username")

				users = append(users, User{Username: u.(string)})
			}

			return users, nil
		})

	if err != nil {
		return nil, err
	}

	return users.([]User), nil

}

func (mr *SocialGraphRepo) GetFollowersForUser(username string) ([]User, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := "MATCH (user:User)<-[:FOLLOW]-(request) WHERE user.username = $username RETURN request.username as username"

	users, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx, query, map[string]interface{}{"username": username})

			if err != nil {
				return nil, err
			}
			var users []User
			for result.Next(ctx) {
				r := result.Record()
				u, _ := r.Get("username")

				users = append(users, User{Username: u.(string)})
			}

			return users, nil
		})

	if err != nil {
		return nil, err
	}

	return users.([]User), nil

}

// korisnici koje prati korisnik
func (mr *SocialGraphRepo) GetFollowingUsers(username string) ([]User, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := "MATCH (user:User)-[:FOLLOW]->(request) WHERE user.username = $username RETURN request.username as username"

	users, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx, query, map[string]interface{}{"username": username})

			if err != nil {
				return nil, err
			}
			var users []User
			for result.Next(ctx) {
				r := result.Record()
				u, _ := r.Get("username")

				users = append(users, User{Username: u.(string)})
			}

			return users, nil
		})

	if err != nil {
		return nil, err
	}

	return users.([]User), nil

}

func (mr *SocialGraphRepo) ChangePrivacy(username, isPrivate string) error {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := fmt.Sprintf("MATCH (u:User {username:$username}) set u.privacy= $isPrivate")

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx, query, map[string]interface{}{"username": username, "isPrivate": isPrivate})
			if err != nil {
				return nil, err
			}

			return nil, nil
		})

	if err != nil {
		return err
	}

	return nil

}
