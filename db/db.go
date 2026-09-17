package db

import (
	"context"
	"log"
	"os"
	"time" // Import time for the context timeout

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// collection is currently unused package state, left over from an earlier
// version of this file — nothing assigns to or reads it today.
var collection *mongo.Collection

// ConnectToMongo opens and verifies (via Ping) a connection to the local
// MongoDB instance, authenticating with the MONGO_DB_USERNAME /
// MONGO_DB_PASSWORD env vars (see .env / docker-compose.yml).
//
// The host is hardcoded to localhost:27017 rather than read from an env var.
// That's intentional for the current setup: only Mongo runs in Docker
// (docker-compose.yml), while this API process runs natively on the host and
// reaches Mongo through its published port — see docs/RUNNING.md. If the API
// is ever containerized too, this URI will need to become configurable
// (e.g. pointing at a compose service name instead of localhost).
func ConnectToMongo() (*mongo.Client, error) {
	// MongoDB connection string
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Getting username and password from .env
	username := os.Getenv("MONGO_DB_USERNAME")
	password := os.Getenv("MONGO_DB_PASSWORD")

	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	// Set up a context with a timeout for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping database: ", err)
		return nil, err
	}

	log.Println("Successfully connected and pinged MongoDB!")

	return client, nil
}
