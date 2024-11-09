// Credit to https://www.youtube.com/watch?v=8uiZC0l4Ajw for basic Go API

package main

import (
	"GO_PRACTICE_PROJECT/internal/handlers"
	"GO_PRACTICE_PROJECT/internal/tools"
	"fmt"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Initial DB setup
	store, err := tools.NewPostgresStore() // Store is a *PostgresStore that is *Mongo.Database
	if err != nil {
		log.Fatal(err)
	}

	usersCollection := store.DB.Collection("users")

	// Test attributes
	email := "sarah.tran1@ufl.edu"
	password := "meowmeowmeoW@13"
	firstName := "Sarah"
	lastName := "Tran"
	tookDSA := true
	classYear := 2

	// CREATE OPERATION
	// TODO: Should I try to make it return an error and like how should I let front-end know what information given was invalid?
	// TEMPORARILY: Use a string
	tools.InsertAccount(email, password, firstName, lastName, tookDSA, classYear, usersCollection)
	// CREATE OPERATION

	fmt.Println()

	// READ OPERATION (https://www.youtube.com/watch?v=Ap8elI7ePt4&t=7s)
	email = "meowmeow@ufl.edu"
	account := tools.EmailInDatabase(email, usersCollection)
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
	classYear = 1

	tools.InsertAccount(email, password, firstName, lastName, tookDSA, classYear, usersCollection)
	deleteErr := tools.DeleteAccount(email, usersCollection)
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
	classYear = 3

	tools.InsertAccount(email, password, firstName, lastName, tookDSA, classYear, usersCollection)
	password, _ = argon2id.CreateHash(password, argon2id.DefaultParams)
	fmt.Println(password)

	// UPDATE OPERATION

	fmt.Println()
	fmt.Println("I am not dead. Yay!")

	// Initial API setup
	var router *chi.Mux = chi.NewRouter()
	handlers.Handler(router)

	fmt.Println("Starting API service on port 8000...")
	fmt.Println(`
_____/\\\\\\\\\\\\_______/\\\\\____________________/\\\\\\\\\_____/\\\\\\\\\\\\\____/\\\\\\\\\\\_        
 ___/\\\//////////______/\\\///\\\________________/\\\\\\\\\\\\\__\/\\\/////////\\\_\/////\\\///__       
  __/\\\_______________/\\\/__\///\\\_____________/\\\/////////\\\_\/\\\_______\/\\\_____\/\\\_____      
   _\/\\\____/\\\\\\\__/\\\______\//\\\___________\/\\\_______\/\\\_\/\\\\\\\\\\\\\/______\/\\\_____     
    _\/\\\___\/////\\\_\/\\\_______\/\\\___________\/\\\\\\\\\\\\\\\_\/\\\/////////________\/\\\_____    
     _\/\\\_______\/\\\_\//\\\______/\\\____________\/\\\/////////\\\_\/\\\_________________\/\\\_____   
      _\/\\\_______\/\\\__\///\\\__/\\\______________\/\\\_______\/\\\_\/\\\_________________\/\\\_____  
       _\//\\\\\\\\\\\\/_____\///\\\\\/_______________\/\\\_______\/\\\_\/\\\______________/\\\\\\\\\\\_ 
        __\////////////_________\/////_________________\///________\///__\///______________\///////////__
		`)

	var serverErr error = http.ListenAndServe("localhost:8000", router)
	if serverErr != nil {
		log.Error(serverErr)
	}

}
