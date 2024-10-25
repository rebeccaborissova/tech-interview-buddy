package main

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostgresStore struct {
	db *mongo.Database
}

// This won't be here for the actual project.
// $ docker run --name some-postgres -e POSTGRES_PASSWORD=gobank -p 5432:5432 -d postgres
// Silly PostgreSGL password: DQ8nZxjXf
func NewPostgresStore() (*PostgresStore, error) {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://sarah:pf%5FdnZAB2FA8SKy@cen3031-testing.wq4wc.mongodb.net/?retryWrites=true&w=majority&appName=CEN3031-Testing").SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("atlasAdmin@admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database!")
	return &PostgresStore{
		db: client.Database("tester"), // RETURNS A DATABASE
	}, nil
}
