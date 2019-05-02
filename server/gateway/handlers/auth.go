package handlers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"github.com/nimajalali/go-force/sobjects"
)

type Credentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type Record struct {
	sobjects.BaseSObject

	Email string `force:"Email"`
	Password string `force:"password__c"`
}

type Records struct {
	sobjects.BaseQuery

	Records []Record `force:"records"`
}

type AuthToken struct {
	Default_password bool `json:"default_password"`
	Access_token string `json:"access_token"`
	Token_type string `json:"token_type"`
}

var bcryptCost = 13

func (ctx *SessionContext)Authorize(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//fetch email and password from salesforce contact table
		//if the password field is empty, then the user needs to create a password.
		//Authenticate compares the plaintext password against the stored hash
		//and returns an error if they don't match, or nil if they do
		err := checkContentType(r, contentTypeJSON)
		if err != nil {
			http.Error(w, fmt.Sprintf("Content-type/request body must be application/json type"),
				http.StatusUnsupportedMediaType)
			return
		}
		userCreds := &Credentials{}
		err = decodeJsonAndPopulate(r, userCreds)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}
		salesforceCredentialSObjects := &Records{}
		err = ctx.forceApi.Query("SELECT Email, password__c FROM Contact WHERE Email=" +
			userCreds.Email, salesforceCredentialSObjects)
		if err != nil {
			http.Error(w, fmt.Sprintf("username does not exist: %v", err), http.StatusBadRequest)
			return
		}
		authResponse := &AuthToken{}
		//check password hash
		err = bcrypt.CompareHashAndPassword([]byte(salesforceCredentialSObjects.Records[0].Password),
			[]byte(userCreds.Password))
		if err != nil {
			http.Error(w, fmt.Sprintf("Incorrect password: %v", err), http.StatusUnauthorized)
			return
		} else {
			//if password is collaborator2019, then the user needs to be prompted to change
			if userCreds.Password == "collaborators2019" {
				authResponse.Default_password = true
			}
		}
		authResponse.Access_token = ctx.forceApi.GetAccessToken()
		authResponse.Token_type = "Bearer"
		respond(w, authResponse, http.StatusOK)
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
