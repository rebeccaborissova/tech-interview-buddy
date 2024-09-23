package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://sarah:pf%5FdnZAB2FA8SKy@cen3031-testing.wq4wc.mongodb.net/?retryWrites=true&w=majority&appName=CEN3031-Testing").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("atlasAdmin@admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}

	// Tutorials from MongoDB use ctx, use context.TODO() instead
	// Gives access to the specific collection and sub-collection
	wholeCollection := client.Database("sample_mflix") // Declares a new operator. Single is just reassignment
	usersCollection := wholeCollection.Collection("users")

	// CRUD operation: READ (begin)
	// Getting all the data from user sub-collection of sample_mflix collection
	// Video on how to do this: https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-read-documents
	cursor, err := usersCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}

	// Reading each entry one-by-one of the whole sub-collection. Don't worry about overflowing memory
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var user bson.M
		if err = cursor.Decode(&user); err != nil {
			panic(err)
		}
		fmt.Println(user) // Just user prints out the entire entry. Do user["name"] to specifically print out the name of each entry
	}
	fmt.Println("")

	// Filtering a collection
	// If no one fits the criteria, you just get an empty array
	fmt.Println("Now printing one user, with filter.")
	var filterUsers []bson.M

	// Use bson.M for the filter
	// We can also use userCollection.FindOne. Only use this though if you anticipated only 1 entry back
	filterCursor, err := usersCollection.Find(context.TODO(), bson.M{"name": "Patrick Knight"})
	if err != nil {
		panic(err)
	}
	if err = filterCursor.All(context.TODO(), &filterUsers); err != nil {
		panic(err)
	}
	fmt.Println(filterUsers)
	fmt.Println("")

	// Attemping to make a function to find a user.
	fmt.Println("Making a function to find a user")
	var findUserResult bson.M = findUser("meow meow", usersCollection, err)
	if len(findUserResult) == 0 {
		fmt.Println("USER IS NOT FOUND")
	} else {
		fmt.Println(findUserResult)
	}
	// CRUD operation: READ (end)

	// CRUD operation: UPDATE (begin)

	fmt.Println("We survived, yay!")
}

// Input: Email address
// Output: Ideally, one entry returned with the account associated with that email addreess.
func findUser(email string, database *mongo.Collection, err error) bson.M {
	var toReturn bson.M
	if err = database.FindOne(context.TODO(), bson.M{"email": email}).Decode(&toReturn); err != nil {
		return bson.M{}
	}
	return toReturn
}
