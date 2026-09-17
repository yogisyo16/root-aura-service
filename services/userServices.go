package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User is both the Mongo document shape (todos_db.users) and the JSON wire
// shape returned by the user endpoints.
//
// TODO(security): Password is tagged `json:"password,omitempty"`, so the
// bcrypt hash currently gets serialized straight back out in GetAllUsers /
// GetUserByID responses. Should be `json:"-"`. See docs/BACKLOG.md.
type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string    `json:"first_name,omitempty" bson:"_first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty" bson:"_last_name,omitempty"`
	Email     string    `json:"email,omitempty" bson:"_email,omitempty"`
	Password  string    `json:"password,omitempty" bson:"_password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"_created_at,omitempty"`
	UpdatedAt time.Time `json:"update_at,omitempty" bson:"_update_at,omitempty"`
}

type UserService interface {
	GetAllUsers() ([]User, error)
	InsertUser(entry User) error
}

// retunrUserCollection (sic) resolves a collection off the shared Mongo
// client. That client is package-level state set once by services.New() /
// services.NewTodoDetailsService() — all three service files
// (todoServices.go, userServices.go, todoDetailsServices.go) reuse it.
func retunrUserCollection(collection string) *mongo.Collection {
	return client.Database("todos_db").Collection(collection)
}

func (u *User) GetAllUsers() ([]User, error) {
	collection := retunrUserCollection("users")
	var users []User
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}

	return users, nil
}

// InsertUser stores a new user document. It expects entry.Password to
// already be a bcrypt hash — hashing happens in the handler
// (handlers/userHandler.go's insertUser), not here.
func (u *User) InsertUser(entry User) error {
	collection := retunrUserCollection("users")
	_, err := collection.InsertOne(context.TODO(), User{
		FirstName: entry.FirstName,
		LastName:  entry.LastName,
		Email:     entry.Email,
		Password:  entry.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error: ", err)
		return err
	}

	return nil
}

// GetUserByID looks a user up by their Mongo ObjectID hex string.
//
// TODO(auth): there's no GetUserByEmail yet, which a login flow will need
// to look up a user by the credential they actually log in with. See
// docs/BACKLOG.md.
func (u *User) GetUserByID(id string) (User, error) {
	collection := retunrUserCollection("users")
	var user User

	// Convert string id to MongoDB ObjectID
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid Hex ID:", err)
		return User{}, err
	}

	// Use mongoID in the query instead of the raw string id
	err = collection.FindOne(context.TODO(), bson.M{"_id": mongoID}).Decode(&user)
	if err != nil {
		log.Println("Error finding user: ", err)
		return User{}, err
	}
	return user, nil
}
