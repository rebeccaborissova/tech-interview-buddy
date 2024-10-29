package handlers

import (
	"encoding/json"
	"net/http"

	"GO_PRACTICE_PROJECT/api"
	"GO_PRACTICE_PROJECT/internal/tools"

	// "github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func getLoginReponse(writer http.ResponseWriter, router *http.Request) {
	var params = api.LoginParams{}
	// var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	params.Username = router.Header.Get("Username")
	// err = decoder.Decode(&params, router.Header.Get("Username"))
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}

	// Get an instance of the database
	var store *tools.PostgresStore
	store, err = tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}

	usersCollection := store.DB.Collection("users")
	account := tools.EmailInDatabase(params.Username, usersCollection)
	if account == nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}

	var response = api.LoginResponse{
		Code:     http.StatusOK,
		Username: params.Username,
		Message:  "Successfully authenticated " + params.Username,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
	writer.Header().Set("Access-Control-Max-Age", "15")
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
