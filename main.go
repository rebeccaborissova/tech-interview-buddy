package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"unicode"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	store, err := NewPostgresStore() // Store is a *PostgresStore that is *Mongo.Database
	if err != nil {
		log.Fatal(err)
	}

	usersCollection := store.db.Collection("users")

	// Test attributes
	email := "sarah.tran1@ufl.edu"
	password := "meowmeowmeoW@13"
	firstName := "Sarah"
	lastName := "Tran"
	tookDSA := true
	knowPython := true
	knowCPP := true
	classYear := 2

	// CREATE OPERATION
	// TODO: Should I try to make it return an error and like how should I let front-end know what information given was invalid?
	// TEMPORARILY: Use a string
	insertAccount(email, password, firstName, lastName, tookDSA, knowPython, knowCPP, classYear, usersCollection)
	// CREATE OPERATION

	fmt.Println()

	// READ OPERATION (https://www.youtube.com/watch?v=Ap8elI7ePt4&t=7s)
	email = "meowmeow@ufl.edu"
	account := emailInDatabase(email, usersCollection)
	if account == nil {
		fmt.Println("Account is NOT in database.")
	} else {
		fmt.Println("Account IS in databse.")
	}
	// READ OPERATION

	fmt.Println("")

	// DELETE OPERATION
	// Account will be deleted by email since that will be our primary key.
	email = "blehbleh@ufl.edu"
	password = "oiewjfaw)((HF1))"
	firstName = "Meow"
	lastName = "Meow"
	tookDSA = false
	knowPython = false
	knowCPP = false
	classYear = 1

	insertAccount(email, password, firstName, lastName, tookDSA, knowPython, knowCPP, classYear, usersCollection)
	deleteErr := deleteAccount(email, usersCollection)
	if deleteErr != nil {
		log.Fatal(deleteErr)
	}
	// DELETE OPERATION

	// UPDATE OPERATION
	email = "mother3@ufl.edu"
	password = "QWERTYy19$@"
	firstName = "Mother"
	lastName = "Three"
	tookDSA = false
	knowPython = true
	knowCPP = true
	classYear = 3

	insertAccount(email, password, firstName, lastName, tookDSA, knowPython, knowCPP, classYear, usersCollection)

	// UPDATE OPERATION

	fmt.Println()
	fmt.Println("I am not dead. Yay!")
}

// TODO: CHANGE PASSWORD
/*func updatePassword(email, password string, uses *mongo.Collection) (passwordChanged bool) {
	filter := bson.D{{"email", email}}
}*/

// Assumption that the user will be deleting their account while logged in so there is always be an account.
// Just in case, an err will return.
func deleteAccount(email string, user *mongo.Collection) (err error) {
	filter := bson.D{{Key: "email", Value: email}}

	result, err := user.DeleteOne(context.TODO(), filter)
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)
	return err
}

// Takes the user's password and makes sure that it'll be a little strong
// A user's password is at least 10 characters, contains a special character, a digit, and at least one capital and lowercase letter.
func passwordValidation(password string) (validPassword bool) {
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
func emailInDatabase(email string, user *mongo.Collection) (account *Account) {
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
func insertAccount(email, password, first, last string, dsa, python, cpp bool, year int, users *mongo.Collection) (err error) {
	validation := validateAccount(email, password, first, last, users)

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
func validateAccount(email, password, first, last string, users *mongo.Collection) (validationString [4]bool) {
	valid := [4]bool{true, true, true, true}
	pattern := `@ufl\.edu$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(email) || emailInDatabase(email, users) != nil {
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
	if !passwordValidation(password) {
		valid[3] = false
	}
	return valid
}
