package tools

import (
	"context"
	"fmt"
	"regexp"
	"unicode"

	"github.com/alexedwards/argon2id"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostgresStore struct {
	DB *mongo.Database
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
		DB: client.Database("tester"), // RETURNS A DATABASE
	}, nil
}

func UpdatePassword(email, password string, users *mongo.Collection) (passErr error) {
	filter := bson.D{{Key: "email", Value: email}}

	password, _ = argon2id.CreateHash(password, argon2id.DefaultParams)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: password}}}}

	_, err := users.UpdateOne(context.TODO(), filter, update)
	return err
}

func IsCorrectPassword(email, password string, users *mongo.Collection) (isValid bool, err error) {
	filter := bson.D{{Key: "email", Value: email}}
	account := users.FindOne(context.TODO(), filter)

	attemptHash, _ := argon2id.CreateHash(password, argon2id.DefaultParams)
	var decodedAccount bson.M
	err = account.Decode(&decodedAccount)
	actualHash := decodedAccount["password"]
	fmt.Printf("actualHash: %v\n", actualHash)
	fmt.Printf("attemptHash: %v\n", attemptHash)

	var matches bool = (attemptHash == actualHash)
	return matches, err
}

// Assumption that the user will be deleting their account while logged in so there is always be an account.
// Just in case, an err will return.
func DeleteAccount(email string, user *mongo.Collection) (err error) {
	filter := bson.D{{Key: "email", Value: email}}

	result, err := user.DeleteOne(context.TODO(), filter)
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)
	return err
}

// Takes the user's password and makes sure that it'll be a little strong
// A user's password is at least 10 characters, contains a special character, a digit, and at least one capital and lowercase letter.
func PasswordValidation(password string) (validPassword bool) {
	re := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	containsSpecialCharas := re.MatchString(password)

	re = regexp.MustCompile(`[0-9]`)
	containsDigits := re.MatchString(password)

	re = regexp.MustCompile(`[A-Z]`)
	containsCapital := re.MatchString(password)

	re = regexp.MustCompile(`[a-z]`)
	containsLowercase := re.MatchString(password)

	return containsSpecialCharas && containsDigits && containsCapital && containsLowercase && (len(password) >= 10)
}

// Takes an email and see if it is in the database.
// Exists - Return an Account object
// Does not Exist - Return NIL
// DOES NOT check if the email is a UFL email.
func EmailInDatabase(email string, user *mongo.Collection) (account *Account) {
	var acc Account
	filter := bson.M{"email": email}

	err := user.FindOne(context.TODO(), filter).Decode(&acc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return nil
	}

	return &acc
}

// Takes a new account created and inserts it into the collection of users.
// NOTE/TODO: Double check that the email is not in the database.
func InsertAccount(email, password, first, last string, dsa, python, cpp bool, year int, users *mongo.Collection) (err error) {
	validation := ValidateAccount(email, password, first, last, users)

	isValid := true
	for i := 0; i < len(validation); i++ {
		if !validation[i] {
			isValid = false
			if i == 0 {
				fmt.Println("Invalid Email. Must be an UFL email or this email is already being used.")
			} else if i == 1 {
				fmt.Println("First name must only contain letters.")
			} else if i == 2 {
				fmt.Println("Last name must only contain letters.")
			} else if i == 3 {
				fmt.Println("Password is invalid. Must be at least 16 characters and contains at least one special character, digit, capital letter, and lowercase letter.")
			}
		}
	}

	if isValid {
		// Make password encryption here.
		password, _ = argon2id.CreateHash(password, argon2id.DefaultParams)
		user := NewAccount(email, password, first, last, dsa, python, cpp, year)
		_, err := users.InsertOne(context.TODO(), user)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ensures that attributes are valid. The following are the requirements:
// Email: Must be a UFL email. Hence the email must contain @ufl.edu and must not already be in the database
// First Name: Must only contain letters.
// Last Name: Must only contain letters.
// NOTE/TODO: I might change this function so that it auto capitalizes the first letter of the name.
func ValidateAccount(email, password, first, last string, users *mongo.Collection) (validationString [4]bool) {
	valid := [4]bool{true, true, true, true}
	pattern := `@ufl\.edu$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(email) || EmailInDatabase(email, users) != nil {
		valid[0] = false
	}
	for _, r := range first {
		if !unicode.IsLetter(r) {
			valid[1] = false
		}
	}
	for _, r := range last {
		if !unicode.IsLetter(r) {
			valid[2] = false
		}
	}
	if !PasswordValidation(password) {
		valid[3] = false
	}
	return valid
}
