package handlers

import (
	"encoding/json"
	"net/http"

	"GO_PRACTICE_PROJECT/api"
	"GO_PRACTICE_PROJECT/internal/tools"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func getLoginReponse(writer http.ResponseWriter, router *http.Request) {
	var params = api.LoginParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, router.Header)
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
		Code: http.StatusOK,
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
