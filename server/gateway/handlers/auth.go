package handlers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/models/users"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/indexes"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/sessions"
	"time"
	"path"
	"strconv"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"strings"
	"sort"
	"github.com/nbutton23/zxcvbn-go"
)

var bcryptCost = 13

//helper function that would store or remove users from trie structure.
func storeOrRemoveFromTrie(user *users.User, trie *indexes.Trie, remove bool) {
	names := []string{user.FirstName, user.LastName, user.UserName}
	for _, name := range names {
		for _, key := range strings.Split(name, " ") {
			if remove {
				trie.Remove(key, user.UserID)
			} else {
				trie.Add(key, user.UserID)
			}
		}
	}
}


func (ctx *SessionContext)UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, err := sessions.GetSessionID(r, ctx.signInKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("unauthorized user: %v", err), http.StatusUnauthorized)
			return
		}
		var currentState SessionState
		_, err = sessions.GetState(r, ctx.signInKey, ctx.sessionStore, &currentState)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting the current state of user: %v", err), http.StatusInternalServerError)
			return
		}
		queryInfo := r.URL.Query().Get("q")
		if len(queryInfo) < 1 {
			http.Error(w, fmt.Sprintf("missing prefix query string"), http.StatusBadRequest)
			return
		}
		userIDsWithGivenPrefix := ctx.lookUpTrie.HasPrefix(strings.ToLower(queryInfo), 20)
		var retrievedUsers []*users.User
		for _, id := range userIDsWithGivenPrefix {
			userWithPrefix, err := ctx.userStore.GetByID(id)
			if err != nil {
				http.Error(w, fmt.Sprintf("error while fetching user with given id %d: %v", id, err), http.StatusInternalServerError)
				return
			}
			retrievedUsers = append(retrievedUsers, userWithPrefix)
		}
		if retrievedUsers != nil {
			sort.Slice(retrievedUsers, func(i, j int) bool { return retrievedUsers[i].UserName < retrievedUsers[j].UserName })
		} else {
			retrievedUsers = make([]*users.User, 0)
		}
		respond(w, retrievedUsers, http.StatusOK)
	case http.MethodPost:
		err := checkContentType(r, contentTypeJSON)
		if err != nil {
			http.Error(w, fmt.Sprintf("Content-type/request body must be application/json type"),
				http.StatusUnsupportedMediaType)
			return
		}
		newUser := &users.NewUser{}
		err = decodeJsonAndPopulate(r, newUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}
		if err := newUser.Validate(); err != nil {
			http.Error(w, fmt.Sprintf("error not a valid new user profile: %v", err), http.StatusBadRequest)
			return
		}
		//password strength check
		zxcvbnRes := zxcvbn.PasswordStrength(newUser.Password, nil)
		if zxcvbnRes.Score < 3 {
			//not safely unguessable so reject request
			http.Error(w, fmt.Sprintf("Given password is not safely unguessable. Time taken to crack your password: %d",
				zxcvbnRes.CrackTime), http.StatusBadRequest)
			return
		}
		user, err := newUser.ToUser()
		if err != nil {
			http.Error(w, fmt.Sprintf("error converting new user to user profile: %v", err), http.StatusInternalServerError)
			return
		}
		//check if the new user's email and username values are unique
		uniqueUserCheck, err := ctx.userStore.GetByEmail(user.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting user by email: %v", err), http.StatusInternalServerError)
			return
		}
		if uniqueUserCheck != nil {
			http.Error(w, fmt.Sprintf("error requested new user's email or username is not unique"), http.StatusBadRequest)
			return
		}
		uniqueUserCheck, err = ctx.userStore.GetByUserName(user.UserName)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting user by username: %v", err), http.StatusInternalServerError)
			return
		}
		if uniqueUserCheck != nil {
			http.Error(w, fmt.Sprintf("error requested new user's email or username is not unique"), http.StatusBadRequest)
			return
		}
		user, err = ctx.userStore.Insert(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting new user into the database: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = sessions.BeginSession(ctx.signInKey, ctx.sessionStore, &SessionState{Time:time.Now(), AuthUser:user}, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("error initializing a new session: %v", err), http.StatusInternalServerError)
			return
		}
		storeOrRemoveFromTrie(user, ctx.lookUpTrie, false)
		respond(w, user, http.StatusCreated)
	default:
		http.Error(w, fmt.Sprintf("invalid methods"), http.StatusMethodNotAllowed)
		return
	}
}

func (ctx *SessionContext)SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	_, err := sessions.GetSessionID(r, ctx.signInKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("unauthorized user: %v", err), http.StatusUnauthorized)
		return
	}
	var currentState SessionState
	_, err = sessions.GetState(r, ctx.signInKey, ctx.sessionStore, &currentState)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting the current state of user: %v", err), http.StatusInternalServerError)
		return
	}
	requestedId := path.Base(r.URL.Path)
	var intId int64
	switch requestedId {
	case "me":
		intId = currentState.AuthUser.UserID
	default:
		intId, err = strconv.ParseInt(requestedId, 10, 64)
		if err != nil {
		http.Error(w, fmt.Sprintf("error converting the requested user id: %v", err), http.StatusBadRequest)
		return
		}
	}
	switch r.Method {
	case http.MethodGet:
		user, err := ctx.userStore.GetByID(intId)
		if err != nil {
			http.Error(w, fmt.Sprintf("error fetching the requested user profile from store: %v", err), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, fmt.Sprintf("error couldn't find the requested user profile fromt store"), http.StatusNotFound)
			return
		}
		respond(w, user, http.StatusOK)
	case http.MethodPatch:
		if requestedId != "me" && currentState.AuthUser.UserID != intId {
			http.Error(w, fmt.Sprintf("error requested id and currently authenticated user it does not match"),
				http.StatusForbidden)
			return
		}
		err = checkContentType(r, contentTypeJSON)
		if err != nil {
			http.Error(w, fmt.Sprintf("Content-type/request body must be application/json type"),
				http.StatusUnsupportedMediaType)
			return
		}
		userUpdate := &users.Updates{}
		err = decodeJsonAndPopulate(r, userUpdate)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}
		storeOrRemoveFromTrie(currentState.AuthUser, ctx.lookUpTrie, true)
		err = currentState.AuthUser.ApplyUpdates(userUpdate)
		//if fail to apply update, roll back by adding
		if err != nil {
			storeOrRemoveFromTrie(currentState.AuthUser, ctx.lookUpTrie, false)
			http.Error(w, fmt.Sprintf("error updating current user, update data not valid: %v", err), http.StatusBadRequest)
			return
		}
		updatedUser, err := ctx.userStore.Update(intId, userUpdate)
		if err != nil {
			http.Error(w, fmt.Sprintf("error updating user in the database: %v", err), http.StatusInternalServerError)
			return
		}
		//update successful so patch trie
		storeOrRemoveFromTrie(updatedUser, ctx.lookUpTrie, false)
		respond(w, updatedUser, http.StatusOK)
	default:
		http.Error(w, fmt.Sprintf("invalid methods"), http.StatusMethodNotAllowed)
		return
	}
}

func (ctx *SessionContext)SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := checkContentType(r, contentTypeJSON)
		if err != nil {
			http.Error(w, fmt.Sprintf("Content-type/request body must be application/json type"),
				http.StatusUnsupportedMediaType)
			return
		}
		userCreds := &users.Credentials{}
		err = decodeJsonAndPopulate(r, userCreds)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}
		user, err := ctx.userStore.GetByEmail(userCreds.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("error fetching user by email for authentication: %v", err), http.StatusInternalServerError)
			return
		}
		if user == nil {
			bcrypt.GenerateFromPassword([]byte(userCreds.Password), bcryptCost)
			http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
			return
		}
		if err = user.Authenticate(userCreds.Password); err != nil {
			http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
			return
		}
		_, err = sessions.BeginSession(ctx.signInKey, ctx.sessionStore, &SessionState{Time:time.Now(), AuthUser:user}, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("error initializing a new session: %v", err), http.StatusInternalServerError)
			return
		}
		respond(w, user, http.StatusCreated)
	default:
		http.Error(w, fmt.Sprintf("invalid methods"), http.StatusMethodNotAllowed)
		return
	}
}

func (ctx *SessionContext)SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		p := path.Base(r.URL.Path)
		if p != "mine" {
			http.Error(w, fmt.Sprintf("error trying to end other user's session with no permission"), http.StatusForbidden)
			return
		}
		_, err := sessions.EndSession(r, ctx.signInKey, ctx.sessionStore)
		if err != nil {
			http.Error(w, fmt.Sprintf("error trying to end session: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Add(headerContentType, contentTypeText)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "signed out")
	default:
		http.Error(w, fmt.Sprintf("invalid methods"), http.StatusMethodNotAllowed)
		return
	}
}

func decodeJsonAndPopulate(r *http.Request, value interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return err
	}
	return nil
}

func checkContentType(r *http.Request, contentType string) error {
	if r.Header.Get(headerContentType) != contentType {
		return errors.New(fmt.Sprintf("Content-type/request body must be application/json type"))
	}
	return nil
}
