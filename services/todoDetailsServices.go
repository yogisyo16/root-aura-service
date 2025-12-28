package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoDetails struct {
	todoId        string    `json:"todo_id,omitempty" bson:"_todo_id,omitempty"`
	todo          Todo      `json:"todo,omitempty" bson:"_todo,omitempty"`
	user          User      `json:"user,omitempty" bson:"_user,omitempty"`
	taskDetails   string    `json:"task_details,omitempty" bson:"_task_details,omitempty"`
	collaborators []User    `json:"collaborators,omitempty" bson:"_collaborators,omitempty"`
	createdAt     time.Time `json:"created_at,omitempty" bson:"_created_at,omitempty"`
	updatedAt     time.Time `json:"update_at,omitempty" bson:"_update_at,omitempty"`
}

func NewTodoDetailsService(mongo *mongo.Client) TodoDetails {
	client = mongo
	return TodoDetails{}
}

func (t *TodoDetails) getAllTodosDetails() ([]TodoDetails, error) {
	collection := returnTodosCollection("todos")
	var todoDetails []TodoDetails
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todoDetail TodoDetails
		cursor.Decode(&todoDetail)
		todoDetails = append(todoDetails, todoDetail)
	}

	return todoDetails, nil
}

func (t *TodoDetails) getTodoDetailsById(id string) (TodoDetails, error) {
	collection := returnTodosCollection("todos")
	var todoDetail TodoDetails

	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return TodoDetails{}, err
	}

	err = collection.FindOne(context.TODO(), bson.M{"_id": mongoID}).Decode(&todoDetail)
	if err != nil {
		log.Println(err)
		return TodoDetails{}, err
	}

	return todoDetail, nil
}
