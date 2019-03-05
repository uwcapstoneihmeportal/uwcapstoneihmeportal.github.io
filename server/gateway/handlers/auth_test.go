package handlers

import (
	"testing"
	"net/http"
	"fmt"
	"net/http/httptest"
	"time"
	"os"
	"log"
	"github.com/go-sql-driver/mysql"
	"database/sql"
	"encoding/json"
	"bytes"
	"strconv"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"strings"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/models/users"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/indexes"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/sessions"
)

func reqEnv(name string) string {
	val := os.Getenv(name)
	if len(val) == 0{
		log.Fatalf("please set %s variable", name)
	}
	return val
}

func TestSessionContext_UsersHandler(t *testing.T) {
	sessStore := sessions.NewMemStore(time.Hour, time.Minute * 3)
	mysqlAddr := reqEnv("MYSQL_ADDR")
	mysqlDB := reqEnv("MYSQL_DATABASE")
	mysqlPwd := reqEnv("MYSQL_ROOT_PASSWORD")
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}
	config := mysql.Config{
		Addr: mysqlAddr,
		User: "root",
		Passwd: mysqlPwd,
		DBName: mysqlDB,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}
	defer db.Close()
	userStore := users.NewMySQLStore(db)
	testTrie := indexes.NewTrie()
	hctx := NewContext("testKey", sessStore, userStore, testTrie, nil)
	validUser := users.User{int64(1), "test@uw.edu", []byte{1,2,3,4,5,6,7}, "test10",
		"tester", "testing", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"}
	u, err := hctx.userStore.Insert(&validUser)
	cases := []struct {
		name string
		query string
		methodType string
		newUserProfile users.NewUser
		expectedStatusCode int
		expectedContentType string
		expectedUser users.User
	}{
		{
			"Invalid Method",
			"",
			http.MethodDelete,
			users.NewUser{"kwontae@uw.edu", "1234","1234","test1",
			"tester","testing",},
			http.StatusMethodNotAllowed,
			contentTypeText,
			users.User{1, "kwontae@uw.edu", []byte{1,2,3,4}, "test1",
			"tester", "testing", ""},
		},
		{
			"Valid POST",
			"",
			http.MethodPost,
			users.NewUser{"kwontae45@uw.edu", "1234567", "1234567", "test45",
			"tester","testing"},
			http.StatusCreated,
			contentTypeJSON,
			users.User{1, "kwontae45@uw.edu", []byte{1,2,3,4,5,6,7}, "test45",
			"tester", "testing", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"},
		},
		{
			"POST with a bad new user profile",
			"",
			http.MethodPost,
			users.NewUser{"kwontae2@uw.edu", "12345", "12345", "test2",
			"tester2","testing2"},
			http.StatusBadRequest,
			contentTypeText,
			users.User{2, "kwontae2@uw.edu", []byte{1,2,3,4,5}, "test2", "tester2",
			"testing2", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"},
		},
		{
			"Valid POST with a same email",
			"",
			http.MethodPost,
			users.NewUser{"test@uw.edu", "1234567", "1234567", "test2",
				"tester2","testing2"},
				http.StatusBadRequest,
				contentTypeText,
			users.User{1, "test@uw.edu", []byte{1,2,3,4,5,6,7}, "test2",
				"tester2", "testing2", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"},
		},
		{ "Valid POST with a same username",
		"",
		http.MethodPost,
			users.NewUser{"kwontae3@uw.edu", "1234567", "1234567", "test10",
				"tester1","testing1"},
			http.StatusBadRequest,
			contentTypeText,
			users.User{1, "kwontae3@uw.edu", []byte{1,2,3,4,5,6,7}, "test1",
				"tester1", "testing1", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"},
		},
	}
	for _, c:= range cases {
		URL := fmt.Sprintf("/users")
		//create a body to post
		body, _ := json.Marshal(&c.newUserProfile)
		req := httptest.NewRequest(c.methodType, URL, bytes.NewReader(body))
		req.Header.Add(headerContentType, contentTypeJSON)
		respRec := httptest.NewRecorder()
		hctx.UsersHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d, but got %d", c.name, c.expectedStatusCode, resp.StatusCode)
		}
		//check the content type header
		contentType := resp.Header.Get(headerContentType)
		if contentType != c.expectedContentType {
			t.Errorf("case %s: incorrect Content-Type header: expected %s, but got %s", c.name, c.expectedContentType, contentType)
		}
		if resp.StatusCode == http.StatusCreated {
			var returnedUser users.User
			if err := json.NewDecoder(resp.Body).Decode(&returnedUser); err != nil {
				t.Errorf(fmt.Sprintf("case %s: error decoding JSON: %v", c.name, err), http.StatusBadRequest)
			}
			//delete so it can be run without restarting the database instance
			hctx.userStore.Delete(int64(returnedUser.UserID))
			if returnedUser.FirstName != c.expectedUser.FirstName || returnedUser.UserName != c.expectedUser.UserName ||
				returnedUser.LastName != c.expectedUser.LastName {
				t.Errorf(fmt.Sprintf("case %s: error expected email: %s and username: %s to match, but got: %s and %s",
					c.name, c.expectedUser.Email, c.expectedUser.UserName, returnedUser.Email, returnedUser.UserName))
			}
		}
	}
	hctx.userStore.Delete(u.UserID)
}

func TestSessionContext_SpecificUserHandler(t *testing.T) {
	sessStore := sessions.NewMemStore(time.Hour, time.Minute * 3)
	mysqlAddr := reqEnv("MYSQL_ADDR")
	mysqlDB := reqEnv("MYSQL_DATABASE")
	mysqlPwd := reqEnv("MYSQL_ROOT_PASSWORD")
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}
	config := mysql.Config{
		Addr: mysqlAddr,
		User: "root",
		Passwd: mysqlPwd,
		DBName: mysqlDB,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}
	defer db.Close()
	userStore := users.NewMySQLStore(db)
	testTrie := indexes.NewTrie()
	hctx := NewContext("testKey", sessStore, userStore, testTrie, nil)
	sid, err := sessions.NewSessionID(hctx.signInKey)
	if err != nil {
		t.Fatalf("error generating SessionID: %v", err)
	}
	//insert some valid users
	validUser := users.User{int64(1), "test@uw.edu", []byte{1,2,3,4,5,6,7}, "test10",
		"tester", "testing", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"}
	u, err := hctx.userStore.Insert(&validUser)
	validUser2 := users.User{int64(2), "test2@uw.edu", []byte{1,2,3,4,5,6,7}, "test11",
		"tester", "testing", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"}
	u2, err := hctx.userStore.Insert(&validUser2)
	//begin session for this user
	err = hctx.sessionStore.Save(sid, &SessionState{Time:time.Now(), AuthUser:&validUser})
	if err != nil {
		t.Errorf(fmt.Sprintf("Error in uploading the session state with key %s: %v", sid, err))
		return
	}
	getCases := []struct {
		name string
		header string
		methodType string
		query string
		expectedUser users.User
		expectedContentType string
		expectedStatusCode int
	}{
		{
			"Valid SessionID and Scheme with GET method",
			schemeBearer + string(sid),
			http.MethodGet,
			strconv.FormatInt(u2.UserID, 10),
			validUser2,
			contentTypeJSON,
			http.StatusOK,
		},
		{
			"Valid SessionID and Scheme with GET method",
			schemeBearer + string(sid),
			http.MethodGet,
			"me",
			validUser,
			contentTypeJSON,
			http.StatusOK,
		},
		{
			"Valid SessionID and Scheme with Get method with no user",
			schemeBearer + string(sid),
			http.MethodGet,
			strconv.FormatInt(int64(0),10),
			validUser2,
			contentTypeText,
			http.StatusNotFound,
		},
		{
			"Invalid Method",
			schemeBearer + string(sid),
			http.MethodDelete,
			"2",
			validUser2,
			contentTypeText,
			http.StatusMethodNotAllowed,
		},
	}
	for _, c := range getCases {
		urlPath := fmt.Sprintf("/users/%s", c.query)
		req := httptest.NewRequest(c.methodType, urlPath, nil)
		req.Header.Add(authHeader,c.header)
		req.Header.Add(headerContentType, contentTypeJSON)
		w := httptest.NewRecorder()
		hctx.SpecificUserHandler(w, req)
		resp := w.Result()
		
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d, but got %d", c.name, c.expectedStatusCode, resp.StatusCode)
		}
		//check the content type header
		contentType := resp.Header.Get(headerContentType)
		if contentType != c.expectedContentType {
			t.Errorf("case %s: incorrect Content-Type header: expected %s, but got %s", c.name, c.expectedContentType, contentType)
		}
		if c.expectedStatusCode == http.StatusOK {
			var returnedUser users.User
			if err := json.NewDecoder(resp.Body).Decode(&returnedUser); err != nil {
				t.Errorf(fmt.Sprintf("error decoding JSON: %v", err))
				return
			}
			if returnedUser.FirstName != c.expectedUser.FirstName || returnedUser.UserName != c.expectedUser.UserName ||
				returnedUser.LastName != c.expectedUser.LastName {
				t.Errorf(fmt.Sprintf("case %s: error expected email: %s and username: %s to match, but got: %s and %s",
					c.name, c.expectedUser.Email, c.expectedUser.UserName, returnedUser.Email, returnedUser.UserName))
			}
		}
	}
	patchCases := []struct {
		name string
		header string
		methodType string
		query string
		update users.Updates
		expectedUser users.User
		expectedContentType string
		expectedStatusCode int
	}{
		{
			"Patch Method with forbidden access",
			schemeBearer + string(sid),
			http.MethodPatch,
			strconv.FormatInt(u2.UserID, 10),
			users.Updates{"updatedFirstname", "updatedLastName"},
			users.User{u2.UserID, u2.Email, u2.PassHash, u2.UserName,
				"updatedFirstname", "updatedLastName", u2.PhotoURL},
			contentTypeText,
			http.StatusForbidden,
		},
		{
			"Patch Method with 'me' ",
			schemeBearer + string(sid),
			http.MethodPatch,
			"me",
			users.Updates{"updatedFirstname", "updatedLastName"},
			users.User{u.UserID, u.Email, u.PassHash, u.UserName,
				"updatedFirstname", "updatedLastName", u.PhotoURL},
			contentTypeJSON,
			http.StatusOK,
		},
		{
			"Patch Method with a valid access",
			schemeBearer + string(sid),
			http.MethodPatch,
			strconv.FormatInt(u.UserID, 10),
			users.Updates{"updatedFirstname2", "updatedLastName2"},
			users.User{u.UserID, u.Email, u.PassHash, u.UserName,
				"updatedFirstname2", "updatedLastName2", u.PhotoURL},
			contentTypeJSON,
			http.StatusOK,
		},
		{
			"Patch Method with a invalid update",
			schemeBearer + string(sid),
			http.MethodPatch,
			strconv.FormatInt(u.UserID, 10),
			users.Updates{"", ""},
			users.User{u.UserID, u.Email, u.PassHash, u.UserName,
				"updatedFirstname", "updatedLastName", u.PhotoURL},
			contentTypeText,
			http.StatusBadRequest,
		},
	}
	for _, c := range patchCases {
		urlPath := fmt.Sprintf("/users/%s", c.query)
		//create a body to post
		body, _ := json.Marshal(&c.update)
		req := httptest.NewRequest(c.methodType, urlPath, bytes.NewReader(body))
		req.Header.Add(headerContentType, contentTypeJSON)
		req.Header.Add(authHeader,c.header)
		w := httptest.NewRecorder()
		hctx.SpecificUserHandler(w, req)
		resp := w.Result()
		
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d, but got %d", c.name, c.expectedStatusCode, resp.StatusCode)
		}
		//check the content type header
		contentType := resp.Header.Get(headerContentType)
		if contentType != c.expectedContentType {
			t.Errorf("case %s: incorrect Content-Type header: expected %s, but got %s", c.name, c.expectedContentType, contentType)
		}
		if c.expectedStatusCode == http.StatusOK {
			var returnedUser users.User
			if err := json.NewDecoder(resp.Body).Decode(&returnedUser); err != nil {
				t.Errorf(fmt.Sprintf("error decoding JSON: %v", err))
				return
			}
			if returnedUser.FirstName != c.expectedUser.FirstName || returnedUser.UserName != c.expectedUser.UserName ||
				returnedUser.LastName != c.expectedUser.LastName {
				t.Errorf(fmt.Sprintf("case %s: error expected first name: %s, last name: %s," +
					" and username: %s to match, but got: %s, %s and %s",
					c.name, c.expectedUser.FirstName,c.expectedUser.LastName, c.expectedUser.UserName,
						returnedUser.FirstName, returnedUser.LastName, returnedUser.UserName))
			}
		}
	}
	hctx.userStore.Delete(u.UserID)
	hctx.userStore.Delete(u2.UserID)
}

func TestSessionContext_SessionsHandler(t *testing.T) {
	sessStore := sessions.NewMemStore(time.Hour, time.Minute * 3)
	mysqlAddr := reqEnv("MYSQL_ADDR")
	mysqlDB := reqEnv("MYSQL_DATABASE")
	mysqlPwd := reqEnv("MYSQL_ROOT_PASSWORD")
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}
	config := mysql.Config{
		Addr: mysqlAddr,
		User: "root",
		Passwd: mysqlPwd,
		DBName: mysqlDB,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}
	defer db.Close()
	userStore := users.NewMySQLStore(db)
	testTrie := indexes.NewTrie()
	hctx := NewContext("testKey", sessStore, userStore, testTrie, nil)
	passHash, _ := bcrypt.GenerateFromPassword([]byte("1234567"),13)
	//insert some valid users
	validUser := users.User{int64(1), "test@uw.edu", passHash, "test10",
		"tester", "testing", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"}
	u, err := hctx.userStore.Insert(&validUser)
	validUser2 := users.User{int64(2), "test2@uw.edu", passHash, "test11",
		"tester", "testing", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"}
	u2, err := hctx.userStore.Insert(&validUser2)
	cases := []struct {
		name string
		methodType string
		query string
		cred users.Credentials
		headerContentType string
		expectedContentType string
		expectedStatusCode int
		expectedUser users.User
	}{
		{
			"Valid POST",
			http.MethodPost,
			"",
			users.Credentials{"test@uw.edu", "1234567"},
			contentTypeJSON,
			contentTypeJSON,
			http.StatusCreated,
			validUser,
		},
		{
			"Invalid method",
			http.MethodGet,
			"",
			users.Credentials{"", ""},
			contentTypeJSON,
			contentTypeText,
			http.StatusMethodNotAllowed,
			validUser,
		},
		{
			"Invalid POST with bad header content type",
			http.MethodPost,
			"",
			users.Credentials{"test@uw.edu","1234567"},
			contentTypeText,
			contentTypeText,
			http.StatusUnsupportedMediaType,
			validUser,
		},
		{
			"Invalid POST with bad email credentials",
			http.MethodPost,
			"",
			users.Credentials{"tester@example.com", "1234567"},
			contentTypeJSON,
			contentTypeText,
			http.StatusUnauthorized,
			users.User{},
		},
		{
			"Invalid POST with bad password credentials",
			http.MethodPost,
			"",
			users.Credentials{"test@uw.edu", "12345678"},
			contentTypeJSON,
			contentTypeText,
			http.StatusUnauthorized,
			users.User{},
		},
	}

	for _, c := range cases {
		URL := fmt.Sprintf("/sessions")
		//create a body to post
		body, _ := json.Marshal(&c.cred)
		req := httptest.NewRequest(c.methodType, URL, bytes.NewReader(body))
		req.Header.Add(headerContentType, c.headerContentType)
		respRec := httptest.NewRecorder()
		hctx.SessionsHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d, but got %d", c.name, c.expectedStatusCode, resp.StatusCode)
		}
		//check the content type header
		contentType := resp.Header.Get(headerContentType)
		if contentType != c.expectedContentType {
			t.Errorf("case %s: incorrect Content-Type header: expected %s, but got %s", c.name, c.expectedContentType, contentType)
		}
		if resp.StatusCode == http.StatusCreated {
			var returnedUser users.User
			if err := json.NewDecoder(resp.Body).Decode(&returnedUser); err != nil {
				t.Errorf(fmt.Sprintf("case %s: error decoding JSON: %v", c.name, err), http.StatusBadRequest)
			}
			if returnedUser.FirstName != c.expectedUser.FirstName || returnedUser.UserName != c.expectedUser.UserName ||
				returnedUser.LastName != c.expectedUser.LastName {
				t.Errorf(fmt.Sprintf("case %s: error expected email: %s and username: %s to match, but got: %s and %s",
					c.name, c.expectedUser.Email, c.expectedUser.UserName, returnedUser.Email, returnedUser.UserName))
			}
		}
	}
	//delete users
	hctx.userStore.Delete(u.UserID)
	hctx.userStore.Delete(u2.UserID)
}

func TestSessionContext_SpecificSessionHandler(t *testing.T) {
	sessStore := sessions.NewMemStore(time.Hour, time.Minute * 3)
	mysqlAddr := reqEnv("MYSQL_ADDR")
	mysqlDB := reqEnv("MYSQL_DATABASE")
	mysqlPwd := reqEnv("MYSQL_ROOT_PASSWORD")
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}
	config := mysql.Config{
		Addr: mysqlAddr,
		User: "root",
		Passwd: mysqlPwd,
		DBName: mysqlDB,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}
	defer db.Close()
	userStore := users.NewMySQLStore(db)
	testTrie := indexes.NewTrie()
	hctx := NewContext("testKey", sessStore, userStore, testTrie, nil)
	passHash, _ := bcrypt.GenerateFromPassword([]byte("1234567"),13)
	//insert a user
	validUser := users.User{int64(1), "test@uw.edu", passHash, "test10",
		"tester", "testing", "https://www.gravatar.com/avatar/cd305e38164a0ea6e7873b1e6c722090"}
	u, err := hctx.userStore.Insert(&validUser)
	sid, err := sessions.NewSessionID(hctx.signInKey)
	if err != nil {
		t.Fatalf("error generating SessionID: %v", err)
	}
	//begin session for this user
	err = hctx.sessionStore.Save(sid, &SessionState{Time:time.Now(), AuthUser:&validUser})
	if err != nil {
		t.Errorf(fmt.Sprintf("Error in uploading the session state with key %s: %v", sid, err))
		return
	}
	cases := []struct {
		name string
		header string
		methodType string
		query string
		expectedOutput string
		expectedStatusCode int
		expectedContentType string
	}{
		{
			"Valid Delete",
			schemeBearer + string(sid),
			http.MethodDelete,
			"mine",
			"signed out",
			http.StatusOK,
			contentTypeText,
		},
		{
			"Delete with invalid query",
			schemeBearer + string(sid),
			http.MethodDelete,
			"wrongQuery",
			"signed out",
			http.StatusForbidden,
			contentTypeText,
		},
		{
			"Delete with no header auth",
			"",
			http.MethodDelete,
			"mine",
			"signed out",
			http.StatusInternalServerError,
			contentTypeText,
		},
		{
			"Invalid Method",
			schemeBearer + string(sid),
			http.MethodPost,
			"mine",
			"signed out",
			http.StatusMethodNotAllowed,
			contentTypeText,
		},
	}
	for _, c:= range cases {
		urlPath := fmt.Sprintf("/sessions/%s", c.query)
		//create a body to post
		req := httptest.NewRequest(c.methodType, urlPath, nil)
		req.Header.Add(authHeader,c.header)
		w := httptest.NewRecorder()
		hctx.SpecificSessionHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d, but got %d", c.name, c.expectedStatusCode, resp.StatusCode)
		}
		//check the content type header
		contentType := resp.Header.Get(headerContentType)
		if contentType != c.expectedContentType {
			t.Errorf("case %s: incorrect Content-Type header: expected %s, but got %s", c.name, c.expectedContentType, contentType)
		}
		if c.expectedStatusCode == http.StatusOK {
			respData, _ := ioutil.ReadAll(resp.Body)
			if strings.Compare(c.expectedOutput, string(respData)) != 0 {
				t.Errorf("case %s: incorrect output. Expected %s, but got %s", c.name, c.expectedOutput, string(respData))
			}
		}
	}
	//delete user
	hctx.userStore.Delete(u.UserID)
}
