package tools

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/alexedwards/argon2id"
	"github.com/gofrs/uuid/v5"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
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

func GetUserCollection(db *mongo.Database) (collection *mongo.Collection) {
	return db.Collection("users")
}

func GetSessionCollection(db *mongo.Database) (collection *mongo.Collection) {
	return db.Collection("sessions")
}

// FUNCTIONS FOR ACCOUNTS ENDS HERE
// Takes a user's email attempted password, and a Mongo collection
// Returns true is the attempted password matches the stored hashed password
// Inspired by https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
func IsCorrectPassword(email string, password string, users *mongo.Collection) (match bool, err error) {
	// var (
	// 	ErrInvalidHash = errors.New("The encoded hash is not in the correct format")
	// 	ErrIncompatibleVersion = errors.New("Incompatible version of argon2")
	// )

	filter := bson.D{{Key: "email", Value: email}}
	account := users.FindOne(context.TODO(), filter)

	var decodedAccount bson.M
	err = account.Decode(&decodedAccount)
	if err != nil {
		log.Error(err)
	}

	actualHash := decodedAccount["password"].(string) // type assertion to string
	match, err = argon2id.ComparePasswordAndHash(password, actualHash)

	fmt.Printf("actualHash: %v\n", actualHash)
	fmt.Printf("Match: %v\n", match)

	return match, err
}

// Assumption that the user will be deleting their account while logged in so there is always be an account.
// Just in case, an err will return.
func DeleteAccount(email string, user *mongo.Collection) (err error) {
	filter := bson.D{{Key: "email", Value: email}}

	result, err := user.DeleteOne(context.TODO(), filter)
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)
	return err
}

// Takes a new account created and inserts it into the collection of users.
func InsertAccount(email, password, first, last string, dsa bool, year int, users *mongo.Collection) (err error) {
	validation := ValidateAccount(email, password, first, last, users)

	// TODO: Return the println statements as error types instead.
	isValid := true
	for i := 0; i < len(validation); i++ {
		if !validation[i] {
			isValid = false
			return errors.New("Does not meet new account requirements")
		}
	}

	if isValid {
		first = strings.ToUpper(string(first[0])) + strings.ToLower(first[1:])
		last = strings.ToUpper(string(last[0])) + strings.ToLower(last[1:])

		// Make password encryption here.
		password, _ = argon2id.CreateHash(password, argon2id.DefaultParams)
		user := NewAccount(email, password, first, last, dsa, year)
		_, err := users.InsertOne(context.TODO(), user)
		if err != nil {
			return err
		}
	}

	return nil
}

// v ALL THE UPDATE FUNCTIONS FOR EACH OF THE FIELDS BESIDES EMAIL v
func UpdatePassword(email, password string, users *mongo.Collection) (passErr error) {
	var InvalidPasswordError = errors.New("password does not meet security requirements")

	filter := bson.D{{Key: "email", Value: email}}

	if PasswordValidation(password) {
		// TODO: argon2id.DefaultParams should be changed once in production, okay for dev
		password, _ = argon2id.CreateHash(password, argon2id.DefaultParams)
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: password}}}}

		_, err := users.UpdateOne(context.TODO(), filter, update)
		return err
	} else {
		return InvalidPasswordError
	}

}

func UpdateDSA(email string, dsa bool, users *mongo.Collection) (err error) {
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "TakenDSA", Value: dsa}}}}

	_, err = users.UpdateOne(context.TODO(), filter, update)
	return err
}

func UpdateYear(email string, year int, users *mongo.Collection) (err error) {
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "Year", Value: year}}}}

	_, err = users.UpdateOne(context.TODO(), filter, update)
	return err
}

func UpdateFirstName(email, firstname string, users *mongo.Collection) (err error) {
	var InvalidFirstError = errors.New("first name can only contain letters")

	filter := bson.D{{Key: "email", Value: email}}
	if ContainsLettersOnly(firstname) {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "FirstName", Value: firstname}}}}

		_, err = users.UpdateOne(context.TODO(), filter, update)
		return err
	} else {
		return InvalidFirstError
	}
}

func UpdateLastName(email, lastname string, users *mongo.Collection) (err error) {
	var InvalidFirstError = errors.New("last name can only contain letters")

	filter := bson.D{{Key: "email", Value: email}}
	if ContainsLettersOnly(lastname) {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "FirstName", Value: lastname}}}}

		_, err = users.UpdateOne(context.TODO(), filter, update)
		return err
	} else {
		return InvalidFirstError
	}
}

// ^ ALL THE UPDATE FUNCTIONS FOR EACH OF THE FIELDS BESIDES EMAIL ^

// Ensures that attributes are valid. The following are the requirements:
// Email: Must be a UFL email. Hence the email must contain @ufl.edu and must not already be in the database
// First Name: Must only contain letters.
// Last Name: Must only contain letters.
func ValidateAccount(email, password, first, last string, users *mongo.Collection) (validationString [4]bool) {
	valid := [4]bool{true, true, true, true}
	pattern := `@ufl\.edu$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(email) || EmailInDatabase(email, users) != nil {
		valid[0] = false
	}
	valid[1] = ContainsLettersOnly(first)
	valid[2] = ContainsLettersOnly(last)
	valid[3] = PasswordValidation(password)
	return valid
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
func EmailInDatabase(email string, user *mongo.Collection) (account *Account) {
	var acc Account
	filter := bson.M{"email": email}

	err := user.FindOne(context.TODO(), filter).Decode(&acc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}

	return &acc
}

func ContainsLettersOnly(str string) (applies bool) {
	for _, r := range str {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func GetOnlineAccounts(user *mongo.Collection) (accounts []Account, err error) {
	filter := bson.D{{Key: "online", Value: true}}

	cursor, err := user.Find(context.TODO(), filter)

	var results []Account
	err = cursor.All(context.TODO(), &results)

	return results, err
}

// FUNCTIONS FOR ACCOUNTS ENDS HERE //

// SESSION HANDLING BEGINS HERE //
func AddSession(sessionToken uuid.UUID, username string, expiresAt time.Time, sessions, users *mongo.Collection) (err error) {
	var UsernameNotFound = errors.New("Email was not found")
	if EmailInDatabase(username, users) == nil {
		return UsernameNotFound
	}

	filter := bson.D{{Key: "email", Value: username}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "online", Value: true}}}}

	_, err = users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}

	session := NewSession(sessionToken, username, expiresAt)
	_, err = sessions.InsertOne(context.TODO(), session)
	if err != nil {
		return err
	}

	return nil
}

func CheckSession(token Session, sessions, users *mongo.Collection) (err error) {
	var TokenNotFound = errors.New("Token not found.")
	var UserNotFound = errors.New("User for token not found.")
	var TokenExpired = errors.New("Session expired.")

	var sesh Session
	filter := bson.M{"username": token.Username}

	err = sessions.FindOne(context.TODO(), filter).Decode(&sesh)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return TokenNotFound
		}
	}

	if sesh.IsExpired() {
		DeleteSession(token, sessions)
		if EmailInDatabase(token.Username, users) == nil {
			return UserNotFound
		}
		filter := bson.D{{Key: "email", Value: token.Username}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "online", Value: false}}}}

		_, err = users.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil
			}
		}
		return TokenExpired
	}
	return nil
}

// Assumes that session exists in the database.
func DeleteSession(token Session, sessions *mongo.Collection) (err error) {
	var TokenNotFound = errors.New("Token not found.")
	filter := bson.D{{Key: "username", Value: token.Username}}

	sessionFilter := bson.D{{Key: "email", Value: token.Username}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "online", Value: false}}}}

	_, err = sessions.UpdateOne(context.TODO(), sessionFilter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return TokenNotFound
		}
	}

	_, err = sessions.DeleteOne(context.TODO(), filter)
	return err
}

func GetSession(uuid uuid.UUID, sessions *mongo.Collection) (session Session) {
	var sesh Session
	filter := bson.D{{Key: "token", Value: uuid}}

	err := sessions.FindOne(context.TODO(), filter).Decode(&sesh)
	if err != nil {
		log.Error(err)
	}

	return sesh
}

// SESSION HANDLING ENDS HERE //
