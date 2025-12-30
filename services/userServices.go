package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
