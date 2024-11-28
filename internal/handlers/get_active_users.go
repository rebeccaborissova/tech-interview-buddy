package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "strings"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	_ "github.com/gofrs/uuid/v5"
	log "github.com/sirupsen/logrus"
)

func getActiveUsers(writer http.ResponseWriter, request *http.Request) {
	// Get the cookie from the HTTP request
	cookie, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		api.InternalErrorHandler(writer)
		return
	}

	fmt.Println(cookie)
	// sessionToken := cookie.Value
	// sessionUUID, err := uuid.FromString(sessionToken)
	// if err != nil {
	// 	api.InternalErrorHandler(writer)
	// 	return
	// }

	// Get an instance of the database
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}

	usersCollection := tools.GetUserCollection(store.DB)
	// sessionCollection := tools.GetSessionCollection(store.DB)

	// userSession := tools.GetSession(sessionUUID, sessionCollection)
	// err = tools.CheckSession(userSession, sessionCollection, usersCollection)
	// if err != nil {
	// 	log.Error(err)
	// 	api.RequestErrorHandler(writer, err)
	// 	return
	// }

	// Get a list of active users
	activeAccounts, err := tools.GetOnlineAccounts(usersCollection)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}

	activeUsersJSON, err := json.Marshal(activeAccounts)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}

	activeUsersString := string(activeUsersJSON)

	// strings.ReplaceAll(activeUsersString, `\`, "")
	// strings.ReplaceAll(activeUsersString, "[", "{")
	// strings.ReplaceAll(activeUsersString, "]", "}")

	var response = api.SimpleResponse{
		Code:    http.StatusOK,
		Message: activeUsersString,
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
