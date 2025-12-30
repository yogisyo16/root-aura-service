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
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	TodoID      string    `json:"todo_id,omitempty" bson:"_todo_id,omitempty"`
	TaskDetails string    `json:"task_details" bson:"_task_details"`
	Notes       string    `json:"notes" bson:"_notes"`
	Status      string    `json:"status" bson:"_status"`
	Priority    string    `json:"priority" bson:"_priority"`
	CreatedAt   time.Time `json:"created_at,omitempty" bson:"_created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"_updated_at,omitempty"`
}

func NewTodoDetailsService(mongo *mongo.Client) TodoDetails {
	client = mongo
	return TodoDetails{}
}

func returnTodoDetailsCollection(collection string) *mongo.Collection {
	return client.Database("todos_db").Collection(collection)
}

// GetAllTodosDetails - now exported
func (t *TodoDetails) GetAllTodosDetails() ([]TodoDetails, error) {
	collection := returnTodoDetailsCollection("todo_details")
	var todoDetails []TodoDetails
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todoDetail TodoDetails
		if err := cursor.Decode(&todoDetail); err != nil {
			log.Println(err)
			continue
		}
		todoDetails = append(todoDetails, todoDetail)
	}

	return todoDetails, nil
}

// GetTodoDetailsById - now exported
func (t *TodoDetails) GetTodoDetailsById(id string) (TodoDetails, error) {
	collection := returnTodoDetailsCollection("todo_details")
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

// InsertTodoDetails - create todo details
func (t *TodoDetails) InsertTodoDetails(entry TodoDetails) error {
	collection := returnTodoDetailsCollection("todo_details")
	_, err := collection.InsertOne(context.TODO(), TodoDetails{
		TodoID:      entry.TodoID,
		TaskDetails: entry.TaskDetails,
		Notes:       entry.Notes,
		Status:      entry.Status,
		Priority:    entry.Priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}

// UpdateTodoDetails - update todo details
func (t *TodoDetails) UpdateTodoDetails(id string, entry TodoDetails) (*mongo.UpdateResult, error) {
	collection := returnTodoDetailsCollection("todo_details")
	mongoID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "_task_details", Value: entry.TaskDetails},
			{Key: "_notes", Value: entry.Notes},
			{Key: "_status", Value: entry.Status},
			{Key: "_priority", Value: entry.Priority},
			{Key: "_updated_at", Value: time.Now()},
		}},
	}

	res, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": mongoID},
		update,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return res, nil
}

// DeleteTodoDetails
func (t *TodoDetails) DeleteTodoDetails(id string) error {
	collection := returnTodoDetailsCollection("todo_details")
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = collection.DeleteOne(
		context.Background(),
		bson.M{"_id": mongoID},
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
