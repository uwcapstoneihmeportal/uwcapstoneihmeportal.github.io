package sessions

import (
	"errors"
	"net/http"
	"fmt"
)

//necessary authorization requirements for http call.
const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	//create a new session id with a given signingKey
	sessionId, err := NewSessionID(signingKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in generating new session Id: %v", err), http.StatusInternalServerError)
		return InvalidSessionID, err
	}

	//save the session id to the cache store
	err = store.Save(sessionId, sessionState)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in uploading the session state with key %s: %v", sessionId, err), http.StatusInternalServerError)
		return InvalidSessionID, err
	}
	//include the authorization header with a proper session id and scheme.
	w.Header().Add(headerAuthorization, schemeBearer + sessionId.String())
	return sessionId, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//Check if the authorization field exists in the given headers.
	authHeader := r.Header.Get(headerAuthorization)
	//if not provided, check if the auth query was included in the url.
	if authHeader == "" {
		authHeader = r.URL.Query().Get(paramAuthorization)
	}
	if len(authHeader) > 0 {
		//check if Bearer scheme was provided
		scheme := authHeader[0:len(schemeBearer)]
		if scheme != schemeBearer {
			return InvalidSessionID, ErrInvalidScheme
		}
		//with the given session id, validate the ID.
		authHeader = authHeader[len(schemeBearer):]
		_, err := ValidateID(authHeader, signingKey)
		if err != nil {
			return InvalidSessionID, err
		}
	}
	return SessionID(authHeader), nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//get session id with the given signingKey.
	sessionId, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	//grab the session state using the session ID.
	err = store.Get(sessionId, &sessionState)
	if err != nil {
		return InvalidSessionID, err
	}
	return sessionId, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//grab session ID
	SessionId, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	//if authorized, proceed to delete from cache store.
	err = store.Delete(SessionId)
	if err != nil {
		return InvalidSessionID, err
	}
	return SessionId, nil
}
