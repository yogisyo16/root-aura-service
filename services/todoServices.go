package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database Todo : Struct modal
//
// TODO(auth): UserID exists on the model but nothing currently sets it
// (createTodo never populates it) or filters by it (GetAllTodos returns
// every todo regardless of owner). Once JWT auth lands, todo reads/writes
// should be scoped to the authenticated user's ID. See docs/BACKLOG.md.
type Todo struct {
	ID        string     `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string     `json:"user_id" bson:"_user_id"`
	Task      string     `json:"task" bson:"_task"`
	DateStart *time.Time `json:"date_start,omitempty" bson:"_date_start,omitempty"`
	DateDue   *time.Time `json:"date_due,omitempty" bson:"_date_due,omitempty"`
	Completed bool       `json:"completed" bson:"_completed"`
	CreatedAt time.Time  `json:"created_at,omitempty" bson:"_created_at,omitempty"`
	UpdatedAt time.Time  `json:"update_at,omitempty" bson:"_update_at,omitempty"`
}

// client is the shared Mongo client for the whole services package, set once
// here and reused by userServices.go and todoDetailsServices.go too — there's
// only ever one client for the process's lifetime.
var client *mongo.Client

// New sets the package-level Mongo client and returns a zero-value Todo
// service. Called once at startup from api/main.go; the User and
// TodoDetails services rely on the same client being set here first.
func New(mongo *mongo.Client) Todo {
	client = mongo
	return Todo{}
}

func returnTodosCollection(collection string) *mongo.Collection {
	return client.Database("todos_db").Collection(collection)
}

// GetAllTodos
func (t *Todo) GetAllTodos() ([]Todo, error) {
	collection := returnTodosCollection("todos")
	var todos []Todo
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	return todos, nil
}

// GetTodoById
func (t *Todo) GetTodoById(id string) (Todo, error) {
	collection := returnTodosCollection("todos")
	var todo Todo

	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Todo{}, err
	}

	err = collection.FindOne(context.TODO(), bson.M{"_id": mongoID}).Decode(&todo)
	if err != nil {
		log.Println(err)
		return Todo{}, err
	}

	return todo, nil
}

// InsertTodo
func (t *Todo) InsertTodo(entry Todo) (*Todo, error) {
	collection := returnTodosCollection("todos")
	now := time.Now()
	docToInsert := Todo{
		Task:      entry.Task,
		DateStart: entry.DateStart,
		DateDue:   entry.DateDue,
		Completed: entry.Completed,
		CreatedAt: now,
		UpdatedAt: now,
	}

	res, err := collection.InsertOne(context.TODO(), docToInsert)
	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		docToInsert.ID = oid.Hex()
	}

	return &docToInsert, nil
}

// UpdatedTodo

func (t *Todo) UpdatedTodo(id string, entry Todo) (*mongo.UpdateResult, error) {
	collection := returnTodosCollection("todos")
	mongoID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "_task", Value: entry.Task},
			{Key: "_date_start", Value: entry.DateStart},
			{Key: "_date_due", Value: entry.DateDue},
			{Key: "_completed", Value: entry.Completed},
			{Key: "_update_at", Value: time.Now()},
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

// DeleteTodo

func (t *Todo) DeleteTodo(id string) error {
	collection := returnTodosCollection("todos")
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
