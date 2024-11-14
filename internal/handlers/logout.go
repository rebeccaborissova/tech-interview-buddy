package handlers

import (
	"net/http"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"
)

func Logout(writer http.ResponseWriter, request *http.Request) {
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

	sessionToken := cookie.Value

	// Remove the user's session from the database collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	sessionCollection := tools.GetSessionCollection(store.DB)
	DeleteSession(/* TODO: add session struct*/, sessionCollection)
	
	// TODO set cookie to blank
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   "",
	// 	Expires: time.Now(),
	// })
}	