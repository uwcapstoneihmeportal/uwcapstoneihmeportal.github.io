package handlers

import (
	"testing"
	"net/http"
	"time"
	"os"
	"github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"encoding/json"
	"bytes"
	"net/http/httptest"
	"fmt"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/models/users"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/indexes"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/sessions")

func TestCorsHandler_ServeHTTP(t *testing.T) {
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
	hctx := NewContext("testKey", sessStore, userStore, testTrie)
	nu := users.NewUser{"kwontae1@uw.edu", "1234567", "1234567", "test1",
		"tester","testing"}
	body, _ := json.Marshal(&nu)
	cases := []struct {
		name string
		methodType string
		allowOrigin string
		allowMethods string
		allowHeaders string
		exposeHeaders string
		maxAge string
	}{
		{"Valid Cors",
		http.MethodPost,
			originAny,
			allowMethods,
			allowHeaders,
			authHeader,
			maxAge,
		},
		{"Method Options shouldn't do anything, but still get the correct header values written in",
			http.MethodOptions,
			originAny,
			allowMethods,
			allowHeaders,
			authHeader,
			maxAge,
		},
	}
	for _, c := range cases {
		req, err := http.NewRequest(c.methodType, "/users", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add(headerContentType, contentTypeJSON)
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		ch := NewCorsHandler(http.HandlerFunc(hctx.UsersHandler))
		ch.ServeHTTP(rr, req)
		if rr.Header().Get(headerAccessControlAllowOrigin) != c.allowOrigin &&
			rr.Header().Get(headerAccessControlAllowMethods) != c.allowMethods &&
				rr.Header().Get(headerAccessControlAllowHeaders) != c.allowHeaders &&
					rr.Header().Get(headerAccessControlExposeHeaders) != c.exposeHeaders &&
						rr.Header().Get(headerAccessControlMaxAge) != c.maxAge {
							t.Errorf(fmt.Sprintf("case %s error missing headers", c.name	))
		}
	}
}