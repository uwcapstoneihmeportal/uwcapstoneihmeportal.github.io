package handlers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"github.com/nimajalali/go-force/sobjects"
)

//successful response interface for user authorization
type AuthToken struct {
	Default_password bool `json:"default_password"`
	Access_token string `json:"access_token"`
	Token_type string `json:"token_type"`
	Sojbect_url string `json:"sobject_url"`
}

//login input json interface from the client
type Credentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

//payload interface used to send patch request to Salesforce REST Api
type NewPassword struct {
	NewPassHash string	`json:"password__c"`
}

//password update input json interface from the client
type PasswordRequest struct {
	Access_token string `json:"access_token"`
	New_password string `json:"new_password"`
	Sobject_url string `json:"sobject_url"`
}

//specific record interface to decode json response from Salesforce query
type Record struct {
	sobjects.BaseSObject
	Email string `force:"Email"`
	Password string `force:"password__c"`
}

//interface to decode json response from Salesforce query
type Records struct {
	sobjects.BaseQuery
	Records []Record `force:"records"`
}

var bcryptCost = 13

func (ctx *SessionContext)Authorize(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//fetch email and password from salesforce contact table
		//compares the bcrypted password input with the pre-encrypted password fetched from salesforce
		//if the password is valid, then check if the password is a default value.
		//if default, then set a flag to remind the client that the user needs to update the password.
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
			if userCreds.Password == "collaborator2019" {
				authResponse.Default_password = true
			}
		}
		//url used for possible password update.
		authResponse.Sojbect_url = salesforceCredentialSObjects.Records[0].Attributes.Url
		authResponse.Access_token = ctx.forceApi.GetAccessToken()
		authResponse.Token_type = "Bearer"
		respond(w, authResponse, http.StatusOK)
	default:
		http.Error(w, fmt.Sprintf("invalid methods"), http.StatusMethodNotAllowed)
		return
	}
}

func (ctx *SessionContext)PasswordUpdate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := checkContentType(r, contentTypeJSON)
		if err != nil {
			http.Error(w, fmt.Sprintf("Content-type/request body must be application/json type"),
				http.StatusUnsupportedMediaType)
			return
		}
		newPassword := &PasswordRequest{}
		err = decodeJsonAndPopulate(r, newPassword)
		if err != nil {
			http.Error(w, fmt.Sprintf("username does not exist: %v", err), http.StatusBadRequest)
			return
		}
		//must have a valid access token.
		if newPassword.Access_token != ctx.access_code {
			http.Error(w, fmt.Sprint("unauthorized user"), http.StatusUnauthorized)
			return
		}
		newPassHash, err := bcrypt.GenerateFromPassword([]byte(newPassword.New_password), bcryptCost)
		if err != nil {
			http.Error(w, fmt.Sprintf("error generating encrypted hash: %v", err),
				http.StatusBadRequest)
			return
		}
		rawPatchPayload := `{"password__c":` + string(newPassHash)
		unMarshPatchPayload := NewPassword{}
		json.Unmarshal([]byte(rawPatchPayload), &unMarshPatchPayload)
		err = ctx.forceApi.Patch(newPassword.Sobject_url, nil, unMarshPatchPayload, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("error updating password: %v", err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, "Password updated.")
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
