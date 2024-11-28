package handlers

import (
	"encoding/json"
	"net/http"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	"github.com/gofrs/uuid/v5"
	log "github.com/sirupsen/logrus"
)

func getUserInfo(writer http.ResponseWriter, request *http.Request) {
	// Get the session cookie from the HTTP request
	cookie, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get the session UUID from the cookie
	sessionToken := cookie.Value
	sessionUUID, err := uuid.FromString(sessionToken)
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	
	// Get an instance of the user and session collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	usersCollection := tools.GetSessionCollection(store.DB)
	sessionCollection := tools.GetSessionCollection(store.DB)

	// Get an instance of the user's session
	userSession := tools.GetSession(sessionUUID, sessionCollection)
	userAccount := tools.EmailInDatabase(userSession.Username, usersCollection)
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}

	var response = api.UserInfoResponse{
		Code:		http.StatusOK,
		Username:	userAccount.Email,
		FirstName:	userAccount.FirstName,
		LastName: 	userAccount.LastName,
		InvitedBy:	userAccount.InvitedBy,
		TakenDSA:	userAccount.TakenDSA,
		Year:		userAccount.Year,
		Description: userAccount.Description,
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
